package predicates

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
	"openmcp/openmcp/omcplog"
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
)

type PodFitsHostPorts struct{}

const (
	// ErrReason when cluster ports aren't available
	ErrReason = "cluster(s) didn't have free ports for the requested deployment ports"
)

func (pl *PodFitsHostPorts) Name() string {
	return "PodFitsHostPorts"
}

func (pl *PodFitsHostPorts) Filter(newPod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) bool {
	// Requested ports should be available for newPod
	// Example of *.yaml for a new OpenMCPDeployemt as folllow:
	//   spec:
	//     containers:
	//       ports:
	//         - name: http
	//           containerPort:80
	//         - name: health
	//           containerPort:8080
	// In this case, selected node must be able to use port "80" and "8080"

	// if there is not requested Ports, PodFitHostPorts return true
	wantPorts := getContainerPorts(newPod.Pod)

	if len(wantPorts) == 0 {
		return true
	}

	// It checks all nodes in clusterInfo and returns true if any node has an available port
	for _, node := range clusterInfo.Nodes {

		// newPod cannot be deployed, if one of the wantPorts is used
		if fitsPorts(wantPorts, node) == false {
			omcplog.V(0).Info("%s", ErrReason)
			return false
		}
	}

	return true
}

func fitsPorts(wantPorts []*v1.ContainerPort, nodeInfo *ketiresource.NodeInfo) bool {

	for _, wantPort := range wantPorts {

		// Checks if the wantPort conflict with the existing ones in HostPortInfo
		if wantPort.HostPort <= 0 {
			continue
		}

		for i := range nodeInfo.Pods {
			pod := nodeInfo.Pods[i]

			for j := range pod.Pod.Spec.Containers {
				container := &pod.Pod.Spec.Containers[j]

				for k := range container.Ports {
					existingPort := &container.Ports[k]

					if existingPort.HostPort == wantPort.HostPort {
						klog.Infof("existingPort:%v, wantPort:%v", existingPort, wantPort)
						return false
					}
				}
			}
		}
	}

	return true
}

// getContainerPorts returns the used host ports of Pods from newPod's spec
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
