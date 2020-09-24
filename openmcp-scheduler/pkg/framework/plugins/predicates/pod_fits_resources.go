package predicates

import (
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
)

type PodFitsResources struct{}

func (pl *PodFitsResources) Name() string {
	return "PodFitsResources"
}
func (pl *PodFitsResources) PreFilter(newPod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) bool {
	for _, node := range clusterInfo.Nodes {
		// check if node has enough CPU
		if node.AllocatableResource.MilliCPU < newPod.RequestedResource.MilliCPU {
			// omcplog.V(0).Info(clusterInfo.ClusterName + "True")
			clusterInfo.PreFilter = false
			return true
		}

	}
	clusterInfo.PreFilter = true

	return false

}

func (pl *PodFitsResources) Filter(newPod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) bool {

	for _, node := range clusterInfo.Nodes {
		node_result := true
		// check if node has enough Memory
		if node.AllocatableResource.Memory < newPod.RequestedResource.Memory {
			node_result = false
		}
		// check if node has enough EphemeralStorage
		if node.AllocatableResource.EphemeralStorage < newPod.RequestedResource.EphemeralStorage {
			node_result = false
		}

		if node_result == true {
			return true
		}
	}
	return false
}

// Return true if there is at least 1 node that have AdditionalResources
func (pl *PodFitsResources) PostFilter(newPod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) (bool, error) {

	for _, node := range clusterInfo.Nodes {
		node_result := true
		if node.CapacityResource.MilliCPU < newPod.RequestedResource.MilliCPU {
			node_result = false
		}
		if node.CapacityResource.Memory < newPod.RequestedResource.Memory {
			node_result = false
		}
		// check if node has enough Memory
		if node.CapacityResource.Memory < newPod.RequestedResource.Memory {
			node_result = false
		}
		// check if node has enough EphemeralStorage
		if node.CapacityResource.EphemeralStorage < newPod.RequestedResource.EphemeralStorage {
			node_result = false
		}

		if node_result == true {
			return false, nil
		}
	}
	return true, nil
}
