package main

import (
	"log"
	"openmcp/openmcp/omcplog"
	openmcpcluster "openmcp/openmcp/openmcp-cluster-manager/src/controller"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/util/controller/logLevel"
	"openmcp/openmcp/util/controller/reshape"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"
)

func main() {
	logLevel.KetiLogInit()

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

	co, err_co := openmcpcluster.NewController(live, ghosts, namespace, cm)
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
