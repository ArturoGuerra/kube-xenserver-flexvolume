package flexvolume

func (n *masterClient) IsAttached(p *JsonParams, nodename string) {
    reply := &DriverReply{
        Status: "Success",
        Message: "",
        Attached: false,
    }

    ok, err := n.API.IsAttached(p.PVOrVolumeName, nodename)
    if err != nil {
        reply.Status = "Failure"
    } else if ok {
        reply.Attached = true
    }

    n.Reply(reply)
}

func (n *nodeClient) IsAttached(p *JsonParams, nodename string) {
    n.NotSupported()
}
