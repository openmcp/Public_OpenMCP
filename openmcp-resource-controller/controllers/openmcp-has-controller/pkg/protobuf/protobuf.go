package protobuf

import (
	"google.golang.org/grpc"
	"openmcp/openmcp/omcplog"
)

func NewGrpcClient(ip, port string) RequestAnalysisClient {
	host := ip + ":" + port
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		omcplog.V(0).Info("did not connect: %v", err)
	}
	c := NewRequestAnalysisClient(conn)
	return c
}
