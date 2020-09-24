package resourceinfo

import (
	corev1 "k8s.io/api/core/v1"
)

var (
	emptyResource = Resource{}
)

// cluster Level
type Cluster struct {
	ClusterName         string
	Nodes               []*NodeInfo
	RequestedResource   *Resource
	AllocatableResource *Resource
	PreFilter           bool
	PreFilterA          bool
}

// NodeInfo is node level aggregated information.
type NodeInfo struct {
	// Overall node information.
	ClusterName string
	NodeName    string

	Node *corev1.Node
	Pods []*Pod

	// Capacity
	CapacityResource *Resource
	// Total requested resource of all pods on this node
	RequestedResource *Resource
	// Total allocatable resource of all pods on this node
	AllocatableResource *Resource
	// Additional resource like nvidia/gpu
	AdditionalResource []string
	// Affinity(Region/Zone)
	Affinity map[string]string
	// Score to Update Resourcese
	NodeScore int64
	UpdateTX  int64
	UpdateRX  int64
	//if PreFilter is true, return Nodeis false

}

type Pod struct {
	// Overall pod informtation.
	ClusterName string
	NodeName    string
	PodName     string

	Pod                *corev1.Pod
	RequestedResource  *Resource
	AdditionalResource []string
	Affinity           map[string][]string
}

type Resource struct {
	MilliCPU         int64
	Memory           int64
	EphemeralStorage int64
}

func NewResource() *Resource {
	return &Resource{
		MilliCPU:         0,
		Memory:           0,
		EphemeralStorage: 0,
	}
}

func AddResources(res, new *Resource) *Resource {
	return &Resource{
		MilliCPU:         res.MilliCPU + new.MilliCPU,
		Memory:           res.Memory + new.Memory,
		EphemeralStorage: res.EphemeralStorage + new.EphemeralStorage,
	}
}

func GetAllocatable(capacity, request *Resource) *Resource {
	return &Resource{
		MilliCPU:         capacity.MilliCPU - request.MilliCPU,
		Memory:           capacity.Memory - request.Memory,
		EphemeralStorage: capacity.EphemeralStorage - request.EphemeralStorage,
	}
}
