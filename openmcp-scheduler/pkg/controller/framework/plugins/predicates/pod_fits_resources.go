package predicates

import (
	// "k8s.io/klog"
	ketiresource "openmcpscheduler/pkg/controller/resourceinfo"
	// ketiframework "openmcpscheduler/pkg/controller/framework/v1alpha1"
)

type PodFitsResources struct{}

// var _ ketiframework.OpenmcpFilterPlugin = &Fit{}

// Name is the name of the plugin used in the plugin
// const Name = "PodFitsResources"

// Name returns name of the plugin
func (pl *PodFitsResources) Name() string {
	return "PodFitsResources"
}

func (pl *PodFitsResources) Filter(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) bool {
	// klog.Infof("*********** [PodFitsResources] Filter ***********")

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

// func New() ketiframework.OpenmcpPlugin {
// 	return &Fit{}
// }