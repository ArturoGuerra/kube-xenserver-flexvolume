package main

import (
    "fmt"
    "os"
    "os/exec"
    "encoding/json"
)

const debugLogFile = "/tmp/xenserver-driver.log"

func debug(message string) {
    if _, err := os.Stat(debugLogFile); err == nil {
        f, _ := os.OpenFile(debugLogFile, os.O_APPEND|os.O_WRONLY, 0600)
        defer f.Close()
        f.WriteString(fmt.Sprintln(message))
    }
}


func success() {
    debug("SUCCESS")

    fmt.Print("{\"status\": \"Success\"}")
    os.Exit(0)
}

func failure(err error) {
    debug(fmt.Sprintf("FAILURE - %s", err.Error()))

    failureMap := map[string]string{"status": "Failure", "message": err.Error()}
    jsonMessage, _ := json.Marshal(failureMap)
    fmt.Print(string(jsonMessage))

    os.Exit(1)
}

func run(cmd string, args  ...string) (string, error) {
    debug(fmt.Sprintf("Running %s %s", cmd, args))

    out, err := exec.Command(cmd, args...).CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("Error running %s %v: %v, %s", cmd, args, err, out)
    }
    return string(out), nil
}
