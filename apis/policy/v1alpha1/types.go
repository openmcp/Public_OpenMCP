package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type OpenMCPPolicyTemplate struct {
	Spec OpenMCPPolicyTemplateSpec `json:"spec"`
}

type OpenMCPPolicyTemplateSpec struct {
	TargetController OpenMCPPolicyTartgetController `json:"targetController"`
	Policies         []OpenMCPPolicies              `json:"policies"`
}

type OpenMCPPolicyTartgetController struct {
	Kind string `json:"kind"`
}

type OpenMCPPolicies struct {
	Type  string   `json:"type"`
	Value []string `json:"value"`
}

type OpenMCPPolicySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	//Template - 생성
	Template OpenMCPPolicyTemplate `json:"template"`
	/*Template struct {
		Spec struct {
			TargetController struct {
				Kind string `json:"kind"`
			} `json:"targetController"`
			Policies []struct {
				Type string `json:"type"`
				Value string `json:"value"`
			} `json:"policies"`
		} `json:"spec"`
	} `json:"template"`*/
	RangeOfApplication string `json:"rangeOfApplication"`
	PolicyStatus       string `json:"policyStatus"`
	//Placement

}

// OpenMCPPolicyStatus defines the observed state of OpenMCPPolicy
// +k8s:openapi-gen=true
type OpenMCPPolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Replicas    int32            `json:"replicas"`
	ClusterMaps map[string]int32 `json:"clusters"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPPolicy is the Schema for the openmcppolicys API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type OpenMCPPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenMCPPolicySpec   `json:"spec,omitempty"`
	Status OpenMCPPolicyStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPPolicyList contains a list of OpenMCPPolicy
type OpenMCPPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenMCPPolicy `json:"items"`
}
