package priorities

import (
	// "context"
	// "time"
	"strconv"
	"fmt"
	"k8s.io/klog"
	ketiresource "openmcpscheduler/pkg/controller/resourceinfo"
	_ "github.com/influxdata/influxdb1-client"  // this is important because of the buf in go mod
	client "github.com/influxdata/influxdb1-client/v2"
)

type BalancedNetworkAllocation struct{}

const (
	database = "Metrics"
)

func (pl *BalancedNetworkAllocation) Name() string {
	return "BalancedNetworkAllocation"
}

func (pl *BalancedNetworkAllocation) Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) int64 {
	var clutserScore int64

	// Get InfluxDB from openmcp's InfluxDB Pod
	c := influxDBClient()

	for _, node := range clusterInfo.Nodes {	
		// Get Node's current_network_receive_bytes_total, current_network_transmit_bytes_total
		// by using gRPC from resource collector. 
		res, err := getNetworkFromInfluxDB(c, node.NodeName)
		if err != nil{
			klog.Info("Cannot get network information from InfluxDB")
			return 0
		}
		for _, pod := range res[0].Series{
			rx, _ := strconv.ParseFloat(fmt.Sprintf("%s", pod.Values[0][1]), 64)
			tx, _ := strconv.ParseFloat(fmt.Sprintf("%s", pod.Values[0][2]), 64)
			rx = rx / 1000 / 1000 // change to mega
			tx = tx / 1000 / 1000 // change to mega
			podScore := int64(1 / (rx + tx)) * maxScore
			clutserScore += podScore
		}
	}

	return clutserScore
}

func influxDBClient() client.Client {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://10.244.2.53:8086",
	})
	if err != nil{
		klog.Info("Create InfluxDBClient Error")
	}
	defer c.Close()

	return c
}

func getNetworkFromInfluxDB(c client.Client, nodeName string) (res[] client.Result, err error) {
	q := client.Query{
		Command: fmt.Sprintf("select network_rx_usage, network_tx_usage from Pods where node='"+nodeName+"' group by * order by desc limit 1"),
		// Command: fmt.Sprintf("select * from Pods where node='"+nodeName+"' group by * order by desc limit 1"),
		Database: database,
	}

	if response, err := c.Query(q); err == nil{
		if response.Error() != nil{
			klog.Info("getting network data from influxdb is error..")
			return res, response.Error()
		}
		res = response.Results
	} else {
		klog.Info("getting network data from influxdb is error..")
		return res, err
	}
	return res, nil
}

// The content of *.proto file

// service OpenMCPSchedulerResource {
// 	rpc GetCurrentNetwork (Node) returns (Network) {}
// }

// message Node {
// 	string nodeName = 1;
// 	string clusterName = 2;
// }

// message Network {
// 	float64 currentReceiveBytes = 1;
// 	float64 currentTransmitBytes = 2;	
// }

// message Affinity {
// 	string region = 1;
// 	string zone = 2; 
// }