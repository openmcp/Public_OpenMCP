package influx

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"log"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-metric-collector/master/pkg/protobuf"
	"time"
)

type Influx struct {
	inClient client.Client
}

func NewInflux(INFLUX_IP, INFLUX_PORT, username, password string) *Influx {
	omcplog.V(4).Info("NewInflux Called")
	inf := &Influx{
		inClient: InfluxDBClient(INFLUX_IP, INFLUX_PORT, username, password),
	}
	return inf
}
func InfluxDBClient(INFLUX_IP, INFLUX_PORT, username, password string) client.Client {
	omcplog.V(4).Info("InfluxDBClient Called")
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
	omcplog.V(4).Info("CreateDatabase Called")
	q := client.NewQuery("CREATE DATABASE Metrics", "", "")
	if response, err := in.inClient.Query(q); err == nil && response.Error() == nil {
		fmt.Println(response.Results)
	}
}
func (in *Influx) CreateMeasurements() {
	omcplog.V(4).Info("CreateMeasurements Called")
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
	omcplog.V(4).Info("SaveMetrics Called")
	omcplog.V(2).Info("[Save InfluxDB] ClusterName: '", clusterName,"'")

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		//Precision:        "rfc3339", // yyyy-MM-ddTHH:mm:ss
		Database:         "Metrics",
		RetentionPolicy:  "",
		WriteConsistency: "",
	})

	time_t := time.Now()
	for _, batch := range data.Metricsbatchs {

		nodeName := batch.Node.Name

		tags := map[string]string{
			"cluster": clusterName,
			"node":    nodeName,
			//"region": regions[rand.Intn(len(regions))],
		}

		fields := map[string]interface{}{
			"CPUUsageNanoCores": batch.Node.MP.CPUUsageNanoCores,
			"MemoryAvailableBytes": batch.Node.MP.MemoryAvailableBytes,
			"MemoryUsageBytes": batch.Node.MP.MemoryUsageBytes,
			"MemoryWorkingSetBytes": batch.Node.MP.MemoryWorkingSetBytes,
			"NetworkRxBytes": batch.Node.MP.NetworkRxBytes,
			"NetworkTxBytes": batch.Node.MP.NetworkTxBytes,
			"FsAvailableBytes": batch.Node.MP.FsAvailableBytes,
			"FsCapacityBytes": batch.Node.MP.FsCapacityBytes,
			"FsUsedBytes": batch.Node.MP.FsUsedBytes,
			"NetworkLatency" : batch.Node.MP.NetworkLatency,
		}
		pt, err := client.NewPoint(
			"Nodes",
			tags,
			fields,
			time_t,
		)

		if err != nil {
			fmt.Println("err!", err)
		}

		bp.AddPoint(pt)
		for _, pod := range batch.Pods {
			podName := pod.Name
			podNamespace := pod.Namespace

			tags2 := map[string]string{
				"cluster": clusterName,
				"node":    nodeName,
				"pod":     podName,
				"namespace": podNamespace,
				//"region": regions[rand.Intn(len(regions))],
			}
			fields2 := map[string]interface{}{
				"CPUUsageNanoCores": pod.MP.CPUUsageNanoCores,
				"MemoryAvailableBytes": pod.MP.MemoryAvailableBytes,
				"MemoryUsageBytes": pod.MP.MemoryUsageBytes,
				"MemoryWorkingSetBytes":pod.MP.MemoryWorkingSetBytes,
				"NetworkRxBytes": pod.MP.NetworkRxBytes,
				"NetworkTxBytes": pod.MP.NetworkTxBytes,
				"FsAvailableBytes": pod.MP.FsAvailableBytes,
				"FsCapacityBytes": pod.MP.FsCapacityBytes,
				"FsUsedBytes": pod.MP.FsUsedBytes,
				"NetworkLatency" : pod.MP.NetworkLatency,
			}
			pt2, err := client.NewPoint(
				"Pods",
				tags2,
				fields2,
				time_t,
			)
			if err != nil {
				fmt.Println(err)
			}

			bp.AddPoint(pt2)
		}
	}
	in.inClient.Write(bp)

}

