package main

import (
	"context"
	"fpetkovski/prometheus-adapter/pkg/apiserver"
	"fpetkovski/prometheus-adapter/pkg/controllermanager"
	"fpetkovski/prometheus-adapter/pkg/externalmetrics"
	"io/fs"
	"io/ioutil"
	"k8s.io/client-go/util/cert"
	"net"
	"os"
	"sync"

	"github.com/spf13/pflag"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"

	"github.com/go-logr/logr"

	"k8s.io/klog/v2/klogr"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var logger logr.Logger

func main() {
	logger = klogr.New()
	controllerruntime.SetLogger(logger)
	generateSelfSignedCertificate()

	cfg := config.GetConfigOrDie()
	mgr, err := controllermanager.New(cfg)
	if err != nil {
		exit(err)
	}

	prometheusAPI := makePrometheusAPI("http://prometheus-operated.openshift-monitoring:9090")
	apiServer := apiserver.NewMetricsAPIServer()
	if err = externalmetrics.Register(apiServer, mgr, prometheusAPI, logger.WithName("external-metrics")); err != nil {
		exit(err)
	}

	pflag.Parse()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := mgr.Start(context.Background()); err != nil {
			exit(err)
		}
	}()

	wg.Add(1)
	go func() {
		done := make(chan struct{})
		if err := apiServer.Run(done); err != nil {
			exit(err)
		}
		defer wg.Done()
	}()

	wg.Wait()
}

func makePrometheusAPI(prometheusUrl string) v1.API {
	promCfg := api.Config{
		Address: prometheusUrl,
	}
	promClient, err := api.NewClient(promCfg)
	if err != nil {
		exit(err)
	}
	promApi := v1.NewAPI(promClient)
	return promApi
}

func generateSelfSignedCertificate() {
	crt, key, err := cert.GenerateSelfSignedCertKey(
		"custom-metrics.custom-metrics",
		[]net.IP{},
		[]string{"custom-metrics.custom-metrics.svc"},
	)
	if err != nil {
		exit(err)
	}

	if err := ioutil.WriteFile("/var/run/tls/tls.key", key, fs.FileMode(0644)); err != nil {
		exit(err)
	}

	if err := ioutil.WriteFile("/var/run/tls/tls.crt", crt, fs.FileMode(0644)); err != nil {
		exit(err)
	}
}

func exit(err error) {
	logger.Error(err, "Terminating")
	os.Exit(1)
}
