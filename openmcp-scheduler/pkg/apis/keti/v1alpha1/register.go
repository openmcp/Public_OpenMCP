// NOTE: Boilerplate only.  Ignore this file.

// Package v1alpha1 contains API Schema definitions for the keti v1alpha1 API group
// +k8s:deepcopy-gen=package,register
// +groupName=keti.example.com
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/runtime/scheme"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "keti.example.com", Version: "v1alpha1"}
	SchemeGroupVersion2 = schema.GroupVersion{Group: "autoscaling.k8s.io", Version: "v1beta2"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
	SchemeBuilder2 = &scheme.Builder{GroupVersion: SchemeGroupVersion2}
)
