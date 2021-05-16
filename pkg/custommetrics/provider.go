package custommetrics

import (
	"context"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"time"

	"github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/provider"
	prom "github.com/prometheus/client_golang/api/prometheus/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/metrics/pkg/apis/custom_metrics"
)

type metricsProvider struct {
	logger     logr.Logger
	k8sClient  client.Client
	prometheus prom.API
}

func NewMetricsProvider(logger logr.Logger, k8sClient client.Client, prometheus prom.API) provider.CustomMetricsProvider {
	return &metricsProvider{
		logger:     logger,
		k8sClient:  k8sClient,
		prometheus: prometheus,
	}
}

func (m metricsProvider) GetMetricByName(name types.NamespacedName, info provider.CustomMetricInfo, metricSelector labels.Selector) (*custom_metrics.MetricValue, error) {
	m.logger.Info("Fetching metric by name", "Name", name, "Info", info, "Selector", metricSelector)

	return &custom_metrics.MetricValue{
		TypeMeta: metav1.TypeMeta{},
		DescribedObject: custom_metrics.ObjectReference{
			Kind:            "Pod",
			Namespace:       name.Namespace,
			Name:            name.Name,
			UID:             "10",
			APIVersion:      "",
			ResourceVersion: "2",
		},
		Metric: custom_metrics.MetricIdentifier{
			Name: "rps",
		},
		Timestamp: metav1.NewTime(time.Now()),
		Value:     resource.MustParse("4"),
	}, nil
}

func (m metricsProvider) GetMetricBySelector(namespace string, selector labels.Selector, info provider.CustomMetricInfo, metricSelector labels.Selector) (*custom_metrics.MetricValueList, error) {
	m.logger.Info("Fetching metric by selector",
		"Namespace", namespace,
		"Selector", selector,
		"Info", info,
		"MetricSelector", metricSelector)

	var pods v1.PodList
	err := m.k8sClient.List(context.TODO(), &pods, &client.ListOptions{
		LabelSelector: metricSelector,
		Namespace:     namespace,
	})
	if err != nil {
		return nil, err
	}

	metrics := custom_metrics.MetricValueList{
		Items: make([]custom_metrics.MetricValue, len(pods.Items)),
	}
	for i, pod := range pods.Items {
		value := 100 - (i+1)*20
		metrics.Items[i] = custom_metrics.MetricValue{
			TypeMeta: metav1.TypeMeta{},
			DescribedObject: custom_metrics.ObjectReference{
				Kind:            "Pod",
				Namespace:       namespace,
				Name:            pod.Name,
				UID:             pod.UID,
				APIVersion:      pod.APIVersion,
				ResourceVersion: pod.ResourceVersion,
			},
			Metric: custom_metrics.MetricIdentifier{
				Name: info.Metric,
			},
			Timestamp: metav1.NewTime(time.Now()),
			Value:     resource.MustParse(strconv.Itoa(value)),
		}
	}

	return &metrics, nil
}

func (m metricsProvider) ListAllMetrics() []provider.CustomMetricInfo {
	m.logger.Info("Listing all metrics")

	return []provider.CustomMetricInfo{
		{
			GroupResource: schema.GroupResource{
				Group:    "v1",
				Resource: "pods",
			},
			Namespaced: true,
			Metric:     "rps",
		},
		{
			GroupResource: schema.GroupResource{
				Group:    "networking.k8s.io",
				Resource: "ingresses",
			},
			Namespaced: true,
			Metric:     "rps",
		},
	}
}
