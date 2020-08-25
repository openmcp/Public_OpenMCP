package priorities

import (
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
)

type MostRequested struct{}

func (pl *MostRequested) Name() string {
	return "MostRequested"
}

func (pl *MostRequested) Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) int64 {
	var clusterScore int64

	for _, node := range clusterInfo.Nodes {
		nodeScore := mostRequestedScore(pod.RequestedResource.MilliCPU, node.AllocatableResource.MilliCPU)
		nodeScore += mostRequestedScore(pod.RequestedResource.Memory, node.AllocatableResource.Memory)
		nodeScore += mostRequestedScore(pod.RequestedResource.EphemeralStorage, node.AllocatableResource.EphemeralStorage)

		node.NodeScore = nodeScore
		clusterScore += nodeScore
	}

	return clusterScore
}

func mostRequestedScore(requested, allocable int64) int64 {
	if allocable == 0 {
		return 0
	}
	if requested > allocable {
		return 0
	}
	return (requested * int64(100)) / allocable
}
