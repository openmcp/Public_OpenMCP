package priorities

import (

	"math"
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
)

type BalancedNetworkAllocation struct {
	prescoring map[string]int64

	betweenScore int64
}


func (pl *BalancedNetworkAllocation) Name() string {

	return "BalancedNetworkAllocation"
}

func (pl *BalancedNetworkAllocation) PreScore(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, check bool) int64 {
	var clusterScore int64
	clusterScore = 0

	for _, node := range clusterInfo.Nodes {
		clusterScore += node.NodeScore
	}
	if !check {
		if len(pl.prescoring) == 0 {
			pl.prescoring = make(map[string]int64)
		}
		pl.prescoring[clusterInfo.ClusterName] = clusterScore
	} else {

		pl.betweenScore = pl.prescoring[clusterInfo.ClusterName] - int64(clusterScore)
		pl.betweenScore = int64(math.Abs(float64(pl.betweenScore)))
		pl.prescoring[clusterInfo.ClusterName] = clusterScore

	}
	return clusterScore

}
func (pl *BalancedNetworkAllocation) Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, replicas int32, clustername string) int64 {
	if clustername == clusterInfo.ClusterName {
		score := pl.prescoring[clusterInfo.ClusterName] - pl.betweenScore
		return score
	}
	score := pl.prescoring[clusterInfo.ClusterName]
	return score

}
