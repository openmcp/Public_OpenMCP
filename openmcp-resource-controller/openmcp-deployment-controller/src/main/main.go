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
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-resource-controller/openmcp-deployment-controller/src/controller"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/util/controller/logLevel"
	"openmcp/openmcp/util/controller/reshape"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"
	corev1 "k8s.io/api/core/v1"
)

var cm = clusterManager.NewClusterManager()

func main() {
	logLevel.KetiLogInit()

	omcplog.V(2).Info("Start OpenMCP Deployment Controller")

	//cm := clusterManager.NewClusterManager()

	host_ctx := "openmcp"
	namespace := corev1.NamespaceAll

	host_cfg := cm.Host_config

	live := cluster.New(host_ctx, host_cfg, cluster.Options{})

	ghosts := []*cluster.Cluster{}

	for _, ghost_cluster := range cm.Cluster_list.Items {
		ghost_ctx := ghost_cluster.Name
		ghost_cfg := cm.Cluster_configs[ghost_ctx]

		ghost := cluster.New(ghost_ctx, ghost_cfg, cluster.Options{})
		ghosts = append(ghosts, ghost)
	}

	co, _ := controller.NewController(live, ghosts, namespace, cm)
	reshape_cont, _ := reshape.NewController(live, ghosts, namespace, cm)
	loglevel_cont, _ := logLevel.NewController(live, ghosts, namespace)

	m := manager.New()
	m.AddController(co)
	m.AddController(reshape_cont)
	m.AddController(loglevel_cont)

	stop := reshape.SetupSignalHandler()

	if err := m.Start(stop); err != nil {
		log.Fatal(err)
	}

}
