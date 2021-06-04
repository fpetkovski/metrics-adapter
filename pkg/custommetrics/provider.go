package custommetrics

import (
	"context"
	"fmt"
	"fpetkovski/prometheus-adapter/pkg/apis/v1alpha1"
	"fpetkovski/prometheus-adapter/pkg/query"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"

	"github.com/go-logr/logr"
	"github.com/prometheus/common/model"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

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
	logger        logr.Logger
	k8sClient     client.Client
	promClient    PromClient
	restMapper    meta.RESTMapper
	metricsLister *MetricsLister
	queryBuilder  *query.Builder
}

type PromClient interface {
	Query(ctx context.Context, query string, ts time.Time) (model.Value, prom.Warnings, error)
}

type QueryData struct {
	Name      string
	Namespace string
}

func NewMetricsProvider(
	restConfig *rest.Config,
	k8sClient client.Client,
	promClient PromClient,
	metricsLister *MetricsLister,
	logger logr.Logger,
) (provider.CustomMetricsProvider, error) {
	restMapper, err := apiutil.NewDynamicRESTMapper(restConfig)
	if err != nil {
		return nil, err
	}

	return &metricsProvider{
		logger:        logger,
		k8sClient:     k8sClient,
		promClient:    promClient,
		metricsLister: metricsLister,
		restMapper:    restMapper,
		queryBuilder:  query.NewQueryBuilder(),
	}, nil
}

func (m metricsProvider) GetMetricByName(
	name types.NamespacedName,
	info provider.CustomMetricInfo,
	metricSelector labels.Selector,
) (*custom_metrics.MetricValue, error) {
	m.logger.Info("Fetching metric",
		"ResourceName", name,
		"MetricInfo", info,
		"MetricSelector", metricSelector.String())

	metricResource, err := m.getMetric(info.Metric, name.Namespace)
	if err != nil {
		m.logger.Error(err, "could not get metric")
		return nil, err
	}

	k8sResource, err := m.getKubernetesResource(name, info.GroupResource)
	if err != nil {
		m.logger.Error(err, "could not get kubernetes resource")
		return nil, err
	}

	queryTemplate := metricResource.Spec.PrometheusQuery
	queryData := QueryData{
		Name:      k8sResource.GetName(),
		Namespace: k8sResource.GetNamespace(),
	}
	promQuery, err := m.queryBuilder.BuildQuery(queryTemplate, queryData)
	if err != nil {
		return nil, err
	}
	m.logger.Info("executing prometheus query", "query", promQuery)
	queryResult, _, err := m.promClient.Query(context.Background(), promQuery, time.Now())
	if err != nil {
		m.logger.Error(err, "could not execute prometheus query", "Query", promQuery)
		return nil, err
	}
	m.logger.Info("executed prometheus query", "result", queryResult)

	value, err := getQuantity(queryResult)
	if err != nil {
		m.logger.Error(err, "could not get quantity")
		return nil, err
	}

	m.logger.Info("query result", "Value", value)
	if value == nil {
		return nil, fmt.Errorf("not found")
	}

	return &custom_metrics.MetricValue{
		DescribedObject: custom_metrics.ObjectReference{
			Kind:            k8sResource.GetKind(),
			Namespace:       name.Namespace,
			Name:            name.Name,
			UID:             k8sResource.GetUID(),
			APIVersion:      k8sResource.GetAPIVersion(),
			ResourceVersion: k8sResource.GetResourceVersion(),
		},
		Metric: custom_metrics.MetricIdentifier{
			Name: info.Metric,
		},
		Timestamp: metav1.NewTime(time.Now()),
		Value:     *value,
	}, nil
}

func (m metricsProvider) GetMetricBySelector(
	namespace string,
	selector labels.Selector,
	info provider.CustomMetricInfo,
	metricSelector labels.Selector,
) (*custom_metrics.MetricValueList, error) {
	m.logger.Info("Fetching metric for resources",
		"Namespace", namespace,
		"ResourceSelector", selector.String(),
		"MetricInfo", info,
		"MetricSelector", metricSelector.String())

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
	return m.metricsLister.GetMetricInfos()
}

func (m metricsProvider) getMetric(name string, namespace string) (*v1alpha1.CustomMetric, error) {
	metricNamespacedName := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
	var metric v1alpha1.CustomMetric
	if err := m.k8sClient.Get(context.Background(), metricNamespacedName, &metric); err != nil {
		return nil, err
	}

	return &metric, nil
}

func (m metricsProvider) getKubernetesResource(name types.NamespacedName, groupResource schema.GroupResource) (*unstructured.Unstructured, error) {
	kind, err := m.restMapper.KindFor(groupResource.WithVersion(""))
	if err != nil {
		return nil, err
	}

	var k8sResource unstructured.Unstructured
	k8sResource.SetGroupVersionKind(kind)
	if err := m.k8sClient.Get(context.Background(), name, &k8sResource); err != nil {
		return nil, err
	}

	return &k8sResource, nil
}

func getQuantity(samples model.Value) (*resource.Quantity, error) {
	vector, ok := samples.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("query result must be of type ValVector")
	}

	if len(vector) == 0 {
		return nil, nil
	}

	if len(vector) > 1 {
		return nil, fmt.Errorf("query must return one or no samples, got %d", len(vector))
	}

	value, err := resource.ParseQuantity(fmt.Sprintf("%f", vector[0].Value))
	if err != nil {
		return nil, err
	}

	return &value, nil
}
