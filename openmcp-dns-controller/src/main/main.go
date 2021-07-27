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
	"openmcp/openmcp/openmcp-dns-controller/src/controller/OpenMCPService"
	"openmcp/openmcp/openmcp-dns-controller/src/controller/dnsEndpoint"
	"openmcp/openmcp/openmcp-dns-controller/src/controller/domain"
	"openmcp/openmcp/openmcp-dns-controller/src/controller/ingressCluster"
	"openmcp/openmcp/openmcp-dns-controller/src/controller/ingressDNSRecord"
	"openmcp/openmcp/openmcp-dns-controller/src/controller/service"
	"openmcp/openmcp/openmcp-dns-controller/src/controller/serviceDNSRecord"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/util/controller/logLevel"
	"openmcp/openmcp/util/controller/reshape"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"
	corev1 "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	logLevel.KetiLogInit()

	cm := clusterManager.NewClusterManager()

	host_ctx := "openmcp"
	namespace := corev1.NamespaceAll

	host_cfg := cm.Host_config
	live := cluster.New(host_ctx, host_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{}})

	ghosts := []*cluster.Cluster{}

	for _, ghost_cluster := range cm.Cluster_list.Items {
		ghost_ctx := ghost_cluster.Name
		ghost_cfg := cm.Cluster_configs[ghost_ctx]

		ghost := cluster.New(ghost_ctx, ghost_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{}})
		ghosts = append(ghosts, ghost)
	}

	cont_service, err := service.NewController(live, ghosts, namespace, cm)
	if err != nil {
		omcplog.V(0).Info("err New Controller - ServiceDNS", err)
	}
	cont_domain, err := domain.NewController(live, ghosts, namespace, cm)
	if err != nil {
		omcplog.V(0).Info("err New Controller - Domain", err)
	}
	cont_serviceDNSRecord, err := serviceDNSRecord.NewController(live, ghosts, namespace, cm)
	if err != nil {
		omcplog.V(0).Info("err New Controller - IngressDNS", err)
	}
	cont_dnsEndpoint, err := dnsEndpoint.NewController(live, ghosts, namespace, cm)
	if err != nil {
		omcplog.V(0).Info("err New Controller - DNSEndpoint", err)
	}
	cont_ingressDNSRecord, err := ingressDNSRecord.NewController(live, ghosts, namespace, cm)
	if err != nil {
		omcplog.V(0).Info("err New Controller - ExternalDNS", err)
	}
	cont_ingressCluster, err := ingressCluster.NewController(live, ghosts, namespace, cm)
	if err != nil {
		omcplog.V(0).Info("err New Controller - OpenMCPService", err)
	}
	cont_OpenMCPService, err := OpenMCPService.NewController(live, ghosts, namespace, cm)
	if err != nil {
		omcplog.V(0).Info("err New Controller - Service", err)
	}
	cont_reshape, err := reshape.NewController(live, ghosts, namespace, cm)
	if err != nil {
		omcplog.V(0).Info("err New Controller - Reshape", err)
	}
	contLoglevel, err := logLevel.NewController(live, ghosts, namespace)
	if err != nil {
		omcplog.V(0).Info("err New Controller - logLevel", err)
	}

	m := manager.New()
	m.AddController(cont_domain)
	m.AddController(cont_service)
	m.AddController(cont_OpenMCPService)
	m.AddController(cont_ingressCluster)
	m.AddController(cont_serviceDNSRecord)
	m.AddController(cont_ingressDNSRecord)
	m.AddController(cont_dnsEndpoint)

	m.AddController(cont_reshape)
	m.AddController(contLoglevel)

	stop := reshape.SetupSignalHandler()

	if err := m.Start(stop); err != nil {
		log.Fatal(err)
	}

}
