package main

import (
    "os"
    "net"
    "time"
    "errors"
    "strings"
    "syscall"
    "encoding/json"

    xenapi "github.com/terra-farm/go-xen-api-client"
)

func mount(mountDir string, params *jsonParams) {
    var err error

    device = attach(params)

    if err = mountdevice(device, mountdir); err != nil {
        failure(err)
    }

    if err = saveParams(mountDir, params); err != nil {
        failure(err)
    }

    success()
}

func saveParams(mountDir string, params *jsonParams) error {
    debug("ioutil.WriteFile")
    byt := json.Marshal(params)
    if err := ioutil.WriteFile(fmt.Sprinf("%s.json", mountDir), byt, 0600); err != nil {
        return err
    }

    return nil
}

func attach(params *jsonParams) error {
    if params.FSType == "" {
        params.FSType = "ext4"
    }

    var mode xenapi.VbdMode
    switch params.ReadWrite {
    case "ro":
        mode = xenapi.VbdModeRO
    case "rw":
        mode = xenapi.VbdModeRW
    default:
        failure(errors.New("Unkown ReadWrite"))
    }

    xapi, session, err := xapiLogin(params.XenServerHost, params.XenServerUsername, params.XenServerPassword)
    if err != nil {
        failure(fmt.Errof("Could not login at XenServer, error: %s", err.Error()))
    }

    defer func() {
        if err := xapiLogout(xapi, session); err != nil {
            failure(fmt.Errorf("Failed to logout from XenServer, error: %s", err.Error()))
        }
    }()

    vm, err := getVM(xapi, session)
    if err != nil {
        failure(err)
    }

    debug("VM.GetAllowedVBDDevices")
    vbdDevices, err := xapi.VM.GetAllowedVBDDevices(session, vm)
    if err != nil {
        failure(err)
    }

    if len(vbdDevices) < 1 {
        failure(errors.New("no VBD devices are available anymore"))
    }

    debug("VDI.GetAllRecords")
    vdis, err := xapi.VDI.GetAllRecords(session)
    if err != nil {
        failure(err)
    }

    var vdiUUID xenapi.VDIRef
    for ref, vdi := range vdis {
        if vdi.NameLabel == options.PVOrVolumeName && !vdi.IsASnapshot {
            vdiUUID = ref
        }
    }

    if vdiUUID == "" {
        failure(errors.New("Could not find VDI"))
    }

    params.VDIUUID = vdiUUID

    debug("VBD.GetAllRecords")
    vbds, err := xapi.VBD.GetAllRecords(session)
    if err != nil {
        failure(err)
    }

    for ref, vbd := range vbds {
        if vbd.VDI == vdiUUID && vbd.CurrentlyAttached {
            debug("Attempting to safely detach VDI")
            time.Sleep(10 * time.Second)
            if err := detachVBD(ref, xapi, session); err != nil {
                debug("Failed at detaching VDI, will try again soon")
                failure(err)
            }
        }
    }

    debug("VBD.Create")
    vbdUUID, err := xapi.VBD.Create(session, xenapi.VBDRecord{
        Bootable:    false,
        Mode:        mode,
        Type:        xenapi.VbdTypeDisk,
        Unpluggable: true,
        Userdevice:  vbdDevices[0],
        VDI:         vdiUUID,
        VM:          vm,
    })

    if err != nil {
        failure(err)
    }

    debug("VBD.Plug")
    if err = xapi.VBD.Plug(session, vbdUUID); err != nil {
        failure(err)
    }

    debug("VBD.GetDevice")
    device, err := xapi.VBD.GetDevice(session, vbdUUID)
    if err != nil {
        failure(err)
    }

    return fmt.Sprinf("/dev/%s", device)
}

func mountdevice(devicePath, mountDir string) error {
    blkid, err := run("blkid", devicePath)
    if err != nil && !strings.Contains(err.Error(), "exit status 2") {
        return err
    }

    if blkid == "" {
        if _, err := run("mkfs", "-t", params.FSType, devicePath); err != nil {
            return err
        }
    }

    debug("os.MkdirAll")
    if err = os.MkdirAll(mountDir, 0755); err != nil {
        return err
    }

    debug("syscall.Mount")
    if err = syscall.Mount(devicePath, mountDir, params.FSType, 0, ""); err != nil {
        return err
    }

    return nil
}
