package priorities

import (
	"time"
	"k8s.io/klog"

	"context"
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
	"openmcp/openmcp/openmcp-scheduler/pkg/protobuf"
)

type BalancedNetworkAllocation struct{
	GRPCServer		protobuf.RequestAnalysisClient
}

func (pl *BalancedNetworkAllocation) Name() string {
	return "BalancedNetworkAllocation"
}

func (pl *BalancedNetworkAllocation) Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) int64 {
	var clutserScore int64
	klog.Infof("***** [Start] BalancedNetworkAllocation *****")
	startTime := time.Now()

	for _, node := range clusterInfo.Nodes {

		node_info := &protobuf.NodeInfo{ClusterName: clusterInfo.ClusterName, NodeName: node.NodeName}
		grpcClient := pl.GRPCServer

		grpc_starttime := time.Now()
		klog.Infof("***** [Start] gRPC Communicates *****")
		result, _ := grpcClient.SendNetworkAnalysis(context.TODO(), node_info)

		grpc_endtime := time.Since(grpc_starttime)
		klog.Infof("=> Total gRPC Time : %v", grpc_endtime)
		klog.Infof("***** [End] gRPC Communicates *****")

		var nodeScore int64
		rx := result.RX
		tx := result.TX

		if rx == 0 && tx == 0 {
			nodeScore = maxScore
		}else {
			nodeScore = int64((1 / float64(rx + tx)) * float64(maxScore))
		}
		node.NodeScore = nodeScore
		clutserScore += nodeScore
	}

	elapsedTime := time.Since(startTime)
	klog.Infof("=> Total BalancedNetworkAllocation Time : %v", elapsedTime)
	klog.Infof("***** [End] BalancedNetworkAllocation *****")

	return clutserScore
}