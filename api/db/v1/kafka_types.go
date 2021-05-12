package v1

import (
	kommonsv1 "github.com/flanksource/kommons/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Kafka Replicas",type="integer",JSONPath=".spec.kafka.replicas"
// +kubebuilder:printcolumn:name="Mem",type="string",JSONPath=".spec.kafka.resources.requests.memory"
// +kubebuilder:printcolumn:name="CPU",type="string",JSONPath=".spec.kafka.resources.requests.cpu"
// +kubebuilder:printcolumn:name="Storage",type="string",JSONPath=".spec.kafka.storage.size"
// +kubebuilder:printcolumn:name="ZK Replicas",type="integer",JSONPath=".spec.zookeeper.replicas"
// +kubebuilder:printcolumn:name="Mem",type="string",JSONPath=".spec.zookeeper.resources.requests.memory"
// +kubebuilder:printcolumn:name="CPU",type="string",JSONPath=".spec.zookeeper.resources.requests.cpu"
// +kubebuilder:printcolumn:name="Storage",type="string",JSONPath=".spec.zookeeper.storage.size"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[].status"
type KafkaCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              KafkaClusterSpec   `json:"spec,omitempty"`
	Status            KafkaClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type KafkaClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KafkaCluster `json:"items"`
}

type KafkaClusterSpec struct {
	Kafka      KafkaSpec     `json:"kafka"`
	Zoopkeeper ZookeeperSpec `json:"zookeeper"`
}

type KafkaSpec struct {
	PodResources `json:",inline"`
	Storage      Storage `json:"storage,omitempty"`
	// +kubebuilder:validation:Optional
	Version string `json:"version"`
}

type ZookeeperSpec struct {
	PodResources `json:",inline"`
	Storage      Storage `json:"storage,omitempty"`
	// +kubebuilder:validation:Optional
	Version string `json:"version"`
}

type KafkaClusterStatus struct {
	Conditions kommonsv1.ConditionList `json:"conditions"`
}

func init() {
	SchemeBuilder.Register(&KafkaCluster{}, &KafkaClusterList{})
}
