module github.com/zorrofox/csi-goofys-s3

require (
	github.com/Azure/azure-pipeline-go v0.2.1
	github.com/Azure/azure-sdk-for-go v32.1.0+incompatible
	github.com/Azure/azure-storage-blob-go v0.7.0
	github.com/Azure/go-autorest v13.0.2+incompatible
	github.com/Azure/go-autorest/autorest v0.10.0 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.2 // indirect
	github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/aws/aws-sdk-go v1.29.16
	github.com/container-storage-interface/spec v1.1.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.4 // indirect
	github.com/jacobsa/fuse v0.0.0-00010101000000-000000000000 // indirect
	github.com/kahing/goofys v0.23.1
	github.com/kubernetes-csi/csi-lib-utils v0.7.0 // indirect
	github.com/kubernetes-csi/csi-test v2.2.0+incompatible
	github.com/kubernetes-csi/drivers v1.0.2
	github.com/mitchellh/go-ps v0.0.0-20170309133038-4fdf99ab2936
	github.com/onsi/ginkgo v1.10.2
	github.com/onsi/gomega v1.7.0
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/shirou/gopsutil v2.20.2+incompatible // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/urfave/cli v1.22.2 // indirect
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	google.golang.org/genproto v0.0.0-20200303153909-beee998c1893 // indirect
	google.golang.org/grpc v1.27.0
	gopkg.in/ini.v1 v1.41.0
	k8s.io/apimachinery v0.17.3 // indirect
	k8s.io/kubernetes v1.13.4
	k8s.io/utils v0.0.0-20200229041039-0a110f9eb7ab // indirect
)

replace github.com/jacobsa/fuse => github.com/kahing/fusego v0.0.0-20190829004836-cb66eea2377f
