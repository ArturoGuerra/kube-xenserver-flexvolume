package xapi

func (c *xClient) Attach(label, mode, fstype, volumename string) (string, error) {
    if !c.Master {
        return "", errors.New("Not master")
    }

    if fstype == "" {
        fstype = "ext4"
    }

    var xmode xenapi.VbdMode
    switch mode {
    case "ro":
        mode = xenapi.VbdModeRO
    case "rw":
        mode = xenapi.VbdModeRW
    default:
        return "", errors.New("Unkown ReadWrite Mode")
    }

    api, session, err := c.Connect()
    if err != nil {
        return "", err
    }

    defer c.Close(api, session)

    vm, err := c.GetVM(label)
    if err != nil {
        return "", err
    }

    utils.Debug("VM.GetAllAllowedVBDDevices")
    vbdDevices, err := api.VM.GetAllowedVBDDevices(session, vm)
    if err != nil {
        return "", err
    }

    if len(vbdDevices) < 0 {
        return "", errors.New("No VBD Devices are available")
    }

    utils.Debug("VDI.GetAllRecords")
    vdis, err := api.VDI.GetAllRecords(session)
    if err != nil {
        return "", err
    }

    var vdiUUID xenapi.VDIRef
    for ref, vdi := range vdis {
        if vdi.NameLabel == volumename && !vdi.IsSnapshot {
            vdiUUID = ref
        }
    }

    if string(vdiUUID) == "" {
        return "", errors.New("Count not find VDI")
    }

    utils.Debug("VBD.GetAllRecords")
    vbds, err := api.VBD.GetAllRecords(session)
    if err != nil {
        return "", err
    }

    for ref, vbd := range vbds {
        if vbd.VDI == vdiUUID && vbd.CurrentlyAttached {
            utils.Debug("Attempting to safely detach VDI")
            time.Sleep(10 * time.Second)
            if err := f.DetachVBD(ref, api, session); err != nil {
                return "", err
            }
        }
    }

    utils.Debug("VBD.Create")
    vbdUUID, err := api.VBD.Create(session, xenapi.VBDRecord{
        Bootable:    false,
        Mode:        xmode,
        Type:        xenapi.VbdTypeDisk,
        Unpluggable: true,
        Userdevice:  vbdDevices[0],
        VDI:         vdiUUID,
        VM:          vm,
    })
    if err != nil {
        return "", err
    }

    utils.Debug("VBD.Plug")
    if err != api.VBD.Plug(session, vbdUUID); err != nil {
        return "", err
    }

    utils.Debug("VBD.GetDevice")
    device, err := api.VBD.GetDevice(session, vbdUUID)
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("/dev/%s", device), nil
}
