package main

import (
	"log"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-analytic-engine/src/analyticEngine"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/util/controller/logLevel"
	"openmcp/openmcp/util/controller/reshape"
	"os"
	"runtime"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"
)

const (
	GRPC_PORT = "2050"
)

func main() {
	logLevel.KetiLogInit()

	for {
		cm := clusterManager.NewClusterManager()

		quit := make(chan bool)
		quitok := make(chan bool)
		go AnalyticEngine(cm, quit, quitok)

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

		/*
			fmt.Println(live)
			fmt.Println(ghosts)
			fmt.Println(namespace)
		*/

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
		m.AddController(reshape_cont)
		m.AddController(loglevel_cont)

		stop := reshape.SetupSignalHandler()

		if err := m.Start(stop); err != nil {
			log.Fatal(err)
		}
		quit <- true
		quit <- true
		<-quitok
		<-quitok
		//time.Sleep(3600 * time.Second)

	}

}

func AnalyticEngine(cm *clusterManager.ClusterManager, quit, quitok chan bool) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	INFLUX_IP := os.Getenv("INFLUX_IP")
	INFLUX_PORT := os.Getenv("INFLUX_PORT")
	INFLUX_USERNAME := os.Getenv("INFLUX_USERNAME")
	INFLUX_PASSWORD := os.Getenv("INFLUX_PASSWORD")

	//ae := analyticEngine.NewAnalyticEngine()
	ae := analyticEngine.NewAnalyticEngine(INFLUX_IP, INFLUX_PORT, INFLUX_USERNAME, INFLUX_PASSWORD)

	go ae.CalcResourceScore(cm, quit, quitok)

	//a := protobuf.HASInfo{HPANamespace:"openmcp", HPAName:"openmcp-hpa", ClusterName:""}

	//ae.SelectHPACluster(&a)
	//mc.Influx.CreateDatabase()
	//mc.Influx.CreateMeasurements()
	go func() {
		ae.StartGRPC(GRPC_PORT)
	}()

	if <-quit {
		ae.StopGRPC()
		quitok <- true
	}

}
