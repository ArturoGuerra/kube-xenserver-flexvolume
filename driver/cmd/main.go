package main

import (
    "fmt"
    "os"
    "bytes"
    "errors"
    "encoding/json"
    "syscall"
    "strings"
    "net"
    "time"

    xenapi "github.com/terra-farm/go-xen-api-client"
)


type jsonParams struct {
    FSGroup string `json:"kubernetes.io/fsGroup"`
    FSType  string `json:"kubernetes.io/fsType"`
    PVOrVolumeName string `json:"kubernetes.io/pvOrVolumeName"`
    PodName string `json:"kubernetes.io/pod.name"`
    PodNamespace string `json:"kubernetes.io/pod.namespace"`
    PodUID string `json:"kubernentes.io/pod.uid"`
    ReadWrite string `json:"kubernetes.io/readwrite"`
    ServiceAccount string `json:"kubernetes.io/serviceAccount.name"`
    XenServerHost string `json:"Host"`
    XenServerUsername string `json:"Username"`
    XenServerPassword string `json:"Password"`
}

func main() {
    var command string
    var mountDir string
    var jsonOptions string

    if len(os.Args) > 1 {
        command = os.Args[1]
    }

    if len(os.Args) > 2 {
        mountDir = os.Args[2]
    }

    if len(os.Args) > 3 {
        jsonOptions = os.Args[3]
    }

    switch command {
    case "init":
        fmt.Println("{\"status\": \"Success\", \"capabilities\": {\"attach\": false}}")
        os.Exit(0)
    case "mount":
        mount(mountDir, jsonOptions)
    case "unmount":
        unmount(mountDir)
    default:
        fmt.Println("{\"status\": \"Not supported\"}")
        os.Exit(1)
    }
}

func mount(mountDir, jsonOptions string) {
    // Loads config
    byt := []byte(jsonOptions)
    options := new(jsonParams)
    if err := json.Unmarshal(byt, &options); err != nil {
        failure(err)
    }

    if options.FSType == "" {
        options.FSType = "ext4"
    }

    var mode xenapi.VbdMode

    switch options.ReadWrite {
    case "ro":
        mode = xenapi.VbdModeRO
    case "rw":
        mode = xenapi.VbdModeRW
    default:
        failure(errors.New("Unknown ReadWrite"))
    }

    xapi, session, err := xapiLogin(options.XenServerHost, options.XenServerUsername, options.XenServerPassword)
    if err != nil {
        failure(fmt.Errorf("Cound not login at XenServer, error: %s", err.Error()))
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
        failure(errors.New("No VBD devices are available anymore"))
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

    //options.VDIUUID = string(vdiUUID)

    debug("VBD.GetAllRecords")
    vbds, err := xapi.VBD.GetAllRecords(session)
    if err != nil {
        failure(err)
    }

    for ref, vbd := range vbds {
        if vbd.VDI == vdiUUID && vbd.CurrentlyAttached {
            debug("Attempting to safely detached VDI")
            time.Sleep(5 * time.Second)
            if err := detachVBD(ref, xapi, session); err != nil {
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
    if err := xapi.VBD.Plug(session, vbdUUID); err != nil {
        failure(err)
    }

    debug("VBD.GetDevice")
    device, err := xapi.VBD.GetDevice(session, vbdUUID)
    if err != nil {
        failure(err)
    }
    devicePath := fmt.Sprintf("/dev/%s", device)

    blkid, err := run("blkid", devicePath)
    if err != nil && !strings.Contains(err.Error(), "exit status 2") {
        failure(err)
    }

    if blkid == "" {
        if _, err := run("mkfs", "-t", options.FSType, devicePath); err != nil {
            failure(err)
        }
    }

    debug("os.MkdirAll")
    if err := os.MkdirAll(mountDir, 0755); err != nil {
        failure(err)
    }

    debug("syscall.Mount")
    if err := syscall.Mount(devicePath, mountDir, options.FSType, 0, ""); err != nil {
        failure(err)
    }



    success()
}

func unmount(mountDir string) {
    debug("syscall.Unmount")
    if err := syscall.Unmount(mountDir, 0); err != nil {
        failure(err)
    }

    success()
}

func detachVBD(vbd xenapi.VBDRef, xapi *xenapi.Client, session xenapi.SessionRef) error {
    debug("VBD.Unplug")
    if err := xapi.VBD.Unplug(session, vbd); err != nil {
        return err
    }

    debug("VBD.Destroy")
    return xapi.VBD.Destroy(session, vbd)
}

func getMAC() (string, error) {
    debug("net.Interfaces")
    interfaces, err := net.Interfaces()
    if err != nil {
        return "", err
    }

    var mac string
    for _, i := range interfaces {
        if i.Name == "eth0" && i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
            mac = i.HardwareAddr.String()
        }
    }

    if mac == "" {
        return "", errors.New("MAC address not found")
    }

    return mac, nil
}


func getVM(xapi *xenapi.Client, session xenapi.SessionRef) (xenapi.VMRef, error) {
    mac, err := getMAC()
    if err != nil {
        return "", err
    }

    debug("VIF.GetAllRecords")
    vifs, err := xapi.VIF.GetAllRecords(session)
    if err != nil {
        return "", err
    }

    var vm xenapi.VMRef
    for _, vif := range vifs {
        if vif.MAC == mac && vif.CurrentlyAttached {
            vm = vif.VM
        }
    }

    if vm == "" {
        return "", errors.New("Could not find VM with MAC")
    }

    return vm, nil
}


func xapiLogin(host, username, password string) (*xenapi.Client, xenapi.SessionRef, error) {
    xapi, err := xenapi.NewClient(fmt.Sprintf("https://%s", host), nil)
    if err != nil {
        return nil, "", err
    }

    session, err := xapi.Session.LoginWithPassword(username, password, "1.0", "spangenberg.io/xenserver")
    if err != nil {
        return nil, "", err
    }

    return xapi, session, nil
}

func xapiLogout(xapi *xenapi.Client, session xenapi.SessionRef) error {
    return xapi.Session.Logout(session)
}

