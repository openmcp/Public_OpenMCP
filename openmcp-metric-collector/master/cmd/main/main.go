package main

import (
	"openmcp/openmcp/openmcp-metric-collector/master/pkg/metricCollector"
	"os"

	"runtime"
)

const (
	GRPC_PORT = "2051"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	INFLUX_IP := os.Getenv("INFLUX_IP")
	INFLUX_PORT := os.Getenv("INFLUX_PORT")
	INFLUX_USERNAME := os.Getenv("INFLUX_USERNAME")
	INFLUX_PASSWORD := os.Getenv("INFLUX_PASSWORD")

	mc := metricCollector.NewMetricCollector(INFLUX_IP, INFLUX_PORT, INFLUX_USERNAME, INFLUX_PASSWORD)

	mc.Influx.CreateDatabase()
	mc.Influx.CreateMeasurements()

	mc.StartGRPC(GRPC_PORT)

	//mc := &metricCollector.MetricCollector{}
	//mc.StartGRPC(GRPC_PORT)

}
