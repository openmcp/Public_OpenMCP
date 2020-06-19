package metricCollector

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"openmcp-metric-collector/pkg/clusterManager"
	"openmcp-metric-collector/pkg/influx"
	"openmcp-metric-collector/pkg/protobuf"
	"strings"
)

type MetricCollector struct {
	ClusterManager clusterManager.ClusterManager
	Influx         influx.Influx
}

func NewMetricCollector(INFLUX_IP, INFLUX_PORT, INFLUX_USERNAME, INFLUX_PASSWORD string) *MetricCollector {
	mc := &MetricCollector{}
	mc.ClusterManager = *clusterManager.NewClusterManager()
	mc.Influx = *influx.NewInflux(INFLUX_IP, INFLUX_PORT, INFLUX_USERNAME, INFLUX_PASSWORD)

	return mc
}
func (mc *MetricCollector) FindClusterName(data *protobuf.Collection) string {
	IpList := []string{}
	for _, Matricsbatch := range data.Matricsbatchs {
		log.Println("Recieved : ", Matricsbatch.Node.Name, "::::"+Matricsbatch.IP)
		IpList = append(IpList, Matricsbatch.IP)
	}

	find := false
	clusterName := ""
	for _, cluster := range mc.ClusterManager.ClusterList.Items {
		for _, Ip := range IpList {
			if strings.Contains(cluster.Spec.APIEndpoint, Ip) {
				clusterName = cluster.Name
				find = true
				break
			}
		}
		if find {
			break
		}
	}
	log.Printf("Recieved Metric Data from '%s'\n", clusterName)
	return clusterName
}

func (mc *MetricCollector) SendMetrics(ctx context.Context, data *protobuf.Collection) (*protobuf.ReturnValue, error) {
	//startTime := time.Now()
	clusterName := mc.FindClusterName(data)
	mc.Influx.SaveMetrics(clusterName, data)

	timetick := 3600
	clustername := "c1"
	return &protobuf.ReturnValue{
		Tick:        int64(timetick),
		ClusterName: clustername,
	}, nil
	//return nil, nil
}

//func (mc *MetricCollector) SendMetrics(ctx context.Context, data *protobuf.Collection) (*protobuf.ReturnValue, error) {
//	fmt.Println(data.Matricsbatchs[0].Node.Name)
//	fmt.Printf("%s",data.Matricsbatchs[0].Node.MP.CpuUsage)
//
//
//	timetick:= 3600
//	clustername := "c1"
//	return &protobuf.ReturnValue{
//		Tick:                 int64(timetick),
//		ClusterName:          clustername,
//	}, nil
//	//return nil, nil
//}
func (mc *MetricCollector) StartGRPC(GRPC_PORT string) {
	log.Printf("Grpc Server Start at Port %s\n", GRPC_PORT)

	//manager = NewClusterManager()
	l, err := net.Listen("tcp", ":"+GRPC_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	protobuf.RegisterSendMetricsServer(grpcServer, mc)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("fail to serve: %v", err)
	}

}
