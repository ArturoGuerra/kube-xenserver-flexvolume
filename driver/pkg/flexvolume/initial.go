package flexvolume

func (n *nodeClient) Init() {
    c := &DriverCapabilities{
        Attach: true,
    }

    o := &DriverReply{
        Status: "Success",
        Capabilities: c,
    }

    n.Reply(o)
}
