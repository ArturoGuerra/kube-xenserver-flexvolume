package flexvolume

import (
    "github.com/arturoguerra/kube-xenserver-flexvolume/driver/pkg/xapi"
)

type (
    JsonParams struct {
        FSGroup           string `json:"kubernetes.io/fsGroup"`
        FSType            string `json:"kubernetes.io/fsType"`
        PVOrVolumeName    string `json:"kubernetes.io/pvOrVolumeName"`
        PodName           string `json:"pod.name"`
        PodNamespace      string `json:"pod.namespace"`
        PodUID            string `json:"kubernetes.io/pod.uid"`
        ReadWrite         string `json:"kubernetes.io/readwrite"`
        ServiceAccount    string `json:"kubernetes.io/ServiceAccount.name"`
        VDIUUID           string `json:"VDIUUID"`
    }


    nodeClient struct {
        Master bool
    }

    masterClient struct {
        *nodeClient
        API xapi.XClient
    }

    FlexVolume interface {
        Init()
        GetVolumeName(*JsonParams) // Json Parameters
        Attach(*JsonParams, string) // Json Parameters, NodeName //  Master Only
        WaitForAttach(string, *JsonParams) // DeviceName, Json Parameters
        Detach(string, string) // VolumeName, NodeName // Master Only
        IsAttached(*JsonParams, string) // Json Parameters, NodeName
        MountDevice(string, string, *JsonParams) // MountDir, DeviceName, Json Parameters
        UnmountDevice(string) // MountDir
        Options(string) *JsonParams
        NotSupported()
    }
)

func NewMaster(username, password, host string) FlexVolume {
    node := &nodeClient{true}
    api := xapi.New(username, password, host)
    master := &masterClient{node,api}
    return master
//    return &flexVolume{
//        Master:   true,
//        Xapi:     xapi.New(username, password, host),
//    }
}

func NewNode() FlexVolume {
    return &nodeClient{true}
}
