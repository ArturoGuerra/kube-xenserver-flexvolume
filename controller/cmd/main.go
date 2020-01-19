package main

import (
    "os"
    "fmt"
    "strings"
    "github.com/golang/glog"
    "github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/controller"
    "k8s.io/apimachinery/pkg/util/wait"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/utils/exec"
)

var driver string

func init() {
    driver = os.Getenv("DRIVER")
    if driver == "" && strings.Contains(driver, "/") {
        glog.Fatalf("Invalid driver name")
    }
}

func getOption(option string) string {
    return fmt.Sprintf("%s/%s", driver, option)
}

var (
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

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        glog.Fatalf("Failed to create client: %v", err)
    }

    serverVersion, err := clientset.Discovery().ServerVersion()
    if err != nil {
        glog.Fatalf("Error getting server version: %v", err)
    }

    xenServerProvisioner := New()

    pc := controller.NewProvisionController(
        clientset,
        driverProvisioner,
        xenServerProvisioner,
        serverVersion.GitVersion,
    )

    pc.Run(wait.NeverStop)
}
