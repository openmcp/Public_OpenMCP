package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IngressDNSRecordSpec defines the desired state of IngressDNSRecord
type OpenMCPIngressDNSRecordSpec struct {
	// Host from the IngressRule in Cluster Ingress Spec
	// Hosts []string `json:"hosts,omitempty"`
	// RecordTTL is the TTL in seconds for DNS records created for the Ingress, if omitted a default would be used
	RecordTTL TTL `json:"recordTTL,omitempty"`
	// DomainRef string `json:"domainRef"`
}

// IngressDNSRecordStatus defines the observed state of IngressDNSRecord
type OpenMCPIngressDNSRecordStatus struct {
	// Array of Ingress Controller LoadBalancers
	DNS []ClusterIngressDNS `json:"dns,omitempty"`
	//Domain string       `json:"domain,omitempty"`
}

// ClusterIngressDNS defines the observed status of Ingress within a cluster.
type ClusterIngressDNS struct {
	// Cluster name
	Cluster string `json:"cluster,omitempty"`
	// LoadBalancer for the corresponding ingress controller
	LoadBalancer corev1.LoadBalancerStatus `json:"loadBalancer,omitempty"`
	Hosts []string `json:"host,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IngressDNSRecord
// +k8s:openapi-gen=true
// +kubebuilder:resource:path=ingressdnsrecords
// +kubebuilder:subresource:status
type OpenMCPIngressDNSRecord struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenMCPIngressDNSRecordSpec   `json:"spec,omitempty"`
	Status OpenMCPIngressDNSRecordStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IngressDNSRecordList contains a list of IngressDNSRecord
type OpenMCPIngressDNSRecordList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenMCPIngressDNSRecord `json:"items"`
}