package main

import (
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-globalcache/pkg/controller/globalregistry"
	"openmcp/openmcp/openmcp-globalcache/pkg/controller/noderegistry"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/util/controller/logLevel"
	"openmcp/openmcp/util/controller/reshape"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	logLevel.KetiLogInit()

	for {
		cm := clusterManager.NewClusterManager()
		host_ctx := "openmcp"
		namespace := "openmcp"

		host_cfg := cm.Host_config

		live := cluster.New(host_ctx, host_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{}})

		ghosts := []*cluster.Cluster{}
		for _, ghost_cluster := range cm.Cluster_list.Items {
			ghost_ctx := ghost_cluster.Name
			ghost_cfg := cm.Cluster_configs[ghost_ctx]

			ghost := cluster.New(ghost_ctx, ghost_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{}})

			ghosts = append(ghosts, ghost)
		}
		global, _ := globalregistry.NewController(live, ghosts, namespace, cm)
		node, _ := noderegistry.NewController(live, ghosts, namespace, cm)
		loglevel_cont, _ := logLevel.NewController(live, ghosts, namespace)
		m := manager.New()
		m.AddController(global)
		m.AddController(node)
		m.AddController(loglevel_cont)
		stop := reshape.SetupSignalHandler()
		if err := m.Start(stop); err != nil {
			omcplog.V(4).Info("[OpenMCP] error:  ", err)
		}

	}
}
