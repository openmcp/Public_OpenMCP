package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SnapshotRestoreSpec defines the desired state of SnapshotRestore
type SnapshotRestoreSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	SnapshotRestoreSource []SnapshotRestoreSource `json:"snapshotRestoreSource"`
}

// SnapshotRestoreSource contains the supported SnapshotRestore sources.
type SnapshotRestoreSource struct {
	ResourceCluster   string `json:"resourceCluster,omitempty"`
	ResourceNamespace string `json:"resourceNamespace,omitempty"`
	ResourceType      string `json:"resourceType"`
	SnapshotKey       string `json:"snapshotKey"`
}

// SnapshotRestoreStatus defines the observed state of SnapshotRestore
type SnapshotRestoreStatus struct {
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

// SnapshotRestore is the Schema for the openmcpsnapshotrestores API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=openmcpsnapshotrestores,scope=Namespaced
type SnapshotRestore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SnapshotRestoreSpec   `json:"spec,omitempty"`
	Status SnapshotRestoreStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SnapshotRestoreList contains a list of SnapshotRestore
type SnapshotRestoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SnapshotRestore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SnapshotRestore{}, &SnapshotRestoreList{})
}
