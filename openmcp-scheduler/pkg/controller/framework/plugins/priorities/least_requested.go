package priorities

import (
	ketiresource "openmcpscheduler/pkg/controller/resourceinfo"
)

type LeastRequested struct{}

func (pl *LeastRequested) Name() string {
	return "LeastRequested"
}

func (pl *LeastRequested) Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) int64 {
	var clutserScore int64

	for _, node := range clusterInfo.Nodes {
		clutserScore += leastRequestedScore(pod.RequestedResource.MilliCPU, node.AllocatableResource.MilliCPU)
		clutserScore += leastRequestedScore(pod.RequestedResource.Memory, node.AllocatableResource.Memory)
		clutserScore += leastRequestedScore(pod.RequestedResource.EphemeralStorage, node.AllocatableResource.EphemeralStorage)
	}


	return clutserScore
}

func leastRequestedScore(requested, allocable int64) int64 {
	if allocable == 0 {
		return 0
	}
	if requested > allocable {
		return 0
	}
	return ((allocable - requested) * int64(100)) / allocable
}