package predicates

import (
	ketiresource "openmcp/openmcp/openmcp-scheduler/src/resourceinfo"
	"openmcp/openmcp/util/clusterManager"
)

type CheckNeededResources struct{}

func (pl *CheckNeededResources) Name() string {
	return "CheckNeededResources"
}

// Return true if there is at least 1 node that have AdditionalResources
func (pl *CheckNeededResources) Filter(newPod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, cm *clusterManager.ClusterManager) bool {

	// Node must have all of the additional resource
	// Example of *.yaml for a new OpenMCPDeployemt as folllow:
	//     resource:
	//       request:
	//         nvidia.com/gpu: 1
	//         amd.com/gpu: 1
	// In this case, selected node must have both of "nvidia.com/gpu, amd.com/gpu"

	if len(newPod.AdditionalResource) == 0 {
		return true
	}

	for _, node := range clusterInfo.Nodes {
		node_result := true
		for _, resource := range newPod.AdditionalResource {
			if contains(node.AdditionalResource, resource) == false {
				node_result = false
				break
			}
		}

		if node_result == true {
			return true
		}
	}

	return false
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if str == a {
			return true
		}
	}
	return false
}
