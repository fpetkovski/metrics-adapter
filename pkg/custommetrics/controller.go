package custommetrics

import (
	"context"
	"fpetkovski/prometheus-adapter/pkg/apis/v1alpha1"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func RegisterController(logger logr.Logger, mgr manager.Manager) error {
	r := reconciler{
		logger: logger,
	}

	ctrl, err := controller.New("custom-metrics", mgr, controller.Options{
		Reconciler: &r,
	})
	if err != nil {
		return err
	}

	if err := ctrl.Watch(&source.Kind{Type: &v1alpha1.CustomMetric{}}, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	return nil
}

type reconciler struct {
	logger logr.Logger
}

func (r *reconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	r.logger.Info("Reconciling Custom Metric", "Name", request.Name)
	return reconcile.Result{}, nil
}
