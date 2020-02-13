package main

import (
    "fmt"
    "os"
)


type (
    jsonParams struct {
        FSGroup           string `json:"kubernetes.io/fsGroup"`
        FSType            string `json:"kubernetes.io/fsType"`
        PVOrVolumeName    string `json:"kubernetes.io/pvOrVolumeName"`
        PodName           string `json:"kubernetes.io/pod.name"`
        PodNamespace      string `json:"kubernetes.io/pod.namespace"`
        PodUID            string `json:"kubernentes.io/pod.uid"`
        ReadWrite         string `json:"kubernetes.io/readwrite"`
        ServiceAccount    string `json:"kubernetes.io/serviceAccount.name"`
        XenServerHost     string `json:"Host"`
        XenServerUsername string `json:"Username"`
        XenServerPassword string `json:"Password"`
        VDIUUID           string `json:"VDIUUID"`
    }

    driverCapabilities struct {
        Attach bool `json:"attach,omiempty"`
    }

    driverReply struct {
        Status       string `json:"status,omitempty"`
        Message      string `json:"message,omitempty"`
        Device       string `json:"device,omitempty"`
        VolumeName   string `json:"volumeName,omitempty"`
        Attached     string `json:"attached,omitempty"`
        Capabilities driverCapabilities `json:"capabilities,omitempty"`
    }
)

func main() {
    var command string

    if len(os.Args) > 1 {
        command = os.Args[1]
    }

    debug(fmt.Sprintf("%s", command))

    switch command {
    case "init":
        initCommand()

    case "mount":
        mountDir := os.Args[2]
        jsonOptions := loadOptions(os.Args[3])
        mount(mountDir, jsonOptions)

    case "unmount":
        mountDir := os.Args[2]
        unmount(mountDir)

    default:
        notSupported()
    }
}
