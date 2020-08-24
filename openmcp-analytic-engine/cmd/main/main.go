package main

import (
	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"
	"log"
	"openmcp/openmcp/openmcp-analytic-engine/pkg/analyticEngine"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/util/controller/logLevel"
	"openmcp/openmcp/util/controller/reshape"

	//"openmcp-analytic-engine/pkg/protobuf"
	"os"
	"runtime"
)
const (
	GRPC_PORT = "2050"
)

func main() {
	logLevel.KetiLogInit()

	go AnalyticEngine()

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

		reshape_cont, _ := reshape.NewController(live, ghosts, namespace)
		loglevel_cont, _ := logLevel.NewController(live, ghosts, namespace)

		m := manager.New()
		m.AddController(reshape_cont)
		m.AddController(loglevel_cont)

		stop := reshape.SetupSignalHandler()

		if err := m.Start(stop); err != nil {
			log.Fatal(err)
		}
	}

}

func AnalyticEngine(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	INFLUX_IP := os.Getenv("INFLUX_IP")
	INFLUX_PORT := os.Getenv("INFLUX_PORT")
	INFLUX_USERNAME := os.Getenv("INFLUX_USERNAME")
	INFLUX_PASSWORD := os.Getenv("INFLUX_PASSWORD")

	//ae := analyticEngine.NewAnalyticEngine()
	ae := analyticEngine.NewAnalyticEngine(INFLUX_IP, INFLUX_PORT, INFLUX_USERNAME, INFLUX_PASSWORD)


	go ae.CalcResourceScore()

	//a := protobuf.HASInfo{HPANamespace:"openmcp", HPAName:"openmcp-hpa", ClusterName:""}


	//ae.SelectHPACluster(&a)
	//mc.Influx.CreateDatabase()
	//mc.Influx.CreateMeasurements()

	ae.StartGRPC(GRPC_PORT)
}
