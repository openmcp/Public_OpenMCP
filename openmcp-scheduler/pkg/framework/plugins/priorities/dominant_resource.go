package priorities

import (
	"math"
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
)

type DominantResource struct{}

func (pl *DominantResource) Name() string {
	return "DominantResource"
}

func (pl *DominantResource) Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) int64 {
	dominantShareArr := make([]float64, 0)

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
	}

	return int64(math.Round((1 / getMinDominantShare(dominantShareArr)) * math.MaxFloat64) * float64(maxScore))
}

func getMinDominantShare(arr []float64) float64 {
	min := math.MaxFloat64

	for _, a := range arr {
		if a == 0{
			continue
		}
		min = math.Min(min, a)
	}
	return min
}
