package main

import (
	"context"
	"fpetkovski/prometheus-adapter/pkg/apiserver"
	"fpetkovski/prometheus-adapter/pkg/controllermanager"
	"fpetkovski/prometheus-adapter/pkg/custommetrics"
	"fpetkovski/prometheus-adapter/pkg/externalmetrics"
	"os"
	"sync"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"

	"github.com/spf13/pflag"

	"github.com/go-logr/logr"

	"k8s.io/klog/v2/klogr"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var logger logr.Logger

func main() {
	logger = klogr.New()
	controllerruntime.SetLogger(logger)

	cfg := config.GetConfigOrDie()
	mgr, err := controllermanager.New(cfg)
	if err != nil {
		exit(err)
	}

	err = externalmetrics.RegisterController(mgr)
	if err != nil {
		exit(err)
	}
	metricsLister, err := custommetrics.RegisterController(cfg, mgr, logger.WithName("custom-metrics-controller"))
	if err != nil {
		exit(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := mgr.Start(context.Background()); err != nil {
			exit(err)
		}
	}()

	promCfg := api.Config{
		Address: "http://prometheus-operated.openshift-monitoring:9090",
	}
	promClient, err := api.NewClient(promCfg)
	if err != nil {
		exit(err)
	}

	promApi := v1.NewAPI(promClient)
	customMetricsProvider, err := custommetrics.NewMetricsProvider(
		mgr.GetConfig(),
		mgr.GetClient(),
		promApi,
		metricsLister,
		logger.WithName("custom-metrics"),
	)
	if err != nil {
		exit(err)
	}

	externalMetricsProvider := externalmetrics.NewMetricsProvider(
		mgr.GetClient(),
		promApi,
		logger.WithName("external-metrics"),
	)

	apiServer := apiserver.NewMetricsAPIServer("", customMetricsProvider, externalMetricsProvider)
	pflag.Parse()
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

func exit(err error) {
	logger.Error(err, "Terminating")
	os.Exit(1)
}
