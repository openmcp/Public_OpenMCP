/*
Copyright 2018 The Multicluster-Controller Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"os"
	"log"
	"fmt"
	"k8s.io/klog"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"
	"openmcp/openmcp/openmcp-scheduler/pkg"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/util/controller/logLevel"
	"openmcp/openmcp/util/controller/reshape"
	"openmcp/openmcp/openmcp-scheduler/pkg/protobuf"
	"openmcp/openmcp/util/controller/logLevel"
	"google.golang.org/grpc"
)

func main() {
	// setting logging
	logLevel.KetiLogInit()

	// setting grpc connection
	SERVER_IP := os.Getenv("GRPC_SERVER")
	SERVER_PORT := os.Getenv("GRPC_PORT")
	grpcClient := newGrpcClient(SERVER_IP, SERVER_PORT)


	for{
		klog.V(0).Info("***** [START] OpenMCP Scheduler *****")
		cm := clusterManager.NewClusterManager()

		host_ctx := "openmcp"
		namespace := "openmcp"

		host_cfg := cm.Host_config
		live := cluster.New(host_ctx, host_cfg, cluster.Options{})

		ghosts := []*cluster.Cluster{}

		for _, ghost_cluster := range cm.Cluster_list.Items {
			ghost_ctx := ghost_cluster.Name
			ghost_cfg := cm.Cluster_configs[ghost_ctx]

			ghost := cluster.New(ghost_ctx, ghost_cfg, cluster.Options{})
			ghosts = append(ghosts, ghost)
		}

		sched_cont, err := openmcpscheduler.NewController(live, ghosts, namespace, grpcClient) 
		if err != nil {
			klog.V(0).Info("err New Controller - Scheduler", err)
		}
		reshape_cont, err := reshape.NewController(live, ghosts, namespace)
		if err != nil {
			klog.V(0).Info("err New Controller - Reshape", err)
		}
		loglevel_cont, err := logLevel.NewController(live, ghosts, namespace)
		if err != nil {
			klog.V(0).Info("err New Controller - logLevel", err)
		}

		m := manager.New()
		m.AddController(sched_cont)
		m.AddController(reshape_cont)
		m.AddController(loglevel_cont)

		stop := reshape.SetupSignalHandler()

		if err := m.Start(stop); err != nil {
			log.Fatal(err)
		}
	}
}

func newGrpcClient(ip, port string) protobuf.RequestAnalysisClient {
	host := ip + ":" + port
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		klog.V(0).Info("did not connect", err)
	}
	defer conn.Close()

	c := protobuf.NewRequestAnalysisClient(conn)
	return c
}
