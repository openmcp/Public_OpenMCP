package predicates

import (
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
	"k8s.io/apimachinery/pkg/labels"
)

type MatchClusterSelector struct{}

// Name returns name of the plugin
func (pl *MatchClusterSelector) Name() string {
	return "MatchClusterSelector"
}

func (pl *MatchClusterSelector) Filter(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) bool {
	if len(pod.Pod.Spec.NodeSelector) == 0 {
		return true
	}

	for _, node := range clusterInfo.Nodes {
		selector := labels.SelectorFromSet(pod.Pod.Spec.NodeSelector)

		if (selector.Matches(labels.Set(node.Node.Labels))){
			return true
		}
	}
	
	return false
}
