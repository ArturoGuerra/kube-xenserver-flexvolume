package main

import (
    "github.com/golang/glog"
    "github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/controller"
    "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (p *XenServerProvisioner) Provision(options controller.ProvisionOptions) (*v1.PersistentVolume, error) {
    glog.Infof("Provision called for volume: %s", options.PVName)

    if err := p.ProvisionOnXenServer(options); err != nil {
        glog.Errorf("Failed to provision volume %s error %s", options, err.Error())
        return nil, err
    }

    pv := &v1.PersistentVolume{
        ObjectMeta: metav1.ObjectMeta{
            Name: options.PVName,
        },
        Spec: v1.PersistentVolumeSpec{
            AccessModes: options.PVC.Spec.AccessModes,
            Capacity: v1.ResourceList{
                v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
            },
            StorageClassName: options.StorageClass.ObjectMeta.Name,
            PersistentVolumeReclaimPolicy: *options.StorageClass.ReclaimPolicy,
            PersistentVolumeSource: v1.PersistentVolumeSource{
                FlexVolume: &v1.FlexPersistentVolumeSource{
                    Driver: driver,
                    FSType: "ext4",
                    Options: map[string]string{
                        Host:     options.StorageClass.Parameters[Host],
                        srName:   options.StorageClass.Parameters[srName],
                        Username: options.StorageClass.Parameters[Username],
                        Password: options.StorageClass.Parameters[Password],
                    },
                },
            },
        },
    }

    return pv, nil
}

func (p *XenServerProvisioner) Delete(volume *v1.PersistentVolume) error {
    host := volume.Spec.FlexVolume.Options[Host]
    username := volume.Spec.FlexVolume.Options[Username]
    password := volume.Spec.FlexVolume.Options[Password]

    glog.Infof("Delete called for volume: %s", volume.Name)

    if err := p.DeleteFromXenServer(host, username, password, volume.ObjectMeta.Name); err != nil {
        glog.Errorf("Failed to delete volume %s error: %s", volume, err.Error())
    }

    return nil
}
