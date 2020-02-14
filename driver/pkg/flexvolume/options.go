package flexvolume

import "encoding/json"

func (n *nodeClient) Options(raw string) *JsonParams {
    var opts JsonParams
    if err := json.Unmarshal([]byte(raw), &opts); err != nil {
        n.Reply(&DriverReply{
            Status: "Failure",
            Message: err.Error(),
        })
    }

    return &opts
}
