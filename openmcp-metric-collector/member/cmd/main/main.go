
package main

import (
	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"
	"context"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jinzhu/copier"
	"log"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/customMetrics"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/kubeletClient"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/protobuf"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/scrap"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/storage"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/util/controller/logLevel"
	"openmcp/openmcp/util/controller/reshape"

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
	logLevel.KetiLogInit()

	go MemberMetricCollector()

	for {
		cm := clusterManager.NewClusterManager()

		host_ctx := "openmcp"
		namespace := "openmcp"

		host_cfg := cm.Host_config
		//live := cluster.New(host_ctx, host_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{Namespace: namespace}})
		live := cluster.New(host_ctx, host_cfg, cluster.Options{})

		ghosts := []*cluster.Cluster{}

		for _, ghost_cluster := range cm.Cluster_list.Items {
			ghost_ctx := ghost_cluster.Name
			ghost_cfg := cm.Cluster_configs[ghost_ctx]

			//ghost := cluster.New(ghost_ctx, ghost_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{Namespace: namespace}})
			ghost := cluster.New(ghost_ctx, ghost_cfg, cluster.Options{})
			ghosts = append(ghosts, ghost)
		}

		reshape_cont, _ := reshape.NewController(live, ghosts, namespace)
		loglevel_cont, _ := logLevel.NewController(live, ghosts, namespace)

		m := manager.New()
		m.AddController(reshape_cont)
		m.AddController(loglevel_cont)

		stop := reshape.SetupSignalHandler()

		if err := m.Start(stop); err != nil {
			log.Fatal(err)
		}
	}

}
func MemberMetricCollector(){
	SERVER_IP := os.Getenv("GRPC_SERVER")
	SERVER_PORT := os.Getenv("GRPC_PORT")
	omcplog.V(2).Info("ClusterMetricCollector Start")
	grpcClient := protobuf.NewGrpcClient(SERVER_IP, SERVER_PORT)

	for {
		cm := clusterManager.NewClusterManager()
		nodes := cm.Node_list.Items
		omcplog.V(2).Info("Get Metric Data From Kubelet")
		kubeletClient, _ := kubeletClient.NewKubeletClient()
		data, errs := scrap.Scrap(cm.Host_config, kubeletClient, nodes)
		if errs != nil {
			omcplog.V(0).Info(errs)
		}
		omcplog.V(2).Info("Convert Metric Data For gRPC")
		grpc_data := convert(data)

		//fmt.Println("GRPC Data Send")
		omcplog.V(2).Info("[gRPC Start] Send Metric Data")
		r, err := grpcClient.SendMetrics(context.TODO(), grpc_data)
		if err != nil {
			fmt.Printf("could not connect : %v", err)
		}
		omcplog.V(2).Info("[gRPC End] Send Metric Data")
		//period_int64 := r.Tick
		_ = data

		omcplog.V(2).Info("[http Start] Post Metric Data to Custom Metric Server")
		token := cm.Host_config.BearerToken
		host := cm.Host_config.Host
		client := cm.Host_kubeClient
		//fmt.Println("host: ", host)
		//fmt.Println("token: ", token)
		//fmt.Println("client: ", client)

		customMetrics.AddToPodCustomMetricServer(data, token, host)
		customMetrics.AddToDeployCustomMetricServer(data, token, host, client)
		omcplog.V(2).Info("[http End] Post Metric Data to Custom Metric Server")

		period_int64 := r.Tick

		if period_int64 > 0 && err == nil {

			//fmt.Println("period : ",time.Duration(period_int64))
			omcplog.V(2).Info("Wait ", time.Duration(period_int64)* time.Second, "...")
			time.Sleep(time.Duration(period_int64) * time.Second)
		}else {
			omcplog.V(2).Info("--- Fail to get period")
			time.Sleep(5 * time.Second)
		}
	}
}