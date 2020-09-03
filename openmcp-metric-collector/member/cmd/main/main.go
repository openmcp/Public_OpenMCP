
package main

import (
	"context"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jinzhu/copier"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/customMetrics"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/kubeletClient"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/protobuf"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/scrap"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/storage"
	"openmcp/openmcp/util/clusterManager"
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

func convert(data *storage.Collection) *protobuf.Collection{
	//klog.V(0).Info("Convert GRPC Data Structure")

	grpc_data := &protobuf.Collection{}

	copier.Copy(grpc_data, data)

	for i, _ := range grpc_data.Matricsbatchs {

		s := int64(data.Matricsbatchs[i].Node.Timestamp.Second()) // from 'int'
		n := int32(data.Matricsbatchs[i].Node.Timestamp.Nanosecond()) // from 'int'

		ts := &timestamp.Timestamp{Seconds:s, Nanos:n}

		mp := &protobuf.MetricsPoint{
			Timestamp: ts,
			CPUUsageNanoCores: data.Matricsbatchs[i].Node.CPUUsageNanoCores.String(),
			MemoryUsageBytes: data.Matricsbatchs[i].Node.MemoryUsageBytes.String(),
			MemoryAvailableBytes: data.Matricsbatchs[i].Node.MemoryAvailableBytes.String(),
			MemoryWorkingSetBytes: data.Matricsbatchs[i].Node.MemoryWorkingSetBytes.String(),
			NetworkRxBytes: data.Matricsbatchs[i].Node.NetworkRxBytes.String(),
			NetworkTxBytes: data.Matricsbatchs[i].Node.NetworkTxBytes.String(),
			FsAvailableBytes: data.Matricsbatchs[i].Node.FsAvailableBytes.String(),
			FsCapacityBytes: data.Matricsbatchs[i].Node.FsCapacityBytes.String(),
			FsUsedBytes: data.Matricsbatchs[i].Node.FsUsedBytes.String(),
		}
		grpc_data.Matricsbatchs[i].Node.MP = mp

		//fmt.Println(grpc_data.Matricsbatchs[0].IP)
		//fmt.Println(grpc_data.Matricsbatchs[0].Node.Name)
		//fmt.Println(grpc_data.Matricsbatchs[0].Node.MP.Timestamp.String())
		//fmt.Println(grpc_data.Matricsbatchs[0].Node.MP.Timestamp.Seconds)
		//fmt.Println(grpc_data.Matricsbatchs[0].Node.MP.CpuUsage)
		//fmt.Println(grpc_data.Matricsbatchs[0].Node.MP.MemoryUsage)

		podMetricsPoints := []*protobuf.PodMetricsPoint{}

		for j, _ := range data.Matricsbatchs[i].Pods {
			s := int64(data.Matricsbatchs[i].Pods[j].Timestamp.Second()) // from 'int'
			n := int32(data.Matricsbatchs[i].Pods[j].Timestamp.Nanosecond()) // from 'int'

			ts := &timestamp.Timestamp{Seconds:s, Nanos:n}
		
			mp2 := &protobuf.MetricsPoint{
				Timestamp: ts,
				CPUUsageNanoCores: data.Matricsbatchs[i].Pods[j].CPUUsageNanoCores.String(),
				MemoryUsageBytes: data.Matricsbatchs[i].Pods[j].MemoryUsageBytes.String(),
				MemoryAvailableBytes: data.Matricsbatchs[i].Pods[j].MemoryAvailableBytes.String(),
				MemoryWorkingSetBytes: data.Matricsbatchs[i].Pods[j].MemoryWorkingSetBytes.String(),
				NetworkRxBytes: data.Matricsbatchs[i].Pods[j].NetworkRxBytes.String(),
				NetworkTxBytes: data.Matricsbatchs[i].Pods[j].NetworkTxBytes.String(),
				FsAvailableBytes: data.Matricsbatchs[i].Pods[j].FsAvailableBytes.String(),
				FsCapacityBytes: data.Matricsbatchs[i].Pods[j].FsCapacityBytes.String(),
				FsUsedBytes: data.Matricsbatchs[i].Pods[j].FsUsedBytes.String(),
			}
			pmp := &protobuf.PodMetricsPoint{
				Name:       data.Matricsbatchs[i].Pods[j].Name,
				Namespace:  data.Matricsbatchs[i].Pods[j].Namespace,
				MP:         mp2,
				Containers: nil,
			}
			podMetricsPoints = append(podMetricsPoints, pmp)

		}
		grpc_data.Matricsbatchs[i].Pods = podMetricsPoints

		//fmt.Println(grpc_data.Matricsbatchs[0].IP)
		//fmt.Println(grpc_data.Matricsbatchs[0].Pods[0].Name)
		//fmt.Println(grpc_data.Matricsbatchs[0].Pods[0].MP.Timestamp.String())
		//fmt.Println(grpc_data.Matricsbatchs[0].Pods[0].MP.Timestamp.Seconds)
		//fmt.Println(grpc_data.Matricsbatchs[0].Pods[0].MP.CpuUsage)
		//fmt.Println(grpc_data.Matricsbatchs[0].Pods[0].MP.MemoryUsage)

	}


	return grpc_data

}
func main() {
	MemberMetricCollector()
}

func MemberMetricCollector(){
	SERVER_IP := os.Getenv("GRPC_SERVER")
	SERVER_PORT := os.Getenv("GRPC_PORT")
	fmt.Println("ClusterMetricCollector Start")
	grpcClient := protobuf.NewGrpcClient(SERVER_IP, SERVER_PORT)

	for {
		cm := clusterManager.NewClusterManager()
		nodes := cm.Node_list.Items
		fmt.Println("Get Metric Data From Kubelet")
		kubeletClient, _ := kubeletClient.NewKubeletClient()
		data, errs := scrap.Scrap(cm.Host_config, kubeletClient, nodes)
		if errs != nil {
			fmt.Println(errs)
		}
		fmt.Println("Convert Metric Data For gRPC")
		grpc_data := convert(data)

		//fmt.Println("GRPC Data Send")
		fmt.Println("[gRPC Start] Send Metric Data")
		r, err := grpcClient.SendMetrics(context.TODO(), grpc_data)
		if err != nil {
			fmt.Printf("could not connect : %v", err)
		}
		fmt.Println("[gRPC End] Send Metric Data")
		//period_int64 := r.Tick
		_ = data

		fmt.Println("[http Start] Post Metric Data to Custom Metric Server")
		token := cm.Host_config.BearerToken
		host := cm.Host_config.Host
		client := cm.Host_kubeClient
		//fmt.Println("host: ", host)
		//fmt.Println("token: ", token)
		//fmt.Println("client: ", client)

		customMetrics.AddToPodCustomMetricServer(data, token, host)
		customMetrics.AddToDeployCustomMetricServer(data, token, host, client)
		fmt.Println("[http End] Post Metric Data to Custom Metric Server")

		period_int64 := r.Tick

		if period_int64 > 0 && err == nil {

			//fmt.Println("period : ",time.Duration(period_int64))
			fmt.Println("Wait ", time.Duration(period_int64)* time.Second, "...")
			time.Sleep(time.Duration(period_int64) * time.Second)
		}else {
			fmt.Println("--- Fail to get period")
			time.Sleep(5 * time.Second)
		}
	}
}