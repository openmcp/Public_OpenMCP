package priorities

import (
	ketiresource "openmcpscheduler/pkg/controller/resourceinfo"
)

type MostRequested struct{}

func (pl *MostRequested) Name() string {
	return "MostRequested"
}

func (pl *MostRequested) Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) int64 {
	var clutserScore int64

	for _, node := range clusterInfo.Nodes {
		clutserScore += mostRequestedScore(pod.RequestedResource.MilliCPU, node.AllocatableResource.MilliCPU)
		clutserScore += mostRequestedScore(pod.RequestedResource.Memory, node.AllocatableResource.Memory)
		clutserScore += mostRequestedScore(pod.RequestedResource.EphemeralStorage, node.AllocatableResource.EphemeralStorage)
	}

	return clutserScore
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