package v1alpha1

import (
	//appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OpenMCPServiceDNSRecordSpec defines the desired state of OpenMCPServiceDNSRecord
// +k8s:openapi-gen=true
type OpenMCPServiceDNSRecordSpec struct {
	// DomainRef is the name of the domain object to which the corresponding federated service belongs
	DomainRef string `json:"domainRef"`
	// RecordTTL is the TTL in seconds for DNS records created for this Service, if omitted a default would be used
	RecordTTL TTL `json:"recordTTL,omitempty"`
	// DNSPrefix when specified, an additional DNS record would be created with <DNSPrefix>.<KubeFedDomain>
	DNSPrefix string `json:"dnsPrefix,omitempty"`
	// ExternalName when specified, replaces the service name portion of a resource record
	// with the value of ExternalName.
	ExternalName string `json:"externalName,omitempty"`
	// AllowServiceWithoutEndpoints allows DNS records to be written for Service shards without endpoints
	AllowServiceWithoutEndpoints bool `json:"allowServiceWithoutEndpoints,omitempty"`
}
type ClusterDNS struct {
	// Cluster name
	Cluster string `json:"cluster,omitempty"`
	// LoadBalancer for the corresponding service
	LoadBalancer corev1.LoadBalancerStatus `json:"loadBalancer,omitempty"`
	// Zones to which the cluster belongs
	Zones []string `json:"zones,omitempty"`
	// Region to which the cluster belongs
	Region string `json:"region,omitempty"`
}

// SyncStatus defines the observed state of Sync
// +k8s:openapi-gen=true
type OpenMCPServiceDNSRecordStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Domain string       `json:"domain,omitempty"`
	DNS    []ClusterDNS `json:"dns,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Sync is the Schema for the syncs API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type OpenMCPServiceDNSRecord struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenMCPServiceDNSRecordSpec   `json:"spec,omitempty"`
	Status OpenMCPServiceDNSRecordStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SyncList contains a list of Sync
type OpenMCPServiceDNSRecordList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenMCPServiceDNSRecord `json:"items"`
}

//func init() {
//	SchemeBuilder.Register(&Sync{}, &SyncList{})
//}
