package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GlobalRegistrySpec defines the desired state of GlobalRegistry
type GlobalRegistrySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Command string `json:"command"` //delete - tagName null 일 경우 전체 삭체, list, tagList

	//Command 가 delete 인 경우에만 포함.
	ImageName string `json:"imageName,omitempty"`
	TagName   string `json:"tagName,omitempty"`
}

// GlobalRegistryStatus defines the observed state of GlobalRegistry
type GlobalRegistryStatus struct {
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
// +kubebuilder:resource:path=noderegistries,scope=Namespaced
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// NodeRegistry is the Schema for the NodeRegistry API
// +kubebuilder:subresource:status
// +k8s:openapi-gen=true
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Created time stamp"
// +kubebuilder:printcolumn:name="IsSuccess",type="boolean",JSONPath=".status.GlobalRegistryStatus",description="-"
// +kubebuilder:printcolumn:name="Command",type="string",JSONPath=".spec.GlobalRegistrySpec[*].Command",description="-"
// +kubebuilder:printcolumn:name="ImageName",type="string",JSONPath=".spec.GlobalRegistrySpec[*].ImageName",description="-"
// +kubebuilder:printcolumn:name="TagName",type="string",JSONPath=".spec.GlobalRegistrySpec[*].TagName",description="-"
// +kubebuilder:printcolumn:name="REASON",type="string",JSONPath=".status.Reason",description="-"
// +kubebuilder:printcolumn:name="ElapsedTime",type="string",JSONPath=".status.ElapsedTime",description="ElapsedTime"
type GlobalRegistry struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GlobalRegistrySpec   `json:"spec,omitempty"`
	Status GlobalRegistryStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GlobalRegistryList contains a list of GlobalRegistry
type GlobalRegistryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GlobalRegistry `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GlobalRegistry{}, &GlobalRegistryList{})
}
