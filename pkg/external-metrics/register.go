package external_metrics

import (
	"fpetkovski/prometheus-adapter/pkg/metrics-server"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func Register(
	server *metrics_server.MetricsAPIServer,
	controllerManager manager.Manager,
	prometheusAPI PromClient,
	logger logr.Logger,
) error {
	if err := registerController(controllerManager); err != nil {
		return err
	}

	externalMetricsProvider := NewMetricsProvider(
		controllerManager.GetClient(),
		prometheusAPI,
		logger,
	)
	server.WithExternalMetrics(externalMetricsProvider)

	return nil
}
