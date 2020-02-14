package flexvolume

import (
    "fmt"
    "github.com/arturoguerra/kube-xenserver-flexvolume/driver/pkg/utils"
)

func (n *nodeClient) WaitForAttach(devicename string, p *JsonParams) {
    utils.Debug(fmt.Sprintf("Drive %s is attached", devicename))
    reply := &DriverReply{
        Status: "Success",
        Device: devicename,
    }

    n.Reply(reply)
}
