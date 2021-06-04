package apiserver

import (
	"github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/apiserver"
	"github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/cmd"
	generatedopenapi "github.com/kubernetes-sigs/custom-metrics-apiserver/test-adapter/generated/openapi"
	openapinamer "k8s.io/apiserver/pkg/endpoints/openapi"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

type MetricsAPIServer struct {
	cmd.AdapterBase
}

func NewMetricsAPIServer() *MetricsAPIServer {
	server := &MetricsAPIServer{
		AdapterBase: cmd.AdapterBase{},
	}

	server.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(generatedopenapi.GetOpenAPIDefinitions, openapinamer.NewDefinitionNamer(apiserver.Scheme))
	server.OpenAPIConfig.Info.Title = "prometheus-metrics-adapter"
	server.OpenAPIConfig.Info.Version = "1.0.0"
	server.InstallFlags()

	return server
}
