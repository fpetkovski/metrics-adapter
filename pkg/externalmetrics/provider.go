package externalmetrics

import (
	"github.com/go-logr/logr"
	"github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/provider"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/metrics/pkg/apis/external_metrics"
	"time"
)

type metricsProvider struct {
	logger logr.Logger
}

func NewMetricsProvider(logger logr.Logger) provider.ExternalMetricsProvider {
	return &metricsProvider{
		logger: logger,
	}
}

func (m metricsProvider) GetExternalMetric(namespace string, metricSelector labels.Selector, info provider.ExternalMetricInfo) (*external_metrics.ExternalMetricValueList, error) {
	return &external_metrics.ExternalMetricValueList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "",
			APIVersion: "",
		},
		Items: []external_metrics.ExternalMetricValue{
			{
				TypeMeta: metav1.TypeMeta{
					Kind:       "",
					APIVersion: "",
				},
				MetricName:   "",
				MetricLabels: nil,
				Timestamp: metav1.Time{
					Time: time.Time{},
				},
				WindowSeconds: nil,
				Value: resource.Quantity{
					Format: "",
				},
			},
		},
	}, nil
}

func (m metricsProvider) ListAllExternalMetrics() []provider.ExternalMetricInfo {
	return []provider.ExternalMetricInfo{
		{
			Metric: "External requests",
		},
	}
}
