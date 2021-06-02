module fpetkovski/prometheus-adapter

go 1.16

require (
	github.com/go-logr/logr v0.4.0
	github.com/kubernetes-sigs/custom-metrics-apiserver v0.0.0-20210311094424-0ca2b1909cdc
	github.com/prometheus/client_golang v1.10.0
	github.com/prometheus/common v0.18.0
	github.com/spf13/pflag v1.0.5
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/apiserver v0.20.1
	k8s.io/client-go v0.20.2
	k8s.io/klog/v2 v2.8.0
	k8s.io/metrics v0.20.0
	sigs.k8s.io/controller-runtime v0.8.3
)

replace github.com/googleapis/gnostic v0.5.1 => github.com/googleapis/gnostic v0.4.1
