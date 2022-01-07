package v1alpha1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type OpenMCPCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenMCPClusterSpec   `json:"spec,omitempty"`
	Status OpenMCPClusterStatus `json:"status,omitempty"`
}

type OpenMCPClusterSpec struct {
	ServerIP               string       `json:"serverIP" protobuf:"bytes,1,opt,name=serverip"`
	ClusterPlatformType    string       `json:"clusterPlatformType" protobuf:"bytes,1,opt,name=clusterplatformtype"`
	ClusterNetworkLocation string       `json:"clusterNetworkLocation" protobuf:"bytes,2,opt,name=clusternetworklocation"`
	JoinStatus             string       `json:"joinStatus" protobuf:"bytes,3,opt,name=joinstatus"`
	JoinStatusTime         time.Time    `json:"joinStatusTime" protobuf:"bytes,3,opt,name=joinstatustime"`
	UnJoinStatusTime       time.Time    `json:"unJoinStatusTime" protobuf:"bytes,3,opt,name=unjoinstatustime"`
	NodeInfo               NodeInfo     `json:"nodeInfo" protobuf:"bytes,4,opt,name=nodeinfo"`
	MetalLBRange           MetalLBRange `json:"metalLBRange" protobuf:"bytes,5,opt,name=metallbrange"`
	KubeconfigInfo         []byte       `json:"kubeconfigInfo,omitempty" protobuf:"bytes,6,opt,name=kubeconfiginfo"`
}

type NodeInfo struct {
	Region string `json:"region,omitempty"`
	Zone   string `json:"zone,omitempty"`
}

type MetalLBRange struct {
	AddressFrom string `json:"addressFrom,omitempty"`
	AddressTo   string `json:"addressTo,omitempty"`
}

/*type ClusterInfo struct {
	APIVersion     string      `json:"apiVersion,omitempty"`
	Clusters       Clusters  `json:"clusters,omitempty"`
	Contexts       Contexts  `json:"contexts,omitempty"`
	CurrentContext string      `json:"current-context,omitempty"`
	Kind           string      `json:"kind,omitempty"`
	Preferences    Preferences `json:"preferences,omitempty"`
	Users          Users      `json:"users,omitempty"`
}*/
/*
type Clusters struct {
	Cluster Cluster `json:"cluster,omitempty"`
	Name    string  `json:"name,omitempty"`
}

type Cluster struct {
	CertificateAuthorityData string `json:"certificate-authority-data,omitempty"`
	Server                   string `json:"server,omitempty"`
}

type Contexts struct {
	Context Context `json:"context,omitempty"`
	Name    string  `json:"name,omitempty"`
}

type Context struct {
	Cluster string `json:"cluster,omitempty"`
	User    string `json:"user,omitempty"`
}

type Preferences struct {
}

type Users struct {
	User    User    `json:"user,omitempty"`
	Name    string  `json:"name,omitempty"`
}

type User struct {
	ClientCertificateData string `json:"client-certificate-data,omitempty"`
	ClientKeyData         string `json:"client-key-data,omitempty"`
}
*/
type OpenMCPClusterStatus struct {
	//ClusterStatus string `json:"clusterStatus,omitempty"`
}

type OpenMCPClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenMCPCluster `json:"items"`
}
