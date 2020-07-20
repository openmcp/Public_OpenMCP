package predicates

import (
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
)

type PodFitsResources struct{}

// Name returns name of the plugin
func (pl *PodFitsResources) Name() string {
	return "PodFitsResources"
}

func (pl *PodFitsResources) Filter(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) bool {

	for _, node := range clusterInfo.Nodes{
		result := true

		// check CPU
		if node.AllocatableResource.MilliCPU < pod.RequestedResource.MilliCPU {
			result = result || false
			continue
		}
		// check Memory
		if node.AllocatableResource.Memory < pod.RequestedResource.Memory {
			result = result || false
			continue
		}
		// check Storage
		if node.AllocatableResource.EphemeralStorage < pod.RequestedResource.EphemeralStorage {
			result = result || false
			continue
		}

		if result == true{
			return true
		}
	}

	return false
}
