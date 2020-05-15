package resourcecollector

// package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the buf in go mod
	client "github.com/influxdata/influxdb1-client/v2"
	ketiresource "resource-controller/controllers/openmcpscheduler/pkg/controller/resourceInfo"
)

const (
	myDB     = "resource"
	host     = "10.0.3.20:8086"
	username = "admin"
	password = "ketilinux"
)

var c client.Client // influxDB client

func NewInfluxDBClient() {
	c, _ = client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://" + host,
		Username: username,
		Password: password,
	})
	/*
	    if err != nil {
			fmt.Println("[DB] cannot create influxDB client!")
	    }
	*/
	defer c.Close()
}

func ClustersTotalCpuRequest(cluster string) int64 {
	cmd := fmt.Sprintf("select sum(request) from cpu where cluster_name='%s'", cluster)
	res := getTotalResource(cmd)
	return res
}

func ClustersTotalMemoryRequest(cluster string) int64 {
	cmd := fmt.Sprintf("select sum(request) from memory where cluster_name='%s'", cluster)
	res := getTotalResource(cmd)
	return res
}

func getTotalResource(command string) int64 {
	q := client.Query{
		Command:  command,
		Database: myDB,
	}

	resp, err := c.Query(q)
	if err != nil || resp.Error() != nil {
		fmt.Println("[DB] cannot get result")
	}

	res, err := resp.Results[0].Series[0].Values[0][1].(json.Number).Int64()
	if err != nil {
		fmt.Println("[DB] cannot get value from database")
	}
	return res
}

func test() {
	clusterList := make(map[string]*ketiresource.ClusterInfo)

	q := client.Query{
		Command:  "select cluster_name,node_name,pod_name,request from cpu",
		Database: myDB,
	}

	resp, err := c.Query(q)
	if err != nil || resp.Error() != nil {
		fmt.Println("[DB] cannot get result")
	}

	tmp := &ketiresource.ClusterInfo{
		ClusterName: "tmp",
	}
	clusterList["tmp"] = tmp
	fmt.Println(tmp)

	// for _, value := range resp.Results[0].Series[0].Values {
	//     cluster_name := fmt.Sprintf("%s", value[1])
	//     node_name := fmt.Sprintf("%f", value[2])
	//     request, _ := value[4].(json.Number).Int64()

	//     // if already cluster_name is not existed
	//     if value, ok := clusterList[cluster_name]; !ok {
	//         clusterList[cluster_name] = &ketiresource.ClusterInfo {
	//             clusterName:    cluster_name,
	//             // clusterList[cluster_name].nodeList[node_name] = &ketiresource.NodeInfo {
	//             //     nodeName:           node_name,
	//             //     requestedResource:  &ketiresource.Resource {
	//             //         MilliCPU:   request,
	//             //     },
	//             // },
	//         }
	//     } else {
	//         // clustername is already existed, but nodename is not existed
	//         if value2, ok2 := clusterList[cluster_name].nodeList[node_name]; !ok {
	//             clusterList[cluster_name].nodeList[node_name] = &ketiresource.NodeInfo {
	//                 nodeName:           node_name,
	//                 requestedResource:  &ketiresource.Resource {
	//                     MilliCPU:   request,
	//                 },
	//             }
	//         } else {
	//             clusterList[cluster_name].nodeList[node_name].requestedResource.MilliCPU += request
	//         }
	//     }
	// }

	// fmt.Println("[DB] %d", len(resp.Results[0].Series[0].Values))
	if err != nil {
		fmt.Println("[DB] cannot get value from database")
	}
}

// func main(){
//     fmt.Println("[DB] create influxDBClient()")
//     NewInfluxDBClient()
//     // fmt.Println("test: ", getClustersTotalCpuRequest("cluster1"))
//     test()
//     fmt.Println("[DB] finish influxDBClient()")
// }
