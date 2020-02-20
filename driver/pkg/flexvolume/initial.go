package flexvolume

func (n *nodeClient) Init() {
    o := &DriverReply{
        Status: "Success",
        Message: "Plz attach",
    }

    n.Reply(o)
}
