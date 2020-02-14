package flexvolume

import (
    "fmt"
    "github.com/arturoguerra/kube-xenserver-flexvolume/driver/pkg/utils"
)

func (n *nodeClient) GetVolumeName(p *JsonParams) {
    utils.Debug(fmt.Sprintf("Sending volumename: %s", p.PVOrVolumeName))
    reply := &DriverReply{
        Status: "Success",
        Message: "returning device name",
        VolumeName: p.PVOrVolumeName,
    }

    n.Reply(reply)
}
