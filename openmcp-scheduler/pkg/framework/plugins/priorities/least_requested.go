package priorities

import (
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
)

type LeastRequested struct{}

func (pl *LeastRequested) Name() string {
	return "LeastRequested"
}

func (pl *LeastRequested) Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) int64 {
	var clusterScore int64

	for _, node := range clusterInfo.Nodes {
		nodeScore := leastRequestedScore(pod.RequestedResource.MilliCPU, node.AllocatableResource.MilliCPU)
		nodeScore += leastRequestedScore(pod.RequestedResource.Memory, node.AllocatableResource.Memory)
		nodeScore += leastRequestedScore(pod.RequestedResource.EphemeralStorage, node.AllocatableResource.EphemeralStorage)

		node.NodeScore = nodeScore
		clusterScore += nodeScore
	}


	return clusterScore
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
