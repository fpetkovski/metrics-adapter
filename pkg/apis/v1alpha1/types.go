package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PrometheusMetric is the Schema for the PrometheusMetric API
// +k8s:openapi-gen=true
type PrometheusMetric struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PrometheusMetricSpec   `json:"spec,omitempty"`
	Status PrometheusMetricStatus `json:"status,omitempty"`
}

type PrometheusMetricSpec struct {
	PrometheusQuery string `json:"prometheusQuery"`
}

type PrometheusMetricStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PrometheusMetricList contains a list of PrometheusMetric
type PrometheusMetricList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PrometheusMetric `json:"items"`
}
