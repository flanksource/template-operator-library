package v1

import (
	kommonsv1 "github.com/flanksource/kommons/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

type ElasticsearchDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ElasticsearchDBSpec   `json:"spec,omitempty"`
	Status            ElasticsearchDBStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type ElasticsearchDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ElasticsearchDB `json:"items"`
}

type ElasticsearchDBSpec struct {
	PodResources `json:",inline"`
	Domain       string `json:"domain,omitempty"`
	Version      string `json:"version,omitempty"`
	// +kubebuilder:validation:Optional
	Storage Storage `json:"storage,omitempty"`
	// +kubebuilder:validation:Optional
	Heap string `json:"heap,omitempty"`
	// +kubebuilder:validation:Optional
	Ingress Ingress `json:"ingress,omitempty"`
}

type ElasticsearchDBStatus struct {
	Conditions kommonsv1.ConditionList `json:"conditions"`
}

type Ingress struct {
	Annotations map[string]string `json:"annotations,omitempty""`
}

func init() {
	SchemeBuilder.Register(&ElasticsearchDB{}, &ElasticsearchDBList{})
}
