package influx

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"log"
	"openmcp-metric-collector/pkg/protobuf"
	"time"
)

type Influx struct {
	inClient client.Client
}

func NewInflux(INFLUX_IP, INFLUX_PORT, username, password string) *Influx {
	inf := &Influx{
		inClient: InfluxDBClient(INFLUX_IP, INFLUX_PORT, username, password),
	}
	return inf
}
func InfluxDBClient(INFLUX_IP, INFLUX_PORT, username, password string) client.Client {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://" + INFLUX_IP + ":" + INFLUX_PORT,
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	return c
}
func (in *Influx) CreateDatabase() {
	q := client.NewQuery("CREATE DATABASE Metrics", "", "")
	if response, err := in.inClient.Query(q); err == nil && response.Error() == nil {
		fmt.Println(response.Results)
	}
}
func (in *Influx) CreateMeasurements() {
	q1 := client.NewQuery("CREATE MEASUREMENTS Nodes", "Metrics", "")
	if response, err := in.inClient.Query(q1); err == nil && response.Error() == nil {
		fmt.Println(response.Results)
	}
	q2 := client.NewQuery("CREATE MEASUREMENTS Pods", "Metrics", "")
	if response, err := in.inClient.Query(q2); err == nil && response.Error() == nil {
		fmt.Println(response.Results)
	}
}
func (in *Influx) SaveMetrics(clusterName string, data *protobuf.Collection) {

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		//Precision:        "rfc3339",
		Database:         "Metrics",
		RetentionPolicy:  "",
		WriteConsistency: "",
	})

	for _, batch := range data.Matricsbatchs {

		nodeName := batch.Node.Name

		tags := map[string]string{
			"cluster": clusterName,
			"node":    nodeName,
			//"region": regions[rand.Intn(len(regions))],
		}
		fields := map[string]interface{}{
			"cpu_usage":        batch.Node.MP.CpuUsage,
			"memory_usage":     batch.Node.MP.MemoryUsage,
			"network_rx_usage": batch.Node.MP.NetworkRxUsage,
			"network_tx_usage": batch.Node.MP.NetworkTxUsage,
			"fs_usage":         batch.Node.MP.FsUsage,
		}
		pt, err := client.NewPoint(
			"Nodes",
			tags,
			fields,
			time.Now(),
		)

		if err != nil {
			fmt.Println("err!", err)
		}

		bp.AddPoint(pt)

		for _, pod := range batch.Pods {
			podName := pod.Name

			tags2 := map[string]string{
				"cluster": clusterName,
				"node":    nodeName,
				"pod":     podName,
				//"region": regions[rand.Intn(len(regions))],
			}
			fields2 := map[string]interface{}{
				"cpu_usage":        pod.MP.CpuUsage,
				"memory_usage":     pod.MP.MemoryUsage,
				"network_rx_usage": pod.MP.NetworkRxUsage,
				"network_tx_usage": pod.MP.NetworkTxUsage,
				"fs_usage":         pod.MP.FsUsage,
			}
			pt2, err := client.NewPoint(
				"Pods",
				tags2,
				fields2,
				time.Now(),
			)
			if err != nil {
				fmt.Println(err)
			}

			bp.AddPoint(pt2)
		}
	}
	in.inClient.Write(bp)

}
