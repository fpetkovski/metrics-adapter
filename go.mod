module fpetkovski/prometheus-adapter

go 1.16

require (
	github.com/armon/go-metrics v0.3.3 // indirect
	github.com/go-logr/logr v0.4.0
	github.com/hashicorp/go-hclog v0.12.2 // indirect
	github.com/hashicorp/go-immutable-radix v1.2.0 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/kubernetes-sigs/custom-metrics-apiserver v0.0.0-20210311094424-0ca2b1909cdc
	github.com/prometheus/client_golang v1.10.0
	github.com/prometheus/common v0.20.0
	github.com/prometheus/prometheus v1.8.2-0.20210331101223-3cafc58827d1
	github.com/samuel/go-zookeeper v0.0.0-20201211165307-7117e9ea2414 // indirect
	github.com/spf13/pflag v1.0.5
	k8s.io/apimachinery v0.20.5
	k8s.io/apiserver v0.20.5
	k8s.io/client-go v0.20.5
	k8s.io/klog/v2 v2.8.0
	k8s.io/metrics v0.20.0
	sigs.k8s.io/controller-runtime v0.8.3
)

replace (
	github.com/googleapis/gnostic v0.5.1 => github.com/googleapis/gnostic v0.4.1
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
