package resourceinfo

import (
//"k8s.io/klog"
)

var (
	emptyResource = Resource{}
)

// cluster Level
type ClusterInfo struct {
	ClusterName         string
	RequestedResource   *Resource
	AllocatableResource *Resource

	NodeList map[string]*NodeInfo
}

// node Level
type NodeInfo struct {
	NodeName            string
	RequestedResource   *Resource
	AllocatableResource *Resource
}

type Resource struct {
	MilliCPU int64
	Memory   int64
	Storage  int64
	Network  int64
}
