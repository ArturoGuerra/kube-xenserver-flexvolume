package flexvolume

import (
    "os"
    "fmt"
    "encoding/json"
    "github.com/arturoguerra/kube-xenserver-flexvolume/driver/pkg/utils"
)

type (
    DriverCapabilities struct {
        Attach bool `json:"attach"`
    }

    DriverReply struct {
        Status       string              `json:"status,omitempty"`
        Message      string              `json:"message,omitempty"`
        Device       string              `json:"device,omitempty"`
        VolumeName   string              `json:"volumeName,omitempty"`
        Attached     bool                `json:"attached,omitempty"`
        Capabilities *DriverCapabilities `json:"capabilities,omitempty"`
    }
)

func Success(msg string) *DriverReply {
    return &DriverReply{
        Status: "Success",
        Message: msg,
    }
}

func Failure(msg string) *DriverReply {
    return &DriverReply{
        Status: "Failure",
        Message: msg,
    }
}

func (n *nodeClient) Reply(message *DriverReply) {
    var exitCode int
    switch message.Status {
    case "Success":
        exitCode = 0
    case "Failure":
        exitCode = 1
    case "Not supported":
        exitCode = 1
    default:
        exitCode = 1
    }

    bdata, _ := json.Marshal(message)
    data := string(bdata)
    fmt.Println(data)
    utils.Debug(fmt.Sprintf("Data: %s Code: %d", data, exitCode))

    os.Exit(exitCode)
}
