package main

import (
    "fmt"
    "errors"
    xenapi "github.com/terra-farm/go-xen-api-client"
    "github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/controller"
    "github.com/golang/glog"
    "k8s.io/api/core/v1"
)

// XAPI Authentication
func (p *XenServerProvisioner) XapiLogin() (*xenapi.Client, xenapi.SessionRef, error) {
    xapi, err := xenapi.NewClient(fmt.Sprintf("https://%s", p.XenServerHost), nil)
	if err != nil {
		return nil, "", err
	}

	session, err := xapi.Session.LoginWithPassword(p.XenServerUsername, p.XenServerPassword, "1.0", provisioner)
	if err != nil {
		return nil, "", err
	}

	return xapi, session, nil
}

func (p *XenServerProvisioner) XapiLogout(client *xenapi.Client, session xenapi.SessionRef) error {
    return  client.Session.Logout(session)
}


// XAPI Functions
func (p *XenServerProvisioner) ProvisionOnXenServer(options controller.ProvisionOptions) error {
    xapi, session, err := p.XapiLogin()
	if err != nil {
		return fmt.Errorf("Could not login at XenServer, error: %s", err.Error())
	}
	defer func() {
		if err := p.XapiLogout(xapi, session); err != nil {
			glog.Errorf("Failed to log out from XenServer, error: %s", err.Error())
		}
	}()

	srNameLabel := options.StorageClass.Parameters[srName]
	srs, err := xapi.SR.GetByNameLabel(session, srNameLabel)
	if err != nil {
		return fmt.Errorf("Could not list SRs for name label %s, error: %s", srNameLabel, err.Error())
	}

	if len(srs) > 1 {
		return fmt.Errorf("Too many SRs where found for name label %s", srNameLabel)
	}

	if len(srs) < 1 {
		return fmt.Errorf("No SR was found for name label %s", srNameLabel)
	}

	capacity, exists := options.PVC.Spec.Resources.Requests[v1.ResourceStorage]
	if !exists {
		return fmt.Errorf("Capacity was not specified for name label %s", options.PVName)
	}

	_, err = xapi.VDI.Create(session, xenapi.VDIRecord{
		NameDescription: "Kubernetes Persisted Volume Claim",
		NameLabel:       options.PVName,
		SR:              srs[0],
		Type:            xenapi.VdiTypeUser,
		VirtualSize:     int(capacity.Value()),
	})
	if err != nil {
		return fmt.Errorf("Could not create VDI for name label %s, error: %s", options.PVName, err.Error())
	}

	return nil
}

func (p *XenServerProvisioner) DeleteFromXenServer(nameLabel string) error {
    xapi, session, err := p.XapiLogin()
    if err != nil {
		return errors.New(fmt.Sprintf("Could not login at XenServer, error: %s", err.Error()))
	}
	defer func() {
		if err := p.XapiLogout(xapi, session); err != nil {
			glog.Errorf("Failed to log out from XenServer, error: %s", err.Error())
		}
	}()

	vdis, err := xapi.VDI.GetByNameLabel(session, nameLabel)
	if err != nil {
		return fmt.Errorf("Could not list VDIs for name label %s, error: %s", nameLabel, err.Error())
	}

	if len(vdis) > 1 {
		return fmt.Errorf("Too many VDIs where found for name label %s", nameLabel)
	}

	if len(vdis) > 0 {
        vbds, err := xapi.VBD.GetAllRecords(session)
        if err != nil {
            return fmt.Errorf("Error getting all VBDs error: %s", err.Error())
        }

        for ref, vbd := range vbds {
            if vbd.VDI == vdis[0] && vbd.CurrentlyAttached {
                if err = xapi.VBD.Unplug(session, ref); err != nil {
                    return fmt.Errorf("Error unpluging VBD error: %s", err.Error())
                }

                if err = xapi.VBD.Destroy(session, ref); err != nil {
                    return fmt.Errorf("Error destroying VBD error: %s", err.Error())
                }
            }
        }

		err := xapi.VDI.Destroy(session, vdis[0])
		if err != nil {
			return fmt.Errorf("Could not destroy VDI for name label %s, error: %s", nameLabel, err.Error())
		}

		glog.Infof("VDI was destroyed for name label %s", nameLabel)
	} else {
		glog.Infof("VDI was already destroyed for name label %s", nameLabel)
	}

	return nil
}

