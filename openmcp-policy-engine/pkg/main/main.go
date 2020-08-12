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
	"openmcp/openmcp/omcplog"

	//"flag"
	"log"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/util/controller/logLevel"

	"admiralty.io/multicluster-controller/pkg/cluster"
	//"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/manager"
	//"admiralty.io/multicluster-controller/pkg/reconcile"
	//"admiralty.io/multicluster-service-account/pkg/config"
	//"k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	//"k8s.io/sample-controller/pkg/signals"

	"openmcp/openmcp/openmcp-policy-engine/pkg/controller"
	"openmcp/openmcp/util/controller/reshape"
	//"k8s.io/client-go/rest"
	//genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
	//fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	//"sigs.k8s.io/kubefed/pkg/controller/util"
)

func main() {
	logLevel.KetiLogInit()

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

		co, err := controller.NewController(live, ghosts, namespace, cm)
		if err != nil {
			omcplog.V(0).Info( "err New Controller - PolicyEngine", err)
		}
		cont_reshape, err := reshape.NewController(live, ghosts, namespace)
		if err != nil {
			omcplog.V(0).Info( "err New Controller - Reshape", err)
		}
		contLoglevel, err := logLevel.NewController(live, ghosts, namespace)
		if err != nil {
			omcplog.V(0).Info( "err New Controller - logLevel", err)
		}

		m := manager.New()
		m.AddController(co)
		m.AddController(cont_reshape)
		m.AddController(contLoglevel)

		stop := reshape.SetupSignalHandler()

		if err := m.Start(stop); err != nil {
			log.Fatal(err)
		}
	}

}