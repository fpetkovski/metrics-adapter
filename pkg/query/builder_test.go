package query_test

import (
	"fpetkovski/prometheus-adapter/pkg/externalmetrics"
	"fpetkovski/prometheus-adapter/pkg/query"
	"testing"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func TestBuildQuery(t *testing.T) {
	builder := query.Builder{}

	template := `cpu_usage{pod="{{ .Name }}", namespace="{{.Namespace}}", label="{{.Label}}"}`
	queryData := struct {
		Name      string
		Namespace string
		Label     string
	}{
		Name:      "test-pod",
		Namespace: "test-namespace",
		Label:     "test-label",
	}
	result := builder.BuildQuery(template, queryData)
	expected := `cpu_usage{pod="test-pod", namespace="test-namespace", label="test-label"}`
	if result != expected {
		t.Fatalf("invalid result, got %s, want %s", result, expected)
	}
}

func TestBuildQueryWithRequirement(t *testing.T) {
	builder := query.Builder{}

	template := `cpu_usage{ {{ .Labels }} }`

	r1, _ := labels.NewRequirement("app", selection.Equals, []string{"nginx"})
	r2, _ := labels.NewRequirement("system", selection.Equals, []string{"monitoring"})
	queryData := struct {
		Labels externalmetrics.LabelSelectorList
	}{
		Labels: []externalmetrics.LabelSelector{
			externalmetrics.LabelSelector(*r1),
			externalmetrics.LabelSelector(*r2),
		},
	}
	result := builder.BuildQuery(template, queryData)
	expected := `cpu_usage{ app="nginx", system="monitoring" }`
	if result != expected {
		t.Fatalf("invalid result, got %s, want %s", result, expected)
	}
}
