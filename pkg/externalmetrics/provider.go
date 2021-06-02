package externalmetrics

import (
	"context"
	"fpetkovski/prometheus-adapter/pkg/apis/v1alpha1"
	"fpetkovski/prometheus-adapter/pkg/query"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/types"

	prom "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	"github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/provider"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/metrics/pkg/apis/external_metrics"
)

type metricsProvider struct {
	logger       logr.Logger
	k8sClient    client.Client
	promClient   PromClient
	queryBuilder *query.Builder
}

type PromClient interface {
	Query(ctx context.Context, query string, ts time.Time) (model.Value, prom.Warnings, error)
}

type QueryData struct {
	Labels string
}

func NewMetricsProvider(k8sClient client.Client, promClient PromClient, logger logr.Logger) provider.ExternalMetricsProvider {
	return &metricsProvider{
		logger:       logger,
		k8sClient:    k8sClient,
		promClient:   promClient,
		queryBuilder: query.NewQueryBuilder(),
	}
}

func (m metricsProvider) GetExternalMetric(namespace string, metricSelector labels.Selector, info provider.ExternalMetricInfo) (*external_metrics.ExternalMetricValueList, error) {
	m.logger.Info("Fetching external metrics",
		"Namespace", namespace,
		"MetricSelector", metricSelector.String(),
		"MetricInfo", info)

	name := types.NamespacedName{
		Namespace: namespace,
		Name:      info.Metric,
	}
	var metric v1alpha1.ExternalMetric
	if err := m.k8sClient.Get(context.Background(), name, &metric); err != nil {
		return nil, err
	}

	selectorRequirements, _ := metricSelector.Requirements()
	selectors := make([]string, len(selectorRequirements))
	for i, r := range selectorRequirements {
		selectors[i] = LabelSelector(r).String()
	}

	queryTpl := metric.Spec.PrometheusQuery
	queryData := QueryData{
		Labels: strings.Join(selectors, ", "),
	}
	promQuery := m.queryBuilder.BuildQuery(queryTpl, queryData)
	m.logger.Info("Executing prometheus query", "query", promQuery)

	return &external_metrics.ExternalMetricValueList{
		Items: []external_metrics.ExternalMetricValue{
			{
				MetricName:   info.Metric,
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
	var metricsList v1alpha1.ExternalMetricList
	_ = m.k8sClient.List(context.Background(), &metricsList)

	metricInfos := make([]provider.ExternalMetricInfo, len(metricsList.Items))
	for i, metric := range metricsList.Items {
		metricInfos[i] = provider.ExternalMetricInfo{
			Metric: metric.Name,
		}
	}

	return metricInfos
}
