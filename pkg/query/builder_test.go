package query_test

import (
	"fpetkovski/prometheus-adapter/pkg/query"
	"testing"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func TestBuildQueryS(t *testing.T) {
	ts := []struct {
		template      string
		data          interface{}
		expectedQuery string
	}{
		{
			template: `cpu_usage{pod="{{ .Name }}", namespace="{{.Namespace}}"}`,
			data: struct {
				Name      string
				Namespace string
			}{
				Name:      "test-pod",
				Namespace: "test-namespace",
			},
			expectedQuery: `cpu_usage{pod="test-pod", namespace="test-namespace"}`,
		},
		{
			template: `cpu_usage{ {{ .Labels }} }`,
			data: struct {
				Labels query.LabelSelectors
			}{
				Labels: makeLabelSelectors(
					makeRequirement("app", selection.In, []string{"haproxy", "nginx"}),
					makeRequirement("system", selection.Equals, []string{"monitoring"}),
				),
			},
			expectedQuery: `cpu_usage{ app=~"^(haproxy|nginx)$", system="monitoring" }`,
		},
		{
			template: `cpu_usage{ {{ index .Labels "app" }} }`,
			data: struct {
				Labels query.LabelSelectors
			}{
				Labels: makeLabelSelectors(
					makeRequirement("app", selection.Equals, []string{"nginx"}),
					makeRequirement("system", selection.Equals, []string{"monitoring"}),
				),
			},
			expectedQuery: `cpu_usage{ app="nginx" }`,
		},
	}

	for _, tt := range ts {
		builder := query.Builder{}
		result, err := builder.BuildQuery(tt.template, tt.data)
		if err != nil {
			t.Fatal(err)
		}
		if result != tt.expectedQuery {
			t.Fatalf("invalid result, got %s, want %s", result, tt.expectedQuery)
		}
	}
}

func makeLabelSelectors(requirements ...labels.Requirement) query.LabelSelectors {
	return query.NewLabelSelectors(requirements)
}

func makeRequirement(key string, op selection.Operator, vals []string) labels.Requirement {
	r, _ := labels.NewRequirement(key, op, vals)
	return *r
}
