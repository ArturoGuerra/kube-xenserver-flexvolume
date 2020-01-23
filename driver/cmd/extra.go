package main

import (
    "os"
    "fmt"
    "encoding/json"
)

func initCommand() {
    capabilities := driverCapabilities{
        Attach: false,
    }

    output := driverReply{
        Status: "Success",
        Capabilities: capabilities,
    }

    data, _ := json.Marshal(output)
    fmt.Println(string(data))
    os.Exit(0)
}

func notSupported() {
    output := driverReply{
        Status: "Not supported",
    }

    data, _ := json.Marshal(output)
    fmt.Println(string(data))
    os.Exit(1)
}
