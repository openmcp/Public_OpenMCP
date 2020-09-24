package priorities

import (
	//"openmcp/openmcp/omcplog"

	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
)

type BalancedNetworkAllocation struct {
	prescoring map[string]int64

	betweenScore int64
}

// func (pl *BalancedNetworkAllocation) PreScore(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) {
// 	var nodeScore int64

// 	for _, node := range clusterInfo.Nodes {
// 		if node.UpdateRX == 0 && node.UpdateTX == 0 {
// 			nodeScore = maxScore
// 		} else {
// 			nodeScore = int64((1 / float64(node.UpdateRX+node.UpdateTX)) * float64(maxScore))
// 		}
// 		//omcplog.V(0).Infof("[%v] node rx [%d] tx [%d]", node.NodeName, node.UpdateRX, node.UpdateTX)
// 		node.NodeScore = nodeScore
// 	}

// }

func (pl *BalancedNetworkAllocation) Name() string {

	return "BalancedNetworkAllocation"
}

func (pl *BalancedNetworkAllocation) PreScore(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, check bool) int64 {
	// startTime := time.Now()
	var clusterScore int64
	clusterScore = 0

	for _, node := range clusterInfo.Nodes {
		clusterScore += node.NodeScore
	}
	// OelapsedTime := time.Since(startTime)
	if !check {
		if len(pl.prescoring) == 0 {
			pl.prescoring = make(map[string]int64)
		}
		pl.prescoring[clusterInfo.ClusterName] = clusterScore
	} else {

		pl.betweenScore = pl.prescoring[clusterInfo.ClusterName] - int64(clusterScore)
		pl.prescoring[clusterInfo.ClusterName] = clusterScore

		//omcplog.V(0).Infof("노드차이", pl.betweenScore)
	}
	return clusterScore

}
func (pl *BalancedNetworkAllocation) Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, replicas int32, clustername string) int64 {
	if clustername == clusterInfo.ClusterName {
		score := pl.prescoring[clusterInfo.ClusterName] - pl.betweenScore
		if score < 0 {
			return 0
		}
	}
	score := pl.prescoring[clusterInfo.ClusterName]
	return score

}
