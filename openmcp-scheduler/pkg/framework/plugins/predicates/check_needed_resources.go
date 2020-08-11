package predicates

import (
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
)

type CheckNeededResources struct{}

// Name returns name of the plugin
func (pl *CheckNeededResources) Name() string {
	return "CheckNeededResources"
}

func (pl *CheckNeededResources) Filter(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) bool {
	// Need datas about hardware spec
	var result bool

	if len(pod.IsNeedResourceMap) == 0 {
		return true
	}

	for _, node := range clusterInfo.Nodes {

		for key, pod_value := range pod.IsNeedResourceMap {
			if node_value, ok := node.IsNeedResourceMap[key]; ok{
				result = (pod_value == node_value)
				if !result {
					break
				}
			} else {
				// if pod doesn't need this resource, node may not have resource
				// if pod need this resource, node must have this resource
				if pod_value == false{
					result = true
				}else {
					result = false
				}
			}
		}

		// If ther is a cluster can be deployed a new Pod, return true
		if result == true{
			break
		}
	}

	return result
}
