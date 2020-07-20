package predicates

import (
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
)

type MatchClusterAffinity struct{}

// Name returns name of the plugin
func (pl *MatchClusterAffinity) Name() string {
	return "MatchClusterAffinity"
}

func (pl *MatchClusterAffinity) Filter(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) bool {

	// check all nodes in this cluster
	for _, node := range clusterInfo.Nodes {
		result := []bool{}

		// check node's affinity value corresponding pod's affinity value
		for key, pod_values := range pod.Affinity{
			if node_value, ok := node.Affinity[key]; ok{
				// if node's affinity has pod's affinity, new deployment can be deploymented
				if containsAffinity(pod_values, node_value) == true{
					result = append(result, true) 
				}else{
					result = append(result, false)
				}
			} else {
				result = append(result, false)
			}
		}
		if !containsBoolean(result, false){
			return true
		}
	}

	return false
}

func containsAffinity(arr []string, str string) bool {
	for _, a := range arr {
		if str == a {
			return true
		}
	}
	return false
}

func containsBoolean(list []bool, it bool) bool {
	for _, l := range list {
		if l == it {
			return true
		}
	}
	return false
}
