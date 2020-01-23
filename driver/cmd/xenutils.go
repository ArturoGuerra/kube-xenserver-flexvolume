package main

import (
    "net"
    "bytes"
    "errors"
    xenapi "github.com/terra-farm/go-xen-api-client"
)

func detachVBD(vbd xenapi.VBDRef, xapi *xenapi.Client, session xenapi.SessionRef) error {
    debug("VBD.Unplug")
    if err := xapi.VBD.Unplug(session, vbd); err != nil {
        return err
    }

    debug("VBD.Destroy")
    return xapi.VBD.Destroy(session, vbd)
}

func forceDetachVBD(vbd xenapi.VBDRef, xapi *xenapi.Client, session xenapi.SessionRef) error {
    debug("VBD.Unplug")
    if err := xapi.VBD.Unplug(session, vbd); err != nil {
        debug("VBD.UnplugForce")
        if err := xapi.VBD.UnplugForce(session, vbd); err != nil {
            return err
        }
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
        return "", errors.New("Count not find VM with MAC")
    }

    return vm, nil
}
