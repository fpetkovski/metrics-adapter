package external_metrics

import (
	"context"
	"fmt"
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
	m.logger.Info("Getting external metrics",
		"Namespace", namespace,
		"MetricSelector", metricSelector.String(),
		"MetricInfo", info)

	name := types.NamespacedName{
		Namespace: namespace,
		Name:      info.Metric,
	}
	var metric v1alpha1.PrometheusMetric
	if err := m.k8sClient.Get(context.Background(), name, &metric); err != nil {
		return nil, err
	}

	selectorRequirements, _ := metricSelector.Requirements()
	selectors := make([]string, len(selectorRequirements))
	for i, r := range selectorRequirements {
		selectors[i] = query.LabelSelector(r).String()
	}

	queryTpl := metric.Spec.PrometheusQuery
	queryData := QueryData{
		Labels: strings.Join(selectors, ", "),
	}
	promQuery, err := m.queryBuilder.BuildQuery(queryTpl, queryData)
	if err != nil {
		return nil, err
	}
	m.logger.Info("Executing prometheus query", "query", promQuery)
	queryResult, _, err := m.promClient.Query(context.TODO(), promQuery, time.Now())
	if err != nil {
		m.logger.Error(err, "could not execute prometheus query", "Query", promQuery)
		return nil, err
	}
	m.logger.Info("executed prometheus query", "result", queryResult)
	metrics, err := resultToMetrics(info.Metric, queryResult)
	if err != nil {
		return nil, err
	}

	return &external_metrics.ExternalMetricValueList{
		Items: metrics,
	}, nil
}

func resultToMetrics(metricName string, samples model.Value) ([]external_metrics.ExternalMetricValue, error) {
	vector, ok := samples.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("query result must be of type ValVector")
	}

	var metrics []external_metrics.ExternalMetricValue

	if len(vector) == 0 {
		return metrics, nil
	}

	for _, item := range vector {
		value, err := resource.ParseQuantity(fmt.Sprintf("%f", item.Value))
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, external_metrics.ExternalMetricValue{
			MetricName:   metricName,
			MetricLabels: extractLabels(item.Metric),
			Value:        value,
		})
	}

	return metrics, nil
}

func extractLabels(m model.Metric) map[string]string {
	labelSet := map[model.LabelName]model.LabelValue(model.LabelSet(m))

	result := make(map[string]string)
	for k, v := range labelSet {
		result[string(k)] = string(v)
	}

	return result
}

func (m metricsProvider) ListAllExternalMetrics() []provider.ExternalMetricInfo {
	var metricsList v1alpha1.PrometheusMetricList
	_ = m.k8sClient.List(context.Background(), &metricsList)

	metricInfos := make([]provider.ExternalMetricInfo, len(metricsList.Items))
	for i, metric := range metricsList.Items {
		metricInfos[i] = provider.ExternalMetricInfo{
			Metric: metric.Name,
		}
	}

	return metricInfos
}
