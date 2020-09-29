package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SyncSpec defines the desired state of Sync
// +k8s:openapi-gen=true
type SyncSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	//Template appsv1.Deployment `json:"template" protobuf:"bytes,3,opt,name=template"`
	//Kind string `json:"kind" protobuf:"varint,1,opt,name=kind"`
	ClusterName string      `json:"clustername" protobuf:"varint,1,opt,name=clustername"`
	Command     string      `json:"command" protobuf:"varint,1,opt,name=command"`
	Template    interface{} `json:"template" protobuf:"bytes,3,opt,name=template"`


}

// SyncStatus defines the observed state of Sync
// +k8s:openapi-gen=true
type SyncStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Sync is the Schema for the syncs API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Sync struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SyncSpec   `json:"spec,omitempty"`
	Status SyncStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SyncList contains a list of Sync
type SyncList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Sync `json:"items"`
}

//func init() {
//	SchemeBuilder.Register(&Sync{}, &SyncList{})
//}
