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
	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"
	"fmt"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"log"
	"openmcp/openmcp/openmcp-dns-controller/pkg/controller/dnsEndpoint"
	"openmcp/openmcp/openmcp-dns-controller/pkg/controller/domain"
	"openmcp/openmcp/openmcp-dns-controller/pkg/controller/externalDNS"
	"openmcp/openmcp/openmcp-dns-controller/pkg/controller/ingressDNS"
	"openmcp/openmcp/openmcp-dns-controller/pkg/controller/serviceDNS"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/util/controller/reshape"
)

func main() {
	for {
		cm := clusterManager.NewClusterManager()

		host_ctx := "openmcp"
		namespace := "openmcp"

		host_cfg := cm.Host_config
		live := cluster.New(host_ctx, host_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{}})
		//live := cluster.New(host_ctx, host_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{Namespace: namespace}})

		ghosts := []*cluster.Cluster{}

		for _, ghost_cluster := range cm.Cluster_list.Items {
			ghost_ctx := ghost_cluster.Name
			ghost_cfg := cm.Cluster_configs[ghost_ctx]

			ghost := cluster.New(ghost_ctx, ghost_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{}})
			//ghost := cluster.New(ghost_ctx, ghost_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{Namespace: namespace}})
			ghosts = append(ghosts, ghost)
		}

		cont_serviceDNS, err := serviceDNS.NewController(live, ghosts, namespace)
		if err != nil {
			fmt.Println("err New Controller - ServiceDNS", err)
		}
		cont_domain, err := domain.NewController(live, ghosts, namespace)
		if err != nil {
			fmt.Println("err New Controller - Domain", err)
		}
		cont_ingressDNS, err := ingressDNS.NewController(live, ghosts, namespace)
		if err != nil {
			fmt.Println("err New Controller - IngressDNS", err)
		}
		cont_dnsEndpoint, err := dnsEndpoint.NewController(live, ghosts, namespace)
		if err != nil {
			fmt.Println("err New Controller - DNSEndpoint", err)
		}
		cont_externalDNS, err := externalDNS.NewController(live, ghosts, namespace)
		if err != nil {
			fmt.Println("err New Controller - ExternalDNS", err)
		}
		cont_reshape, err := reshape.NewController(live, ghosts, namespace)
		if err != nil {
			fmt.Println("err New Controller - Reshape", err)
		}

		m := manager.New()
		m.AddController(cont_serviceDNS)
		m.AddController(cont_domain)
		m.AddController(cont_ingressDNS)
		m.AddController(cont_dnsEndpoint)
		m.AddController(cont_externalDNS)
		m.AddController(cont_reshape)

		stop := reshape.SetupSignalHandler()
		//stop := signals.SetupSignalHandler()

		if err := m.Start(stop); err != nil {
			log.Fatal(err)
		}

	}


}
