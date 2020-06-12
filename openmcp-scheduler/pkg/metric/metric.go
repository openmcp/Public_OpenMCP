package metric

import (
	"os"
	"context"
	"k8s.io/klog"
	"google.golang.org/grpc"
)

func NewGRPCClient() *NodeInfosClient {
	server_ip := os.Getenv("GRPC_SERVER")
	server_port := os.Getenv("GRPC_PORT")
	serverAddr := server_ip + ":" + server_port

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		klog.Infof("cannot connect to grpcServer: %v", err)
	}
	defer conn.Close()

	client := NewNodeInfosClient(conn)
	return &client
}

func GetCurrentNetwork(client *NodeInfosClient, nodeName string) (float32, float32) {
	// Get Network data from GRPC Server (=analytic engine)
	result, err := client.GetNodeInfos(context.TODO(), &RequestedNode{Name:nodeName})
	if err != nil {
		klog.Infof("cannot get Current Network Data from analytic engine: %v", err)
	}

	// Extract Data from result and return datas
	return result.CurrentRXBytes, result.CurrentTXBytes
}