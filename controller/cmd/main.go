package main

import (
    "strings"
    glog "github.com/golang/glog"
    controller "github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/controller"
    wait "k8s.io/apimachinery/pkg/wait"
    kubernetes "k8s.io/client-go/kubernetes"
    rest "k8s.io/client-go/rest"
)

var driverName string

func init() {
    driverName = os.Getenv("DIVIER_NAME")
    if driverName == "" && strings.Contains(driverName, "/") {
        glog.Fatalf("Invalid driver name")
    }
}

func getOption(option string) string {
    return fmt.Sprintf("%s/%s", driverPrefix, option)
}

const (
    driver                        = driverPrefix
    driverOptionXenServerHost     = getOption("host")
    driverOptionXenServerUsername = getOption("username")
    driverOptionXenServerPassword = getOption("password")
    StorageClassParameterSRName   = getOption("srName")
    driverProvisioner             = getOption("provisioner")
    driverFSType                  = "ext4"
)

type XenServerProvisioner struct {
    runner            exec.Interface
    XenServerHost     string
    XenServerUsername string
    XenServerPassword string
}

func New() controller.Provisioner {
    return &XenServerProvisioner{
        runner:            exec.New(),
        XenServerHost:     os.Getenv("XENSERVER_HOST"),
        XenServerUsername: os.Getenv("XENSERVER_USERNAME"),
        XenServerPassword: os.Getenv("XENSERVER_PASSWORD"),
    }
}

func main() {

    config, err := rest.InClusterConfig()
    if err != nil {
        glog.Fatalf("Failed to create config: %v", err)
    }

    client, err := kubernetes.NewForConfig(config)
    if err != nil {
        glog.Fatalf("Failed to create client: %v", err)
    }

    serverVersion, err := client.Discovery().ServerVersion()
    if err != nil {
        glog.Fatalf("Error getting server version: %v", err)
    }

    xenServerProvisioner := New()

    pc := controller.NewProvisionerController(
        client,
        *driverProvisioner,
        xenServerProvisioner,
        serverVersion.GitVersion,
    )

    pc.Run(wait.NeverStop)
}
