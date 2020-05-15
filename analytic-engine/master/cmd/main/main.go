package main

import (
	"openmcp-analytic-engine/pkg/analyticEngine"
	"os"
	"runtime"
)
const (
	GRPC_PORT = "2061"

)

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	INFLUX_IP := os.Getenv("INFLUX_IP")
	INFLUX_PORT := os.Getenv("INFLUX_PORT")
	INFLUX_USERNAME := os.Getenv("INFLUX_USERNAME")
	INFLUX_PASSWORD := os.Getenv("INFLUX_PASSWORD")

	//ae := analyticEngine.NewAnalyticEngine()
	ae := analyticEngine.NewAnalyticEngine(INFLUX_IP, INFLUX_PORT, INFLUX_USERNAME, INFLUX_PASSWORD)

	//mc.Influx.CreateDatabase()
	//mc.Influx.CreateMeasurements()

	ae.StartGRPC(GRPC_PORT)
}
