package main

import (
    "os"
    "fmt"
    "errors"
    "github.com/arturoguerra/kube-xenserver-flexvolume/driver/pkg/flexvolume"
    "github.com/arturoguerra/kube-xenserver-flexvolume/driver/pkg/utils"
)

var (
    master   bool
    client   flexvolume.FlexVolume
)

func init() {
    cfg := utils.Init()
    if cfg.Master {
        utils.Debug("Running as master")
        client = flexvolume.NewMaster(cfg.Username, cfg.Password, cfg.Host)
    } else {
        utils.Debug("Running as node")
        client = flexvolume.NewNode()
    }
}

func main() {
    if len(os.Args) == 0 {
        panic(errors.New("Invalid ammount of args"))
    }

    command := os.Args[1]
    utils.Debug(fmt.Sprintf("Running: %s", command))
    args := make([]string, 0)
    if len(os.Args) > 2 {
        args = append(args, os.Args[2:]...)
    }

    var _ = args

    switch command {
    case "init":
        client.Init()
    case "getvolumename":
        client.GetVolumeName(client.Options(args[0]))
    case "attach":
        client.Attach(client.Options(args[0]), args[1])
    case "waitforattach":
        client.WaitForAttach(args[0], client.Options(args[1]))
    case "detach":
        client.Detach(args[0], args[1])
    case "isattached":
        client.IsAttached(client.Options(args[0]), args[1])
    case "mountdevice":
        client.MountDevice(args[0], args[1], client.Options(args[2]))
    case "unmountdevice":
        client.UnmountDevice(args[0])
    default:
        client.NotSupported()
    }
}
