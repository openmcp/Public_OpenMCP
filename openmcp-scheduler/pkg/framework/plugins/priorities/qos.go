package priorities

import (
	v1 "k8s.io/api/core/v1"
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
)

type QosPriority struct{}

const (
	minScore int64 = 0
	midScore int64 = (minScore + maxScore) / 2
	maxScore int64 = 10
)

func (pl *QosPriority) Name() string {
	return "QosPriority"
}

func (pl *QosPriority) Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) int64 {
	var clusterScore int64

	for _, node := range clusterInfo.Nodes {
		var nodeScore int64
		for _, pod := range node.Pods {

			// get PodQOSClass from v1.Pod
			qos := pod.Pod.Status.QOSClass

			switch qos{
			case v1.PodQOSGuaranteed:
				nodeScore += minScore
			case v1.PodQOSBurstable:
				nodeScore += midScore
			case v1.PodQOSBestEffort:
				nodeScore += maxScore
			}
		}

		node.NodeScore = nodeScore
		clusterScore += nodeScore
	}

	return clusterScore
}
