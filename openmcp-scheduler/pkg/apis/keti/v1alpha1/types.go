package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OpenMCPschedulerSpec defines the desired state of OpenMCPscheduler
// +k8s:openapi-gen=true
type OpenMCPschedulerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Template appsv1.Deployment `json:"template" protobuf:"bytes,3,opt,name=template"`
	Replicas int32             `json:"replicas" protobuf:"varint,1,opt,name=replicas"`
	Policy   []PolicyData      `json:"policy" protobuf:"bytes,1,opt,name=policy"`
	//Placement
}

// OpenMCPschedulerStatus defines the observed state of OpenMCPscheduler
// +k8s:openapi-gen=true
type OpenMCPschedulerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Replicas    int32            `json:"replicas"`
	Policy      []PolicyData     `json:"policy"`
	ClusterMaps map[string]int32 `json:"clusters"`
}

type PolicyData struct {
	Rule string `json:"rule"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPscheduler is the Schema for the openmcpschedulers API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type OpenMCPscheduler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenMCPschedulerSpec   `json:"spec,omitempty"`
	Status OpenMCPschedulerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPschedulerList contains a list of OpenMCPscheduler
type OpenMCPschedulerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenMCPscheduler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenMCPscheduler{}, &OpenMCPschedulerList{})
}
