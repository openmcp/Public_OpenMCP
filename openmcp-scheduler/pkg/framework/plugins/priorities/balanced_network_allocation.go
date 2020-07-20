package priorities

import (
	"os"
	"fmt"
	"context"
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
	"openmcp/openmcp/openmcp-scheduler/pkg/protobuf"
	"google.golang.org/grpc"
)

type BalancedNetworkAllocation struct{}

const (
	database = "Metrics"
)

func (pl *BalancedNetworkAllocation) Name() string {
	return "BalancedNetworkAllocation"
}

func (pl *BalancedNetworkAllocation) Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) int64 {
	var clutserScore int64

	// Get InfluxDB from openmcp's InfluxDB Pod
	SERVER_IP := os.Getenv("GRPC_SERVER")
	SERVER_PORT := os.Getenv("GRPC_PORT")
	grpcClient := NewGrpcClient(SERVER_IP, SERVER_PORT)

	for _, node := range clusterInfo.Nodes {

		node_info := &protobuf.NodeInfo{ClusterName: clusterInfo.ClusterName, NodeName: node.NodeName}
		result, _ := grpcClient.SendNetworkAnalysis(context.TODO(), node_info)

		var nodeScore int64
		rx := result.RX
		tx := result.TX

		if rx == 0 && tx == 0 {
			nodeScore = 0
		}else {
			nodeScore = int64((1 / float64(rx + tx)) * float64(maxScore))
		}
		clutserScore += nodeScore
	}

	return clutserScore
}

func NewGrpcClient(ip, port string) protobuf.RequestAnalysisClient {
	host := ip + ":" + port
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("did not connect: %v", err)
	}
	c := protobuf.NewRequestAnalysisClient(conn)
	return c
}