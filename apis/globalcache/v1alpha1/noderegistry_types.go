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
	ClusterName string      `json:"clusterName,omitempty"`
	NodeName    string      `json:"nodeName,omitempty"`
	ImageLists  []ImageList `json:"ImageList"`
}
type ImageList struct {
	ImageName string `json:"imageName,omitempty"`
	TagName   string `json:"tagName,omitempty"`
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

	ElapsedTime string `json:"ElapsedTime,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeRegistry is the Schema for the noderegistries API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=noderegistries,scope=Namespaced

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeRegistry is the Schema for the NodeRegistry API
// +kubebuilder:subresource:status
// +k8s:openapi-gen=true
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Created time stamp"
// +kubebuilder:printcolumn:name="IsSuccess",type="boolean",JSONPath=".status.NodeRegistryStatus",description="-"
// +kubebuilder:printcolumn:name="Command",type="string",JSONPath=".spec.NodeRegistrySpec[*].Command",description="-"
// +kubebuilder:printcolumn:name="ClusterName",type="string",JSONPath=".spec.NodeRegistrySpec[*].ClusterName",description="-"
// +kubebuilder:printcolumn:name="NodeName",type="string",JSONPath=".spec.NodeRegistrySpec[*].NodeName",description="-"
// +kubebuilder:printcolumn:name="REASON",type="string",JSONPath=".status.Reason",description="-"
// +kubebuilder:printcolumn:name="ElapsedTime",type="string",JSONPath=".status.ElapsedTime",description="ElapsedTime"
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
