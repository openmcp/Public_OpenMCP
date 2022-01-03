package main

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/golang/protobuf/ptypes/timestamp"

	//"github.com/golang/protobuf/ptypes/timestamp"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-metric-collector/member/src/customMetrics"
	"openmcp/openmcp/openmcp-metric-collector/member/src/kubeletClient"
	"openmcp/openmcp/openmcp-metric-collector/member/src/protobuf"
	"openmcp/openmcp/openmcp-metric-collector/member/src/scrap"
	"openmcp/openmcp/openmcp-metric-collector/member/src/storage"
	"openmcp/openmcp/util/clusterManager"

	"github.com/jinzhu/copier"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	//"github.com/jinzhu/copier"

	//"github.com/golang/protobuf/ptypes"
	//"github.com/golang/protobuf/ptypes/timestamp"

	//"github.com/golang/protobuf/ptypes"
	"os"
	//timestamp "github.com/golang/protobuf/ptypes/timestamp"
	//"github.com/jinzhu/copier"

	//"context"
	"fmt"
	"time"
)

func convert(data *storage.Collection, latencyTime string) *protobuf.Collection {
	//klog.V(0).Info("Convert GRPC Data Structure")
	grpc_data := &protobuf.Collection{}

	copier.Copy(grpc_data, data)
	for i, _ := range grpc_data.Metricsbatchs {

		s := int64(data.Metricsbatchs[i].Node.Timestamp.Second())     // from 'int'
		n := int32(data.Metricsbatchs[i].Node.Timestamp.Nanosecond()) // from 'int'

		ts := &timestamp.Timestamp{Seconds: s, Nanos: n}

		mp := &protobuf.MetricsPoint{
			Timestamp:             ts,
			CPUUsageNanoCores:     data.Metricsbatchs[i].Node.CPUUsageNanoCores.String(),
			MemoryUsageBytes:      data.Metricsbatchs[i].Node.MemoryUsageBytes.String(),
			MemoryAvailableBytes:  data.Metricsbatchs[i].Node.MemoryAvailableBytes.String(),
			MemoryWorkingSetBytes: data.Metricsbatchs[i].Node.MemoryWorkingSetBytes.String(),
			NetworkRxBytes:        data.Metricsbatchs[i].Node.NetworkRxBytes.String(),
			NetworkTxBytes:        data.Metricsbatchs[i].Node.NetworkTxBytes.String(),
			FsAvailableBytes:      data.Metricsbatchs[i].Node.FsAvailableBytes.String(),
			FsCapacityBytes:       data.Metricsbatchs[i].Node.FsCapacityBytes.String(),
			FsUsedBytes:           data.Metricsbatchs[i].Node.FsUsedBytes.String(),
			NetworkLatency:        latencyTime,
		}
		grpc_data.Metricsbatchs[i].Node.MP = mp

		//fmt.Println(grpc_data.Metricsbatchs[0].IP)
		//fmt.Println(grpc_data.Metricsbatchs[0].Node.Name)
		//fmt.Println(grpc_data.Metricsbatchs[0].Node.MP.Timestamp.String())
		//fmt.Println(grpc_data.Metricsbatchs[0].Node.MP.Timestamp.Seconds)
		//fmt.Println(grpc_data.Metricsbatchs[0].Node.MP.CpuUsage)
		//fmt.Println(grpc_data.Metricsbatchs[0].Node.MP.MemoryUsage)

		podMetricsPoints := []*protobuf.PodMetricsPoint{}

		for j, _ := range data.Metricsbatchs[i].Pods {
			s := int64(data.Metricsbatchs[i].Pods[j].Timestamp.Second())     // from 'int'
			n := int32(data.Metricsbatchs[i].Pods[j].Timestamp.Nanosecond()) // from 'int'

			ts := &timestamppb.Timestamp{Seconds: s, Nanos: n}

			mp2 := &protobuf.MetricsPoint{
				Timestamp:             ts,
				CPUUsageNanoCores:     data.Metricsbatchs[i].Pods[j].CPUUsageNanoCores.String(),
				MemoryUsageBytes:      data.Metricsbatchs[i].Pods[j].MemoryUsageBytes.String(),
				MemoryAvailableBytes:  data.Metricsbatchs[i].Pods[j].MemoryAvailableBytes.String(),
				MemoryWorkingSetBytes: data.Metricsbatchs[i].Pods[j].MemoryWorkingSetBytes.String(),
				NetworkRxBytes:        data.Metricsbatchs[i].Pods[j].NetworkRxBytes.String(),
				NetworkTxBytes:        data.Metricsbatchs[i].Pods[j].NetworkTxBytes.String(),
				FsAvailableBytes:      data.Metricsbatchs[i].Pods[j].FsAvailableBytes.String(),
				FsCapacityBytes:       data.Metricsbatchs[i].Pods[j].FsCapacityBytes.String(),
				FsUsedBytes:           data.Metricsbatchs[i].Pods[j].FsUsedBytes.String(),
				NetworkLatency:        latencyTime,
			}
			pmp := &protobuf.PodMetricsPoint{
				Name:       data.Metricsbatchs[i].Pods[j].Name,
				Namespace:  data.Metricsbatchs[i].Pods[j].Namespace,
				MP:         mp2,
				Containers: nil,
			}
			podMetricsPoints = append(podMetricsPoints, pmp)

		}
		grpc_data.Metricsbatchs[i].Pods = podMetricsPoints

		//fmt.Println(grpc_data.Metricsbatchs[0].IP)
		//fmt.Println(grpc_data.Metricsbatchs[0].Pods[0].Name)
		//fmt.Println(grpc_data.Metricsbatchs[0].Pods[0].MP.Timestamp.String())
		//fmt.Println(grpc_data.Metricsbatchs[0].Pods[0].MP.Timestamp.Seconds)
		//fmt.Println(grpc_data.Metricsbatchs[0].Pods[0].MP.CpuUsage)
		//fmt.Println(grpc_data.Metricsbatchs[0].Pods[0].MP.MemoryUsage)

	}

	return grpc_data

}
func main() {

	MemberMetricCollector()
}

func MemberMetricCollector() {
	SERVER_IP := os.Getenv("GRPC_SERVER")
	SERVER_PORT := os.Getenv("GRPC_PORT")

	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100

	fmt.Println("ClusterMetricCollector Start")
	grpcClient := protobuf.NewGrpcClient(SERVER_IP, SERVER_PORT)

	var period_int64 int64 = 5
	var latencyTime float64 = 0

	host_config, err := rest.InClusterConfig()
	if err != nil {
		omcplog.V(0).Info(err)
	}
	host_kubeClient := kubernetes.NewForConfigOrDie(host_config)

	token := host_config.BearerToken
	host := host_config.Host
	kubeclient := host_kubeClient

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}

	for {

		Node_list, err := clusterManager.GetNodeList(host_kubeClient)
		if err != nil {
			omcplog.V(0).Info(err)
		}

		nodes := Node_list.Items
		fmt.Println("Get Metric Data From Kubelet")
		kubeletClient, _ := kubeletClient.NewKubeletClient()

		data, errs := scrap.Scrap(host_config, kubeletClient, nodes)

		if errs != nil {
			fmt.Println(errs)
			time.Sleep(time.Duration(period_int64) * time.Second)
			continue
		}
		fmt.Println("Convert Metric Data For gRPC")

		latencyTime_string := fmt.Sprintf("%f", latencyTime)

		grpc_data := convert(data, latencyTime_string)

		fmt.Println("[gRPC Start] Send Metric Data")

		rTime_start := time.Now()
		r, err := grpcClient.SendMetrics(context.TODO(), grpc_data)
		if err == nil {
			rTime_end := time.Since(rTime_start)

			latencyTime = rTime_end.Seconds() - r.ProcessingTime
		} else {
			//fmt.Println("check")
			fmt.Println("could not connect : ", err)
			time.Sleep(time.Duration(period_int64) * time.Second)
			//fmt.Println("check2")
			continue
		}
		fmt.Println("[gRPC End] Send Metric Data")

		//period_int64 := r.Tick
		// _ = data

		fmt.Println("[http Start] Post Metric Data to Custom Metric Server")

		//fmt.Println("host: ", host)
		//fmt.Println("token: ", token)
		//fmt.Println("client: ", client)

		customMetrics.AddToPodCustomMetricServer(data, token, host, client)
		customMetrics.AddToDeployCustomMetricServer(data, token, host, kubeclient, client)
		fmt.Println("[http End] Post Metric Data to Custom Metric Server")

		period_int64 = r.Tick

		if period_int64 > 0 && err == nil {

			//fmt.Println("period : ",time.Duration(period_int64))
			fmt.Println("Wait ", time.Duration(period_int64)*time.Second, "...")
			time.Sleep(time.Duration(period_int64) * time.Second)
		} else {
			fmt.Println("--- Fail to get period")
			time.Sleep(5 * time.Second)
		}
	}
}
