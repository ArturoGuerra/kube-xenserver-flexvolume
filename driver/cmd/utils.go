package main

import (
    "fmt"
    "os"
    "os/exec"
    "encoding/json"
    "errors"
)

const debugLogFile = "/tmp/xenserver-driver.log"

func debug(message string) {
    if _, err := os.Stat(debugLogFile); err == nil {
        f, _ := os.OpenFile(debugLogFile, os.O_APPEND|os.O_WRONLY, 0600)
        defer f.Close()
        f.WriteString(fmt.Sprintln(message))
    }
}


func reply(message driverReply) {
    jsonData, _ := json.Marshal(message)
    fmt.Println(string(jsonData))
}

func loadOptions(data string) *jsonParams {
    var options jsonParams
    if err := json.Unmarshal([]byte(data), &options); err != nil {
        failure(errors.New("Error parsing jsonOptions"))
    }

    return &options
}

func failure(err error) {
    debug(fmt.Sprintf("FAILURE - %s", err.Error()))

    output := driverReply{
        Status: "Failure",
        Message: err.Error(),
    }

    reply(output)
    os.Exit(1)
}

func success(message string) {
    output := driverReply{
        Status: "Success",
        Message: message,
    }

    reply(output)
    os.Exit(0)
}

func run(cmd string, args  ...string) (string, error) {
    debug(fmt.Sprintf("Running %s %s", cmd, args))

    out, err := exec.Command(cmd, args...).CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("Error running %s %v: %v, %s", cmd, args, err, out)
    }
    return string(out), nil
}
