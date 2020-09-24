package priorities

import (
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"

	v1 "k8s.io/api/core/v1"
)

type QosPriority struct {
	prescoring   map[string]int64
	betweenScore int64
}

const (
	minScore int64 = 0
	midScore int64 = (minScore + maxScore) / 2
	maxScore int64 = 10
)

func (pl *QosPriority) Name() string {
	return "QosPriority"
}
func (pl *QosPriority) PreScore(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, check bool) int64 {
	var clusterScore int64

	for _, node := range clusterInfo.Nodes {
		var nodeScore int64
		for _, pod := range node.Pods {

			// get PodQOSClass from v1.Pod
			qos := pod.Pod.Status.QOSClass

			switch qos {
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
	if !check {
		if len(pl.prescoring) == 0 {
			pl.prescoring = make(map[string]int64)
		}
		pl.prescoring[clusterInfo.ClusterName] = clusterScore
	} else {
		// omcplog.V(0).Infof("QOS전", pl.prescoring[clusterInfo.ClusterName], clusterInfo.ClusterName)
		pl.betweenScore = pl.prescoring[clusterInfo.ClusterName] - int64(clusterScore)
		pl.prescoring[clusterInfo.ClusterName] = clusterScore

	}
	return clusterScore
}

func (pl *QosPriority) Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, replicas int32, clustername string) int64 {
	if clustername == clusterInfo.ClusterName {
		score := pl.prescoring[clusterInfo.ClusterName] - pl.betweenScore
		if score < 0 {
			return 0
		}
	}
	score := pl.prescoring[clusterInfo.ClusterName]
	return score
}
