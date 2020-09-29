package metricCollector

import (
	"context"
	"google.golang.org/grpc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"net"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-metric-collector/master/pkg/influx"
	"openmcp/openmcp/openmcp-metric-collector/master/pkg/protobuf"
	"openmcp/openmcp/util/clusterManager"
	"strconv"
)

type MetricCollector struct {
	ClusterManager clusterManager.ClusterManager
	Influx         influx.Influx
}

func NewMetricCollector(INFLUX_IP, INFLUX_PORT, INFLUX_USERNAME, INFLUX_PASSWORD string) *MetricCollector {
	omcplog.V(4).Info("NewMetricCollector Called")
	mc := &MetricCollector{}
	mc.ClusterManager = *clusterManager.NewClusterManager()
	mc.Influx = *influx.NewInflux(INFLUX_IP, INFLUX_PORT, INFLUX_USERNAME, INFLUX_PASSWORD)

	return mc
}
func (mc *MetricCollector) FindClusterName(data *protobuf.Collection) string {
	omcplog.V(4).Info("FindClusterName Called")
	IpList := []string{}
	for _, Matricsbatch := range data.Matricsbatchs {
		omcplog.V(2).Info("[Recieved Data] NodeName: ", Matricsbatch.Node.Name, ", IP: "+ Matricsbatch.IP)
		IpList = append(IpList, Matricsbatch.IP)
	}
	clusterName := data.ClusterName
	klog.V(2).Info("=> Recieved Metric Data From '", clusterName,"'")
	return clusterName
}

func (mc *MetricCollector) SendMetrics(ctx context.Context, data *protobuf.Collection) (*protobuf.ReturnValue, error) {
	omcplog.V(4).Info("SendMetrics Called")
	clusterName := mc.FindClusterName(data)
	mc.Influx.SaveMetrics(clusterName, data)
	var period_int64 int64

	openmcpPolicyInstance, target_cluster_policy_err := mc.ClusterManager.Crd_client.OpenMCPPolicy("openmcp").Get("metric-collector-period", metav1.GetOptions{})

	if target_cluster_policy_err != nil {
		klog.V(0).Info( target_cluster_policy_err)

	} else {
		a := openmcpPolicyInstance.Spec.Template.Spec.Policies
		period := a[0].Value[0]
		klog.V(3).Info("getPeriodPolicy: ",period+" sec")
		period_int64,_ = strconv.ParseInt(period, 10, 64)
	}

	clustername := "c1"

	klog.V(2).Info("gRPC Return Period")
	return &protobuf.ReturnValue{
		Tick:        period_int64,
		ClusterName: clustername,
	}, nil
}

func (mc *MetricCollector) StartGRPC(GRPC_PORT string) {
	omcplog.V(4).Info("StartGRPC Called")
	omcplog.V(2).Info("Grpc Server Start at Port %s\n", GRPC_PORT)
	
	l, err := net.Listen("tcp", ":"+GRPC_PORT)
	if err != nil {
		omcplog.V(0).Info("failed to listen: ", err)

	}
	grpcServer := grpc.NewServer()

	protobuf.RegisterSendMetricsServer(grpcServer, mc)
	if err := grpcServer.Serve(l); err != nil {
		omcplog.V(0).Info("fail to serve: ", err)

	}

}
