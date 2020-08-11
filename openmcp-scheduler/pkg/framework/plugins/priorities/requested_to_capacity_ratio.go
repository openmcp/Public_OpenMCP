package priorities

import (
	"k8s.io/klog"
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
)

type RequestedToCapacityRatio struct{}

type FunctionShape []FunctionShapePoint

type FunctionShapePoint struct {
	// Utilization is function argument
	Utilization int64
	// Score is function value
	Score int64
}

var (
	// give priority to least utilized nodes by default
	defaultFunctionShape = NewFunctionShape([]FunctionShapePoint{
		{
			Utilization: 0,
			Score:       minScore,
		},
		{
			Utilization: 100,
			Score:       maxScore,
		},
	})
)

const (
	minUtilization = 0
	maxUtilization = 100
)

func (pl *RequestedToCapacityRatio) Name() string {
	return "RequestedToCapacityRatio"
}

func (pl *RequestedToCapacityRatio) Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) int64 {
	var clusterScore int64

	for _, node := range clusterInfo.Nodes {
		requested := node.RequestedResource.MilliCPU + pod.RequestedResource.MilliCPU
		nodeScore := RunRequestedToCapacityRatioScorerFunction(node.CapacityResource.MilliCPU, requested)
		requested = node.RequestedResource.Memory + pod.RequestedResource.Memory
		nodeScore += RunRequestedToCapacityRatioScorerFunction(node.CapacityResource.Memory, requested)
		requested = node.RequestedResource.EphemeralStorage + pod.RequestedResource.EphemeralStorage
		nodeScore += RunRequestedToCapacityRatioScorerFunction(node.CapacityResource.EphemeralStorage, requested)

		node.NodeScore = nodeScore
		clusterScore += nodeScore
	}

	return clusterScore
}

func RunRequestedToCapacityRatioScorerFunction (capacity, requested int64) int64 {
	scoringFunctionShape := defaultFunctionShape
	rawScoringFunction := buildBrokenLinearFunction(scoringFunctionShape)

	// 𝑠𝑐𝑜𝑟𝑒=𝑠𝑐𝑜𝑟𝑒+(100−((𝑛𝑜𝑑𝑒.𝑎𝑙𝑙𝑜𝑐𝑎𝑏𝑙𝑒.𝐶𝑃𝑈−(𝑝𝑜𝑑.𝑟𝑒𝑞𝑢𝑒𝑠𝑡.𝐶𝑃𝑈+𝑝𝑜𝑑.𝑢𝑠𝑒𝑑.𝐶𝑃𝑈))/𝑛𝑜𝑑𝑒.𝑎𝑙𝑙𝑜𝑐𝑎𝑏𝑙𝑒.𝐶𝑃𝑈))
	resourceScoringFunction := func(requested, capacity int64) int64 {
		if capacity == 0 || requested > capacity {
			return rawScoringFunction(maxUtilization)
		}
		return rawScoringFunction(maxUtilization - (capacity-requested)*maxUtilization/capacity)
	}

	return int64(resourceScoringFunction(requested, capacity))
}


func buildBrokenLinearFunction(shape FunctionShape) func(int64) int64 {
	n := len(shape)
	return func(p int64) int64 {
		for i := 0; i < n; i++ {
			if p <= shape[i].Utilization {
				if i == 0 {
					return shape[0].Score
				}
				return shape[i-1].Score + (shape[i].Score-shape[i-1].Score)*(p-shape[i-1].Utilization)/(shape[i].Utilization-shape[i-1].Utilization)
			}
		}
		return shape[n-1].Score
	}
}

func NewFunctionShape(points []FunctionShapePoint) FunctionShape {
	n := len(points)

	if n == 0 {
		klog.Info("at least one point must be specified")
		return nil
	}

	for i := 1; i < n; i++ {
		if points[i-1].Utilization >= points[i].Utilization {
			klog.Infof("utilization values must be sorted. Utilization[%v]==%v >= Utilization[%v]==%v", i-1, points[i-1].Utilization, i, points[i].Utilization)
			return nil
		}
	}

	for i, point := range points {
		if point.Utilization < minUtilization {
			klog.Infof("utilization values must not be less than %v. Utilization[%v]==%v", minUtilization, i, point.Utilization)
			return nil
		}
		if point.Utilization > maxUtilization {
			klog.Infof("utilization values must not be greater than %v. Utilization[%v]==%v", maxUtilization, i, point.Utilization)
			return nil
		}
		if point.Score < minScore {
			klog.Infof("score values must not be less than %v. Score[%v]==%v", minScore, i, point.Score)
			return nil
		}
		if point.Score > maxScore {
			klog.Infof("score valuses not be greater than %v. Score[%v]==%v", maxScore, i, point.Score)
			return nil
		}
	}

	// We make defensive copy so we make no assumption if array passed as argument is not changed afterwards
	pointsCopy := make(FunctionShape, n)
	copy(pointsCopy, points)
	return pointsCopy
}
