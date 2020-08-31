package priorities

import (
	"context"
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
	"openmcp/openmcp/openmcp-scheduler/pkg/protobuf"
	"openmcp/openmcp/omcplog"
)

type BalancedNetworkAllocation struct{
	GRPC_Client		protobuf.RequestAnalysisClient
}

func (pl *BalancedNetworkAllocation) Name() string {
	return "BalancedNetworkAllocation"
}

func (pl *BalancedNetworkAllocation) Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) int64 {
	var clutserScore int64

	for _, node := range clusterInfo.Nodes {

		node_info := &protobuf.NodeInfo{ClusterName: clusterInfo.ClusterName, NodeName: node.NodeName}
		client := pl.GRPC_Client
		result, err := client.SendNetworkAnalysis(context.TODO(), node_info)

		if err != nil || result == nil {
			omcplog.V(0).Infof("cannot get %v's data from openmcp-analytic-engine", node.NodeName)
			continue
		}

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

	return clutserScore
}