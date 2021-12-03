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
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/util/controller/logLevel"
	"openmcp/openmcp/util/controller/reshape"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	openmcppv "openmcp/openmcp/openmcp-resource-controller/openmcp-pv-controller/src/controller"
)

var cm = clusterManager.NewClusterManager()

func main() {
	logLevel.KetiLogInit()

	omcplog.V(2).Info("OpenMCP PV Controller Start")

	//cm := clusterManager.NewClusterManager()

	host_ctx := "openmcp"
	namespace := "openmcp"

	//openmcp-has-controller 모듈 생성 시 정의한 sa, secret 기반 config
	host_cfg := cm.Host_config
	live := cluster.New(host_ctx, host_cfg, cluster.Options{})
	ghosts := []*cluster.Cluster{}

	for _, ghost_cluster := range cm.Cluster_list.Items {
		ghost_ctx := ghost_cluster.Name
		ghost_cfg := cm.Cluster_configs[ghost_ctx]

		ghost := cluster.New(ghost_ctx, ghost_cfg, cluster.Options{})
		ghosts = append(ghosts, ghost)
	}

	co, err_co := openmcppv.NewController(live, ghosts, namespace, cm)
	if err_co != nil {
		omcplog.V(2).Info("err_co : ", err_co)
		return
	}

	reshape_cont, err_reshape := reshape.NewController(live, ghosts, namespace, cm)
	if err_reshape != nil {
		omcplog.V(2).Info("err_reshape : ", err_reshape)
		return
	}

	loglevel_cont, err_log := logLevel.NewController(live, ghosts, namespace)
	if err_log != nil {
		omcplog.V(2).Info("err_log : ", err_log)
		return
	}

	m := manager.New()
	m.AddController(co)
	m.AddController(reshape_cont)
	m.AddController(loglevel_cont)

	stop := reshape.SetupSignalHandler()

	if err := m.Start(stop); err != nil {
		log.Fatal(err)
	}
}
