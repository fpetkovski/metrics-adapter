package controller_manager

import (
	"fpetkovski/prometheus-adapter/pkg/apis/v1alpha1"

	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func New(cfg *rest.Config) (manager.Manager, error) {
	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		return nil, err
	}

	if err := v1alpha1.AddToScheme(mgr.GetScheme()); err != nil {
		return nil, err
	}

	return mgr, nil
}
