package predicates

import (
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
)

type PodFitsResources struct{}

func (pl *PodFitsResources) Name() string {
	return "PodFitsResources"
}

func (pl *PodFitsResources) Filter(newPod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) bool {

	for _, node := range clusterInfo.Nodes{
		node_result := true

		// check if node has enough CPU
		if node.AllocatableResource.MilliCPU < newPod.RequestedResource.MilliCPU {
			node_result = false
		}
		// check if node has enough Memory
		if node.AllocatableResource.Memory < newPod.RequestedResource.Memory {
			node_result = false
		}
		// check if node has enough EphemeralStorage
		if node.AllocatableResource.EphemeralStorage < newPod.RequestedResource.EphemeralStorage {
			node_result = false
		}

		if node_result == true{
			return true
		}
	}

	return false
}
