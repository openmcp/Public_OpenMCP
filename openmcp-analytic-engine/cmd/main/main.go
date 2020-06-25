package main

import (
	"openmcp/openmcp/openmcp-analytic-engine/pkg/analyticEngine"
	//"openmcp-analytic-engine/pkg/protobuf"
	"os"
	"runtime"
)
const (
	GRPC_PORT = "2050"

)

func main(){
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
