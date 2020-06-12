package predicates

import (
	// "k8s.io/klog"
	ketiresource "openmcpscheduler/pkg/controller/resourceinfo"
	"k8s.io/apimachinery/pkg/labels"
	// ketiframework "openmcpscheduler/pkg/controller/framework/v1alpha1"
)

type MatchClusterSelector struct{}

// var _ ketiframework.OpenmcpFilterPlugin = &Fit{}

// Name is the name of the plugin used in the plugin
// const Name = "MatchNodeSelector"

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

		if selector.Matches(labels.Set(node.Node.Labels)){
			// check labels & nodeselector
			// klog.Infof("")
			return true
		}
	}

	return false
}
// func New() ketiframework.OpenmcpPlugin {
// 	return &Fit{}
// }