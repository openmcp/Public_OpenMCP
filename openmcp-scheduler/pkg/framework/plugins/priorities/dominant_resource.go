package priorities

import (
	"math"
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
)

type DominantResource struct {
	prescoring map[string]int64

	betweenScore int64
}

func (pl *DominantResource) Name() string {
	return "DominantResource"
}

func (pl *DominantResource) PreScore(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, check bool) int64 {
	dominantShareArr := make([]float64, 0)
	var clusterScore int64

	for _, node := range clusterInfo.Nodes {
		dominantShare := float64(0)

		// get Dominant share
		tmp := (float64(node.RequestedResource.MilliCPU) / float64(node.CapacityResource.MilliCPU)) * 100
		math.Max(dominantShare, tmp)

		tmp = (float64(node.RequestedResource.Memory) / float64(node.CapacityResource.Memory)) * 100
		math.Max(dominantShare, tmp)

		tmp = (float64(node.RequestedResource.EphemeralStorage) / float64(node.CapacityResource.EphemeralStorage)) * 100
		math.Max(dominantShare, tmp)

		dominantShareArr = append(dominantShareArr, dominantShare)
		nodeScore := int64(math.Round((1/getMinDominantShare(dominantShareArr))*math.MaxFloat64) * float64(maxScore))

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
		// omcplog.V(0).Infof("QOS후", pl.prescoring[clusterInfo.ClusterName])
		//omcplog.V(0).Infof("노드차이", pl.betweenScore)
	}
	return clusterScore
}

func (pl *DominantResource) Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, replicas int32, clustername string) int64 {
	dominantShareArr := make([]float64, 0)
	var clusterScore int64

	for _, node := range clusterInfo.Nodes {
		dominantShare := float64(0)

		// get Dominant share
		tmp := (float64(node.RequestedResource.MilliCPU) / float64(node.CapacityResource.MilliCPU)) * 100
		math.Max(dominantShare, tmp)

		tmp = (float64(node.RequestedResource.Memory) / float64(node.CapacityResource.Memory)) * 100
		math.Max(dominantShare, tmp)

		tmp = (float64(node.RequestedResource.EphemeralStorage) / float64(node.CapacityResource.EphemeralStorage)) * 100
		math.Max(dominantShare, tmp)

		dominantShareArr = append(dominantShareArr, dominantShare)
		nodeScore := int64(math.Round((1/getMinDominantShare(dominantShareArr))*math.MaxFloat64) * float64(maxScore))

		node.NodeScore = nodeScore
		clusterScore += nodeScore
	}

	return clusterScore
}

func getMinDominantShare(arr []float64) float64 {
	min := math.MaxFloat64

	for _, a := range arr {
		if a == 0 {
			continue
		}
		min = math.Min(min, a)
	}
	return min
}
