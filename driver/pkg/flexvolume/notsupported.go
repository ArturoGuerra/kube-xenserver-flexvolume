package flexvolume

func (n *nodeClient) NotSupported() {
    reply := &DriverReply{
        Status: "Not supported",
    }

    n.Reply(reply)
}
