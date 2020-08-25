package predicates

import (
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
)

type MatchClusterAffinity struct{}

func (pl *MatchClusterAffinity) Name() string {
	return "MatchClusterAffinity"
}

func (pl *MatchClusterAffinity) Filter(newPod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) bool {

	// Node must have all of the additional resource
	// Examples of *.yaml for a new OpenMCPDeployemt as folllow:
	// # Example 01 #
	//   spec:
	//     affinity:
	//       region:
	//         -AS
	//         -EU
	//       zone:
	//         -KR
	//         -DE
	//         -PT
	// In this case, selected node must have "KR:AS" or "DE:EU" or "PT:EU"
	//
	// # Example 02 #
	//   spec:
	//     affinity:
	//       zone:
	//         -KR
	//         -CH
	// In this case, selected node must have "KR" or "CH"

	if len(newPod.Affinity) == 0 {
		return true
	}

	for _, node := range clusterInfo.Nodes {
		node_result := true

		for key, pod_values := range newPod.Affinity {

			// compare Node's Affinity and Pod's Affinity
			if node_value, ok := node.Affinity[key]; ok {

				// if node's affinity has pod's affinity, new deployment can be deploymented
				if contains(pod_values, node_value) == false {
					node_result = false
				}

			} else {
				node_result = false
			}

			if node_result == false {
				break
			}
		}

		if node_result == true {
			return true
		}
	}

	return false
}
