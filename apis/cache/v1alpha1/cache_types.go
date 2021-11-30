package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CacheSpec defines the desired state of Cache
type CacheSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Timer      string `json:"timer"` //node name
	CacheCount int    `json:"cachecount"`
}

// CacheStatus defines the observed state of Cache
type CacheStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Succeeded indicates if the backup has Succeeded.
	Succeeded bool   `json:"succeeded"`
	History   []Data `json:"History,omitempty"`
	// Reason indicates the reason for any backup related failures.
	Reason string `json:"Reason,omitempty"`
}

type Data struct {
	ImageList []ImageInfo `json:"ImageList,omitempty"`
	Timestamp string      `json:"Timestamp,omitempty"`
}

type ImageInfo struct {
	ImageName  string `json:"ImageNmae,omitempty"`
	ImageCount int64  `json:"ImageCount,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cache is the Schema for the Caches API
// +kubebuilder:subresource:status
// +k8s:openapi-gen=true
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Created time stamp"
// +kubebuilder:printcolumn:name="IsSuccess",type="boolean",JSONPath=".status.succeeded",description="-"
// +kubebuilder:printcolumn:name="REASON",type="string",JSONPath=".status.Reason",description="-"
type Cache struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CacheSpec   `json:"spec,omitempty"`
	Status CacheStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CacheList contains a list of Cache
type CacheList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cache `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cache{}, &CacheList{})
}
