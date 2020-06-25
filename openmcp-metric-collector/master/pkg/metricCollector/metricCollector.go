package metricCollector

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/openmcp-metric-collector/master/pkg/influx"
	"openmcp/openmcp/openmcp-metric-collector/master/pkg/protobuf"
	"strconv"

	//"strconv"
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
	for _, cluster := range mc.ClusterManager.Cluster_list.Items {
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
	var period_int64 int64

	openmcpPolicyInstance, target_cluster_policy_err := mc.ClusterManager.Crd_client.OpenMCPPolicyEngine("openmcp").Get("metric-collector-period", metav1.GetOptions{})

	if target_cluster_policy_err != nil {
		fmt.Println(target_cluster_policy_err)
	} else {
		a := openmcpPolicyInstance.Spec.Template.Spec.Policies
		period := a[0].Value[0]
		fmt.Println("period : ",period)
		period_int64,_ = strconv.ParseInt(period, 10, 64)
	}

	//timetick := 3600
	clustername := "c1"
	return &protobuf.ReturnValue{
		Tick:        period_int64,
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
