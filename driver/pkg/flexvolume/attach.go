package flexvolume

func (n *masterClient) Attach(p *JsonParams, nodename string) {
    device, err := n.API.Attach(nodename, p.ReadWrite, p.FSType, p.PVOrVolumeName)
    if err != nil {
        n.Reply(&DriverReply{
            Status: "Failure",
            Message: err.Error(),
            Device: "",
        })
    }

    n.Reply(&DriverReply{
        Status: "Success",
        Message: "Attached",
        Device: device,
    })
}

func (n *nodeClient)Attach(p *JsonParams, nodename string) {
    n.NotSupported()
}
