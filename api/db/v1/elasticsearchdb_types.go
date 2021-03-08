package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

type ElasticsearchDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ElasticsearchDBSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

type ElasticsearchDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ElasticsearchDB `json:"items"`
}

type ElasticsearchDBSpec struct {
	Domain  string `json:"domain,omitempty"`
	Version string `json:"version,omitempty"`
	// +kubebuilder:validation:Optional
	Replicas int `json:"replicas,omitempty"`
	// +kubebuilder:validation:Optional
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// +kubebuilder:validation:Optional
	Storage Storage `json:"storage,omitempty"`
	// +kubebuilder:validation:Optional
	Heap string `json:"heap,omitempty"`
	// +kubebuilder:validation:Optional
	Ingress Ingress `json:"ingress,omitempty"`
}

type Ingress struct {
	Annotations map[string]string `json:"annotations,omitempty""`
}

func init() {
	SchemeBuilder.Register(&ElasticsearchDB{}, &ElasticsearchDBList{})
}
