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
	"log"
	"k8s.io/klog"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/sample-controller/pkg/signals"

	"openmcpscheduler/pkg/controller"
	"openmcpscheduler/pkg/metric"
)

func main() {
	grpc_client := metric.NewGRPCClient()
	klog.Infof("*********** CREATE GRPC CLIENT ***********")

	cluster_manager := openmcpscheduler.NewClusterManager()

	host_ctx := "openmcp"
	namespace := "openmcp"

	host_cfg := cluster_manager.Host_config
	live := cluster.New(host_ctx, host_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{Namespace: namespace}})

	ghosts := []*cluster.Cluster{}
	for _, ghost_cluster := range cluster_manager.Cluster_list.Items {
		ghost_ctx := ghost_cluster.Name
		ghost_cfg := cluster_manager.Cluster_configs[ghost_ctx]

		ghost := cluster.New(ghost_ctx, ghost_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{Namespace: namespace}})
		ghosts = append(ghosts, ghost)
	}

	// for checking
	for _, ghost := range ghosts {
		klog.Infof("ghost.Name: %v", ghost.Name)
	}

	co, _ := openmcpscheduler.NewController(live, ghosts, namespace)
	m := manager.New()
	m.AddController(co)

	if err := m.Start(signals.SetupSignalHandler()); err != nil {
		log.Fatal(err)
	}
}