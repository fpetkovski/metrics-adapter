package query_test

import (
	"fpetkovski/prometheus-adapter/pkg/query"
	"testing"
)

func TestBuildQuery(t *testing.T) {
	builder := query.Builder{}

	template := `cpu_usage{pod="{{ .Name }}", namespace="{{.Namespace}}", label="{{.Label}}"}`
	queryData := struct {
		Name string
		Namespace string
		Label string
	}{
		Name: "test-pod",
		Namespace: "test-namespace",
		Label: "test-label",
	}
	result := builder.BuildQuery(template, queryData)
	expected := `cpu_usage{pod="test-pod", namespace="test-namespace", label="test-label"}`
	if result != expected {
		t.Fatalf("invalid result, got %s, want %s", result, expected)
	}
}
