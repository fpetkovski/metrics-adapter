package custommetrics

import (
	"context"
	"fpetkovski/prometheus-adapter/pkg/apis/v1alpha1"
	"strings"

	"github.com/go-logr/logr"
	"github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/provider"
	v1 "k8s.io/api/autoscaling/v1"
	"k8s.io/api/autoscaling/v2beta2"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type reconciler struct {
	k8sClient     client.Client
	logger        logr.Logger
	restMapper    meta.RESTMapper
	metricsLister *MetricsLister
}

func RegisterController(config *rest.Config, mgr manager.Manager, logger logr.Logger) (*MetricsLister, error) {
	restMapper, err := apiutil.NewDynamicRESTMapper(config)
	if err != nil {
		return nil, err
	}

	metricsLister := NewMetricsLister()

	r := reconciler{
		k8sClient:     mgr.GetClient(),
		logger:        logger,
		restMapper:    restMapper,
		metricsLister: metricsLister,
	}

	ctrl, err := controller.New("custom-metrics", mgr, controller.Options{
		Reconciler: &r,
	})
	if err != nil {
		return nil, err
	}

	if err := ctrl.Watch(&source.Kind{Type: &v1alpha1.CustomMetric{}}, &handler.EnqueueRequestForObject{}); err != nil {
		return nil, err
	}

	if err := ctrl.Watch(&source.Kind{Type: &v1.HorizontalPodAutoscaler{}}, &handler.EnqueueRequestForObject{}); err != nil {
		return nil, err
	}

	return metricsLister, nil
}

func (r *reconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	var metrics v1alpha1.CustomMetricList
	if err := r.k8sClient.List(ctx, &metrics); err != nil {
		return reconcile.Result{}, err
	}

	var hpas v2beta2.HorizontalPodAutoscalerList
	if err := r.k8sClient.List(ctx, &hpas); err != nil {
		return reconcile.Result{}, nil
	}

	metricInfos := make([]provider.CustomMetricInfo, 0)
	for _, hpa := range hpas.Items {
		infos, err := r.getMetricInfos(hpa.Spec.Metrics)
		if err != nil {
			return reconcile.Result{}, err
		}

		metricInfos = append(metricInfos, infos...)
	}

	r.metricsLister.setMetricInfos(metricInfos)

	return reconcile.Result{}, nil
}

func (r *reconciler) getMetricInfos(metrics []v2beta2.MetricSpec) ([]provider.CustomMetricInfo, error) {
	infos := make([]provider.CustomMetricInfo, 0)
	for _, metric := range metrics {
		// Skip external metrics
		if metric.External != nil {
			continue
		}
		info, err := r.getMetricInfo(metric)
		if err != nil {
			return nil, err
		}

		infos = append(infos, info)
	}

	return infos, nil
}

func (r *reconciler) getMetricInfo(metric v2beta2.MetricSpec) (provider.CustomMetricInfo, error) {
	if metric.Pods != nil {
		return makeObjectMetricInfo(metric.Pods.Metric.Name, "", "pods"), nil
	}

	group := strings.Split(metric.Object.DescribedObject.APIVersion, "/")[0]
	groupKind := schema.GroupKind{
		Group: group,
		Kind:  metric.Object.DescribedObject.Kind,
	}
	mapping, err := r.restMapper.RESTMapping(groupKind)
	if err != nil {
		return provider.CustomMetricInfo{}, err
	}

	return makeObjectMetricInfo(metric.Object.Metric.Name, group, mapping.Resource.Resource), nil
}

func makeObjectMetricInfo(metricName string, group string, resource string) provider.CustomMetricInfo {
	return provider.CustomMetricInfo{
		GroupResource: schema.GroupResource{
			Group:    group,
			Resource: resource,
		},
		Namespaced: true,
		Metric:     metricName,
	}
}
