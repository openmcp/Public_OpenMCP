// +build !ignore_autogenerated

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPDNSEndpoint":       schema_pkg_apis_keti_v1alpha1_OpenMCPDNSEndpoint(ref),
		"openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPDNSEndpointSpec":   schema_pkg_apis_keti_v1alpha1_OpenMCPDNSEndpointSpec(ref),
		"openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPDNSEndpointStatus": schema_pkg_apis_keti_v1alpha1_OpenMCPDNSEndpointStatus(ref),

		"openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPIngressDNSRecord":       schema_pkg_apis_keti_v1alpha1_OpenMCPIngressDNSRecord(ref),
		"openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPIngressDNSRecordSpec":   schema_pkg_apis_keti_v1alpha1_OpenMCPIngressDNSRecordSpec(ref),
		"openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPIngressDNSRecordStatus": schema_pkg_apis_keti_v1alpha1_OpenMCPIngressDNSRecordStatus(ref),

		"openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPServiceDNSRecord":       schema_pkg_apis_keti_v1alpha1_OpenMCPServiceDNSRecord(ref),
		"openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPServiceDNSRecordSpec":   schema_pkg_apis_keti_v1alpha1_OpenMCPServiceDNSRecordSpec(ref),
		"openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPServiceDNSRecordStatus": schema_pkg_apis_keti_v1alpha1_OpenMCPServiceDNSRecordStatus(ref),

		"openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPDomain": schema_pkg_apis_keti_v1alpha1_OpenMCPDomain(ref),
	}
}

func schema_pkg_apis_keti_v1alpha1_OpenMCPDNSEndpoint(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "OpenMCPDNSEndpoint is the Schema for the OpenMCPDNSEndpoints API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPDNSEndpointSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPDNSEndpointStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta", "openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPDNSEndpointSpec", "openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPDNSEndpointStatus"},
	}
}

func schema_pkg_apis_keti_v1alpha1_OpenMCPDNSEndpointSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "OpenMCPDNSEndpointSpec defines the desired state of OpenMCPDNSEndpoint",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_keti_v1alpha1_OpenMCPDNSEndpointStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "OpenMCPDNSEndpointStatus defines the observed state of OpenMCPDNSEndpoint",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_keti_v1alpha1_OpenMCPServiceDNSRecord(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "OpenMCPServiceDNSRecord is the Schema for the OpenMCPServiceDNSRecords API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPServiceDNSRecordSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPServiceDNSRecordStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta", "openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPServiceDNSRecordSpec", "openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPServiceDNSRecordStatus"},
	}
}

func schema_pkg_apis_keti_v1alpha1_OpenMCPServiceDNSRecordSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "OpenMCPServiceDNSRecordSpec defines the desired state of OpenMCPServiceDNSRecord",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_keti_v1alpha1_OpenMCPServiceDNSRecordStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "OpenMCPServiceDNSRecordStatus defines the observed state of OpenMCPServiceDNSRecord",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_keti_v1alpha1_OpenMCPIngressDNSRecord(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "OpenMCPIngressDNSRecord is the Schema for the OpenMCPIngressDNSRecords API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPIngressDNSRecordSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPIngressDNSRecordStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta", "openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPIngressDNSRecordSpec", "openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPIngressDNSRecordStatus"},
	}
}

func schema_pkg_apis_keti_v1alpha1_OpenMCPIngressDNSRecordSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "OpenMCPIngressDNSRecordSpec defines the desired state of OpenMCPIngressDNSRecord",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_keti_v1alpha1_OpenMCPIngressDNSRecordStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "OpenMCPIngressDNSRecordStatus defines the observed state of OpenMCPIngressDNSRecord",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_keti_v1alpha1_OpenMCPDomain(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "OpenMCPDomain is the Schema for the OpenMCPDomains API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPDomainSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPDomainStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta", "openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPDomainSpec", "openmcp-dns-controller/pkg/apis/keti/v1alpha1.OpenMCPDomainStatus"},
	}
}
