package flexvolume

import (
    "syscall"
    "github.com/arturoguerra/kube-xenserver-flexvolume/driver/pkg/utils"
)

func (n *nodeClient) MountDevice(mountdir, devicename string, p *JsonParams) {
    utils.Debug("syscall.Mount")
    if err := syscall.Mount(devicename, mountdir, "auto", 0, ""); err != nil {
        n.Reply(Failure(err.Error()))
    }

    n.Reply(Success("Mounted Device"))
}
