package openmcpscheduler

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"k8s.io/klog"
	pb "openmcp-scheduler/pkg/controller/openmcpscheduler/pb"
)

const (
	serverAddr = "10.109.183.52:3036"
)

func Clients(policy string) map[string]float64 {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())

	if err != nil {
		klog.Infof("[SUJUNE] did not connect: %v", err)
	}
	defer conn.Close()

	klog.Infof("[SUJUNE] connected Server!!")
	c := pb.NewPriorityListClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetPriorityList(ctx, &pb.PriorityRequest{Policy: policy})
	if err != nil {
		klog.Infof("[SUJUNE] could no greet: %v", err)
	}
	klog.Infof("[SUJUNE] list = %s", r.List)

	return r.List
}
