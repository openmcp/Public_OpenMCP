package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NodeRegistrySpec defines the desired state of NodeRegistry
type NodeRegistrySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Command string `json:"command"` //push, pull - nodeName 이 없을 경우 Cluster 단위 명령

	//Command 가 delete 인 경우에만 포함.
	ImageName string `json:"imageName,omitempty"`
	TagName   string `json:"tagName,omitempty"`

	ClusterName string `json:"clusterName,omitempty"`
	NodeName    string `json:"nodeName,omitempty"`
}

// NodeRegistryStatus defines the observed state of NodeRegistry
type NodeRegistryStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Succeeded indicates if the backup has Succeeded.
	Succeeded bool `json:"succeeded"`
	// Reason indicates the reason for any backup related failures.
	Reason string `json:"Reason,omitempty"`
	// LastSuccessDate indicate the time to get snapshot last time
	LastSuccessDate metav1.Time `json:"lastSuccessDate,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeRegistry is the Schema for the noderegistries API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=noderegistries,scope=Namespaced
type NodeRegistry struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeRegistrySpec   `json:"spec,omitempty"`
	Status NodeRegistryStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeRegistryList contains a list of NodeRegistry
type NodeRegistryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeRegistry `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeRegistry{}, &NodeRegistryList{})
}
