package external_metrics

import (
	"context"
	"fpetkovski/prometheus-adapter/pkg/apis/v1alpha1"

	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type reconciler struct {
}

func registerController(mgr manager.Manager) error {
	ctrl, err := controller.New("custom-metrics", mgr, controller.Options{
		Reconciler: &reconciler{},
	})
	if err != nil {
		return err
	}

	if err := ctrl.Watch(&source.Kind{Type: &v1alpha1.PrometheusMetric{}}, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	return nil
}

func (r *reconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}
