package resourceinfo

import (
	v1 "k8s.io/api/core/v1"
)

var (
	emptyResource = Resource{}
)

// cluster Level
type Cluster struct {
	ClusterName			string
	Nodes				[]*NodeInfo
	RequestedResource 	*Resource
	AllocatableResource	*Resource
}

// NodeInfo is node level aggregated information.
type NodeInfo struct {
	// Overall node information.
	Node				*v1.Node

	NodeName			string
	ClusterName			string
	Pods				[]*Pod
	UsedPorts			[]*v1.ContainerPort

	// Capacity
	CapacityResource	*Resource
	// Total requested resource of all pods on this node
	RequestedResource 	*Resource
	// Total allocatable resource of all pods on this node
	AllocatableResource	*Resource
	// For Hardware spec
	IsNeedResourceMap		map[string]bool
	// Affinity -> will be change to GeoAffinity
	Affinity			map[string]string
	// Score to Update Resourcese
	NodeScore			int64
}

type Pod struct {
	// Overall pod informtation.
	Pod					*v1.Pod
	PodName				string
	NodeName			string
	ClusterName			string
	RequestedResource 	*Resource
	Affinity			map[string][]string
	IsNeedResourceMap	map[string]bool
}

type Resource struct {
	MilliCPU				int64
	Memory					int64
	EphemeralStorage		int64
	Network					int64
}

type ProtocolPort struct {
	Protocol	string
	Port		int32
}

func NewResource() *Resource {
	return &Resource {
		MilliCPU:				0,
		Memory:					0,
		EphemeralStorage:		0,
		Network:				0,
	}
}

func AddResources(res, new *Resource) *Resource {
	return &Resource {
		MilliCPU:			res.MilliCPU + new.MilliCPU,
		Memory:				res.Memory + new.Memory,
		EphemeralStorage:	res.EphemeralStorage + new.EphemeralStorage,
	}
}

func GetAllocatable(capacity, request *Resource) *Resource {
	return &Resource {
		MilliCPU:			capacity.MilliCPU - request.MilliCPU,
		Memory:				capacity.Memory - request.Memory,
		EphemeralStorage:	capacity.EphemeralStorage - request.EphemeralStorage,
	}
}