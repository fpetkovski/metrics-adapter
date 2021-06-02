package custommetrics

import (
	"sync"

	"github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/provider"
)

type MetricsLister struct {
	m           sync.Mutex
	metricInfos []provider.CustomMetricInfo
}

func NewMetricsLister() *MetricsLister {
	return &MetricsLister{
		m:           sync.Mutex{},
		metricInfos: make([]provider.CustomMetricInfo, 0),
	}
}

func (l *MetricsLister) setMetricInfos(metricInfos []provider.CustomMetricInfo) {
	l.m.Lock()
	defer l.m.Unlock()

	l.metricInfos = metricInfos
}

func (l *MetricsLister) GetMetricInfos() []provider.CustomMetricInfo {
	l.m.Lock()
	defer l.m.Unlock()

	return l.metricInfos
}
