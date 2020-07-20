package predicates

import (
	v1 "k8s.io/api/core/v1"
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
)

type PodFitsHostPorts struct{}

// Name returns name of the plugin
func (pl *PodFitsHostPorts) Name() string {
	return "PodFitsHostPorts"
}

func (pl *PodFitsHostPorts) Filter(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) bool {
	var wantPorts []*v1.ContainerPort
	wantPorts = getContainerPorts(pod.Pod)
	if len(wantPorts) == 0{
		return true
	}

	for _, node := range clusterInfo.Nodes {
		canDeploy := true
		existingPorts := node.UsedPorts

		for _, wantPort := range wantPorts {
			contain := containsPort(existingPorts, wantPort)
			if contain == true {
				canDeploy = false
				break
			}
		}

		if canDeploy == false{
			return false
		}
	}

	return true
}

func getContainerPorts(pods ...*v1.Pod) []*v1.ContainerPort {
	var ports []*v1.ContainerPort
	for _, pod := range pods {
		for j := range pod.Spec.Containers {
			container := &pod.Spec.Containers[j]
			for k := range container.Ports {
				ports = append(ports, &container.Ports[k])
			}
		}
	}
	return ports
}

func containsPort(arr []*v1.ContainerPort , des *v1.ContainerPort) bool{
	for _, a := range arr {
		if des.ContainerPort == a.ContainerPort {
			return true
		}
	}
	return false
}
