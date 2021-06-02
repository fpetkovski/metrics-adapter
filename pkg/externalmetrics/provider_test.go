package externalmetrics_test

import (
	"context"
	"fmt"
	"fpetkovski/prometheus-adapter/pkg/apis/v1alpha1"
	"fpetkovski/prometheus-adapter/pkg/custommetrics"
	"fpetkovski/prometheus-adapter/pkg/externalmetrics"
	"strings"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/selection"

	"github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/provider"
	prom "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2/klogr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type prometheusStub struct {
	query string
	value string
	err   error
}

func (p prometheusStub) Query(ctx context.Context, query string, ts time.Time) (model.Value, prom.Warnings, error) {
	if p.err != nil {
		return nil, nil, p.err
	}

	if p.query == query {
		timestamp := time.Date(2012, 10, 31, 0, 0, 0, 0, time.UTC).Unix()
		decoder := expfmt.SampleDecoder{
			Dec: expfmt.NewDecoder(strings.NewReader(p.value), expfmt.FmtText),
			Opts: &expfmt.DecodeOptions{
				Timestamp: model.Time(timestamp),
			},
		}
		result := model.Vector{}
		err := decoder.Decode(&result)
		if err != nil {
			return nil, nil, err
		}

		return result, nil, nil
	}

	return nil, nil, fmt.Errorf("query not found")
}

func TestGetExternalMetric(t *testing.T) {
	promStub := prometheusStub{
		query: "container_cpu_usage_seconds_total",
		value: `
# TYPE container_cpu_usage_seconds_total counter
container_cpu_usage_seconds_total{app="nginx"} 10.5
`,
		err: nil,
	}
	k8sClient, metricProvider := makeProvider(t, promStub)
	metric := makeMetric("test-case-metric", "container_cpu_usage_seconds_total{ {{.Labels }} }")
	defer func() {
		_ = k8sClient.Delete(context.Background(), &metric)
	}()
	if err := k8sClient.Create(context.Background(), &metric); err != nil {
		t.Fatal(err)
	}

	metricInfo := provider.ExternalMetricInfo{
		Metric: metric.Name,
	}
	selector := makeMetricSelector(t, map[string]string{"app": "nginx"})
	_, err := metricProvider.GetExternalMetric("default", selector, metricInfo)
	if err != nil {
		t.Fatal(err)
	}

	//expectedValue := resource.MustParse("10.5")
	//if !expectedValue.Equal(result.Value) {
	//	t.Fatalf("Invalid metric result, got %s, want %s", result.Value.String(), expectedValue.String())
	//}
}

func makeMetricSelector(t *testing.T, selectorLabels map[string]string) labels.Selector {
	selector := labels.NewSelector()
	for k, v := range selectorLabels {
		requirement, err := labels.NewRequirement(k, selection.Equals, []string{v})
		if err != nil {
			t.Fatal(err)
		}
		selector = selector.Add(*requirement)
	}

	return selector
}

func makeProvider(t *testing.T, promClient custommetrics.PromClient) (client.Client, provider.ExternalMetricsProvider) {
	cfg := config.GetConfigOrDie()

	k8sClient, err := client.New(cfg, client.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if err := v1alpha1.AddToScheme(k8sClient.Scheme()); err != nil {
		t.Fatal(err)
	}

	metricsProvider := externalmetrics.NewMetricsProvider(k8sClient, promClient, klogr.New())

	return k8sClient, metricsProvider
}

func makeMetric(metricName string, promQuery string) v1alpha1.ExternalMetric {
	metric := v1alpha1.ExternalMetric{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CustomMetric",
			APIVersion: "fpetkovski.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      metricName,
			Namespace: "default",
		},
		Spec: v1alpha1.ExternalMetricSpec{
			PrometheusQuery: promQuery,
		},
	}
	return metric
}
