package main

import (
    "fmt"
    xenapi "github.com/terra-farm/go-xen-api-client"
)

func xapiLogin(params *jsonParams) (*xenapi.Client, xenapi.SessionRef, error) {
    xapi, err := xenapi.NewClient(fmt.Sprintf("https://%s", params.XenServerHost), nil)
    if err != nil {
        return nil, "", err
    }

    session, err := xapi.Session.LoginWithPassword(params.XenServerUsername, params.XenServerPassword, "1.0", "arturoguerra/xenserver")
    if err != nil {
        return nil, "", err
    }

    return xapi, session, nil
}

func xapiLogout(xapi *xenapi.Client, session xenapi.SessionRef) {
    if err := xapi.Session.Logout(session); err != nil {
        failure(fmt.Errorf("Failed to logout from XenServer, Error: %s", err.Error()))
    }
}
