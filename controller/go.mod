module github.com/arturoguerra/kube-xenserver-flexvolume/controller

go 1.13

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/kubernetes-sigs/sig-storage-lib-external-provisioner v4.0.1+incompatible
	github.com/miekg/dns v1.1.27 // indirect
	github.com/prometheus/client_golang v1.3.0 // indirect
	github.com/terra-farm/go-xen-api-client v0.0.0-20191130210227-94f14387b8c2
	k8s.io/api v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0
	k8s.io/utils v0.0.0-20191114184206-e782cd3c129f
	sigs.k8s.io/sig-storage-lib-external-provisioner v4.0.1+incompatible // indirect
)
