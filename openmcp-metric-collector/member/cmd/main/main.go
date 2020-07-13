
package main

import (
	"context"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jinzhu/copier"
	"k8s.io/klog"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/customMetrics"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/kubeletClient"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/protobuf"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/scrap"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/storage"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/util/controller/logLevel"

	//"github.com/jinzhu/copier"

	//"github.com/golang/protobuf/ptypes"
	//"github.com/golang/protobuf/ptypes/timestamp"

	//"github.com/golang/protobuf/ptypes"
	"os"
	//timestamp "github.com/golang/protobuf/ptypes/timestamp"
	//"github.com/jinzhu/copier"

	"time"
)

func convert(data *storage.Collection) *protobuf.Collection{
	klog.V(0).Infof( "Convert GRPC Data Structure")

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

		//klog.V(0).Infof( grpc_data.Matricsbatchs[0].IP)
		//klog.V(0).Infof( grpc_data.Matricsbatchs[0].Node.Name)
		//klog.V(0).Infof( grpc_data.Matricsbatchs[0].Node.MP.Timestamp.String())
		//klog.V(0).Infof( grpc_data.Matricsbatchs[0].Node.MP.Timestamp.Seconds)
		//klog.V(0).Infof( grpc_data.Matricsbatchs[0].Node.MP.CpuUsage)
		//klog.V(0).Infof( grpc_data.Matricsbatchs[0].Node.MP.MemoryUsage)

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

		//klog.V(0).Infof( grpc_data.Matricsbatchs[0].IP)
		//klog.V(0).Infof( grpc_data.Matricsbatchs[0].Pods[0].Name)
		//klog.V(0).Infof( grpc_data.Matricsbatchs[0].Pods[0].MP.Timestamp.String())
		//klog.V(0).Infof( grpc_data.Matricsbatchs[0].Pods[0].MP.Timestamp.Seconds)
		//klog.V(0).Infof( grpc_data.Matricsbatchs[0].Pods[0].MP.CpuUsage)
		//klog.V(0).Infof( grpc_data.Matricsbatchs[0].Pods[0].MP.MemoryUsage)

	}


	return grpc_data

}
var default_period int64 = 5
func main(){
	logLevel.KetiLogInit()
	
	SERVER_IP := os.Getenv("GRPC_SERVER")
	SERVER_PORT := os.Getenv("GRPC_PORT")

	grpcClient := protobuf.NewGrpcClient(SERVER_IP, SERVER_PORT)

	period_int64 := default_period // default

	for {
		cm := clusterManager.NewClusterManager()
		nodes := cm.Node_list.Items

		kubeletClient, _ := kubeletClient.NewKubeletClient()
		data, errs := scrap.Scrap(cm.Host_config, kubeletClient, nodes)
		if errs != nil {
			klog.V(0).Info( errs)
		}

		grpc_data := convert(data)

		klog.V(0).Info( "GRPC Data Send")
		r, err := grpcClient.SendMetrics(context.TODO(), grpc_data)
		if err != nil {
			klog.V(0).Info( "could not connect : ", err)
		} else {
			//period_int64 := r.Tick
			_ = data


			token := cm.Host_config.BearerToken
			host := cm.Host_config.Host
			client := cm.Host_kubeClient
			//klog.V(0).Infof( "host: ", host)
			//klog.V(0).Infof( "token: ", token)
			//klog.V(0).Infof( "client: ", client)

			customMetrics.AddToPodCustomMetricServer(data, token, host)
			customMetrics.AddToDeployCustomMetricServer(data, token, host, client)

			period_int64 = r.Tick
		}


		if period_int64 > 0 && err == nil {

			klog.V(0).Info( "period : ",time.Duration(period_int64) * time.Second)
			time.Sleep(time.Duration(period_int64) * time.Second)
		}else {
			klog.V(0).Info( "--- Fail to get period")
			time.Sleep(time.Duration(default_period) * time.Second)
		}
	}
}