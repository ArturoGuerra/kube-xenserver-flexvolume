package main

import (
    xenapi "github.com/terra-farm/go-xen-api-client"
)

func xapiLogin(options jsonParams) (*xenapi.Client, xenapi.SessionRef, error) {
    xapi, err := xenapi.NewClient(fmt.Sprintf("https://%s", host), nil)
    if err != nil {
        return nil, "", err
    }

    session, err := xapi.Session.LoginWithPassword(options.XenServerHost, options.XenServerPassword, "1.0", driver)
    if err != nil {
        return nil, "", err
    }

    return xapi, session, nil
}

func xapiLogout(xapi, *xenapi.Client, session xenapi.SessionRef) error {
    return xapi.Session.Logout(session)
}
