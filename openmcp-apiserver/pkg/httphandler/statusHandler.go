package httphandler

import (
	"bytes"
	"encoding/json"
	"github.com/influxdata/influxdb/client/v2"
	"net/http"
	"openmcp/openmcp/openmcp-apiserver/pkg/influx"
	"os"
	"strconv"
)

type ClusterStatusList struct {
	Items []ClusterStatus `json:"cluster_status"`
}

type ClusterStatus struct {
	Time      		string      `json:"time"`
	Cluster  	 	string      `json:"cluster"`

	CpuScore  		string		`json:"cpu_score"`
	MemoryScore 	string 		`json:"memory_score"`
	DiskScore  		string		`json:"disk_score"`
	NetworkScore	string		`json:"network_score"`
	LatencyScore	string		`json:"latency_score"`

	CpuWarning  	string		`json:"cpu_warning"`
	MemoryWarning 	string 		`json:"memory_warning"`
	DiskWarning  	string		`json:"disk_warning"`
	NetworkWarning	string		`json:"network_warning"`
	LatencyWarning	string		`json:"latency_warning"`

}

func (h *HttpManager)StatusHandler(w http.ResponseWriter, r *http.Request) {
	// example Node & Pod in cluster1
	// GET http://10.0.3.20:31635/status?clustername=cluster1
	// GET http://10.0.3.20:31635/status?clustername=cluster1&timeStart=2020-09-03_09:00:00
	// GET http://10.0.3.20:31635/status?clustername=cluster1&timeEnd=2020-09-03_09:00:15
	// GET http://10.0.3.20:31635/status?clustername=cluster1&timeStart=2020-09-03_09:00:00&timeEnd=2020-09-03_09:00:15



	clusterNames, ok := r.URL.Query()["clustername"]
	if !ok || len(clusterNames[0]) < 1 {
		w.Write([]byte("Url Param 'clustername' is missing"))
		return
	}
	clusterName := clusterNames[0]

	timeStart := ""
	timeStarts, ok := r.URL.Query()["timeStart"]
	if !ok || len(timeStarts[0]) < 1 {
		timeStart = ""
	} else {
		timeStart = timeStarts[0]
	}
	timeEnd := ""
	timeEnds, ok := r.URL.Query()["timeEnd"]
	if !ok || len(timeEnds[0]) < 1 {
		timeEnd = ""
	} else {
		timeEnd = timeEnds[0]
	}


	jsonByteArray:= getStatus(clusterName, timeStart, timeEnd)


	w.Write(jsonByteArray)
}

func getStatus(clusterName, timeStart, timeEnd string) []byte {
	INFLUX_IP := os.Getenv("INFLUX_IP")
	INFLUX_PORT := os.Getenv("INFLUX_PORT")
	INFLUX_USERNAME := os.Getenv("INFLUX_USERNAME")
	INFLUX_PASSWORD := os.Getenv("INFLUX_PASSWORD")
	inf := influx.NewInflux(INFLUX_IP, INFLUX_PORT, INFLUX_USERNAME, INFLUX_PASSWORD)


	results := inf.GetStatusData(clusterName, timeStart, timeEnd)
	cs := setClusterStatus(results)

	bytesJson, _ := json.Marshal(cs)
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, bytesJson, "", "\t")
	if err != nil {
		panic(err.Error())
	}

	return prettyJSON.Bytes()




}

func setClusterStatus(results []client.Result) *ClusterStatusList {
	csList := &ClusterStatusList{}
	for _, result := range results{
		for _, ser := range result.Series {

			for r, _ := range ser.Values {
				cs := &ClusterStatus{}
				for c, colName := range ser.Columns {
					value := ser.Values[r][c]
					if colName == "time"{
						cs.Time = value.(string)
					} else if colName == "cluster" {
						cs.Cluster = value.(string)
					}else if colName == "cpuScore" {
						cs.CpuScore = value.(json.Number).String()
					}else if colName == "cpuWarning" {
						cs.CpuWarning = strconv.FormatBool(value.(bool))
					}else if colName == "diskScore" {
						cs.DiskScore = value.(json.Number).String()
					}else if colName == "diskWarning" {
						cs.DiskWarning = strconv.FormatBool(value.(bool))
					}else if colName == "latencyScore" {
						cs.LatencyScore = value.(json.Number).String()
					}else if colName == "latencyWarning" {
						cs.LatencyWarning = strconv.FormatBool(value.(bool))
					}else if colName == "memScore" {
						cs.MemoryScore = value.(json.Number).String()
					}else if colName == "memWarning" {
						cs.MemoryWarning= strconv.FormatBool(value.(bool))
					}else if colName == "netScore" {
						cs.NetworkScore= value.(json.Number).String()
					}else if colName == "netWarning" {
						cs.NetworkWarning= strconv.FormatBool(value.(bool))
					}

				}
				csList.Items = append(csList.Items, *cs)
			}


		}
	}



	return csList

}
