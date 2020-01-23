package main

import (
    "syscall"
    "fmt"
    "io/ioutil"
    "encoding/json"
    "os"
    "strings"
    "errors"
)

func unmount(mountDir string) {
    var err error

    debug("syscall.Unmount")
    if err = syscall.Unmount(mountDir, 0); err != nil {
        if err.Error() != "invalid argument" {
            failure(err)
        }
    }

    params := loadParamsFromFile(mountDir)

    defer deleteParamsFile(mountDir)

    detach(params, mountDir)

    success("Unmounted and detached volume")
}

func loadParamsFromFile(mountDir string) *jsonParams {
    debug("ioutil.ReadFile")
    jsonOptionsFile := fmt.Sprintf("%s.json", mountDir)
    byt, err := ioutil.ReadFile(jsonOptionsFile)
    if err != nil {
        resp := "Unable to read config file which means it was already unmounted"
        debug(resp)
        success(resp)
    }

    options := new(jsonParams)
    if err := json.Unmarshal(byt, options); err != nil {
        failure(err)
    }

    return options
}

func deleteParamsFile(mountDir string) {
    debug("os.Remove")
    if err := os.Remove(mountDir); err != nil {
        failure(err)
    }


}

func detach(params *jsonParams, mountDir string) {
    xapi, session, err := xapiLogin(params)
    if err != nil {
        failure(err)
    }

    defer xapiLogout(xapi, session)

    vm, err := getVM(xapi, session)
    if err != nil {
        failure(err)
    }

    devicePath, err := run("findmnt", "-n", "-o", "SOURCE", "--target", mountDir)
    if err != nil {
        failure(err)
    }

    devicePathElements := strings.Split(devicePath, "/")
    if len(devicePathElements) < 3 || len(devicePathElements) > 3 {
        failure(errors.New("Device path is incorrect"))
    }

    debug("VBD.GetAllRecords")
    vbds, err := xapi.VBD.GetAllRecords(session)
    if err != nil {
        failure(err)
    }
    debug("Detaching VDI")
    for ref, vbd := range vbds {
        if vbd.VM == vm && string(vbd.VDI) == params.VDIUUID && vbd.CurrentlyAttached {
            if err := forceDetachVBD(ref, xapi, session); err != nil {
                failure(err)
            }
        }
    }
}
