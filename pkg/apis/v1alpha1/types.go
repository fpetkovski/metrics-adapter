package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// CustomMetricSpec defines the desired state of CustomMetric
type CustomMetricSpec struct {
	PrometheusQuery string `json:"prometheusQuery"`
}

// CustomMetricStatus defines the observed state of CustomMetric.
// It should always be reconstructable from the state of the cluster and/or outside world.
type CustomMetricStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomMetric is the Schema for the custommetrics API
// +k8s:openapi-gen=true
type CustomMetric struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomMetricSpec   `json:"spec"`
	Status CustomMetricStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomMetricList contains a list of CustomMetric
type CustomMetricList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CustomMetric `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExternalMetric is the Schema for the externalmetrics API
// +k8s:openapi-gen=true
type ExternalMetric struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExternalMetricSpec   `json:"spec,omitempty"`
	Status ExternalMetricStatus `json:"status,omitempty"`
}

type ExternalMetricSpec struct {
	PrometheusQuery string `json:"prometheusQuery"`
}

type ExternalMetricStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExternalMetricList contains a list of ExternalMetric
type ExternalMetricList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExternalMetric `json:"items"`
}
