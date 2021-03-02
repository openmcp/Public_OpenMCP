package main

import (
	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"
	"log"
	openmcpcluster "openmcp/openmcp/openmcp-cluster-manager/pkg/controller"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/util/controller/logLevel"
	"openmcp/openmcp/util/controller/reshape"
)

func main() {
	logLevel.KetiLogInit()

	for {
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
		co, _ := openmcpcluster.NewController(live, ghosts, namespace, cm)
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
}
