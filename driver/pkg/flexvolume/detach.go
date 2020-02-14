package flexvolume

import (
    "fmt"
    "github.com/arturoguerra/kube-xenserver-flexvolume/driver/pkg/utils"
)

func (n *masterClient) Detach(volname, nodename string) {
    if err := n.API.Detach(volname, nodename); err != nil {
        n.Reply(&DriverReply{
            Status: "Failure",
            Message: err.Error(),
        })
    }

    msg := fmt.Sprintf("Detached %s from %s", volname, nodename)
    utils.Debug(msg)
    n.Reply(&DriverReply{
        Status: "Success",
        Message: msg,
    })
}

func (n *nodeClient) Detach(vol, name string) {
    n.NotSupported()
}
