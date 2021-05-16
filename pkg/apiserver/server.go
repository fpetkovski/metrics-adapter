package apiserver

import (
	"github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/provider"

	"github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/apiserver"
	"github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/cmd"
	generatedopenapi "github.com/kubernetes-sigs/custom-metrics-apiserver/test-adapter/generated/openapi"
	openapinamer "k8s.io/apiserver/pkg/endpoints/openapi"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

type MetricsAPIServer struct {
	cmd.AdapterBase

	// PrometheusURL is to the Prometheus server to query.
	// Query parameters can be used to configure connection options.
	PrometheusURL string
}

func NewMetricsAPIServer(
	prometheusUrl string,
	customMetricsProvider provider.CustomMetricsProvider,
	externalMetricsProvider provider.ExternalMetricsProvider,
) *MetricsAPIServer {
	server := &MetricsAPIServer{
		AdapterBase:   cmd.AdapterBase{},
		PrometheusURL: prometheusUrl,
	}

	server.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(generatedopenapi.GetOpenAPIDefinitions, openapinamer.NewDefinitionNamer(apiserver.Scheme))
	server.OpenAPIConfig.Info.Title = "prometheus-adapter"
	server.OpenAPIConfig.Info.Version = "1.0.0"
	server.InstallFlags()
	server.WithCustomMetrics(customMetricsProvider)
	server.WithExternalMetrics(externalMetricsProvider)
	
	return server
}
