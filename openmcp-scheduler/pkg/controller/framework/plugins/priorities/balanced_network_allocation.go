package priorities

import (
	"os"
	"fmt"
	"context"
	"strconv"
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/controller/resourceinfo"
	"openmcp/openmcp/openmcp-scheduler/pkg/controller/protobuf"
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
		
		node_networks := &protobuf.SchedInfo{NodeName: node.NodeName}
		result, _ := grpcClient.SendSchedAnalysis(context.TODO(), node_networks)
		rx, _ := strconv.ParseFloat(fmt.Sprintf("%s", result.RX), 64)
		tx, _ := strconv.ParseFloat(fmt.Sprintf("%s", result.TX), 64)

		rx = rx / 1000 / 1000 // change to mega
		tx = tx / 1000 / 1000 // change to mega
		nodeScore := int64(1 / (rx + tx)) * maxScore
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