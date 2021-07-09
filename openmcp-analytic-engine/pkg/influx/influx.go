package influx

import (
	"fmt"
	"openmcp/openmcp/omcplog"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

type Influx struct {
	inClient client.Client
}

func NewInflux(INFLUX_IP, INFLUX_PORT, username, password string) *Influx {
	omcplog.V(4).Info("Func NewInflux Called")
	inf := &Influx{
		inClient: InfluxDBClient(INFLUX_IP, INFLUX_PORT, username, password),
	}
	return inf
}
func InfluxDBClient(INFLUX_IP, INFLUX_PORT, username, password string) client.Client {
	omcplog.V(4).Info("Func InfluxDBClient Called")
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://" + INFLUX_IP + ":" + INFLUX_PORT,
		Username: username,
		Password: password,
	})
	if err != nil {
		omcplog.V(0).Info("Error: ", err)
	}
	return c
}

func (in *Influx) GetCPAMetricsData(cluster string, namespace string, depname string, podnum string) []client.Result {
	omcplog.V(4).Info("Func GetCPAMetricsData Called")

	q := client.NewQuery("select \"CPUUsageNanoCores\", \"MemoryUsageBytes\", \"NetworkLatency\"  from (select * from Pods where \"cluster\"='"+cluster+"' and \"namespace\"='"+namespace+"' and \"pod\"=~/"+depname+"/ order by time DESC limit "+podnum+") order by time desc", "Metrics", "")

	response, err := in.inClient.Query(q)

	if err == nil && response.Error() == nil {
		return response.Results
	} else {
		fmt.Println(err)
		fmt.Println(response.Error())
	}

	return nil
}
func (in *Influx) GetClusterPodsData(clusterName, podName string) ([]client.Result, error) {
	omcplog.V(4).Info("Func GetClusterPodsData Called")

	//fmt.Println("SELECT CPUUsageNanoCores, MemoryUsageBytes  FROM Pods WHERE cluster = '" + clusterName + "' AND pod = '" + podName + "' ORDER BY DESC LIMIT 1")
	q := client.NewQuery("SELECT CPUUsageNanoCores, MemoryUsageBytes, node FROM Pods WHERE cluster = '"+clusterName+"' AND pod = '"+podName+"' ORDER BY DESC LIMIT 1", "Metrics", "")

	response, err := in.inClient.Query(q)

	if err == nil && response.Error() == nil {
		return response.Results, nil
	} else {
		return nil, err
		// fmt.Println(err)
		// fmt.Println(response.Error())
	}

}
func (in *Influx) GetClusterMetricsData(clusterName string) ([]client.Result, error) {
	omcplog.V(4).Info("Func GetClusterMetricsData Called")
	q := client.NewQuery("SELECT * FROM Nodes WHERE cluster = '"+clusterName+"' GROUP BY * ORDER BY DESC LIMIT 5", "Metrics", "")

	response, err := in.inClient.Query(q)

	if err == nil && response.Error() == nil {
		return response.Results, nil
	} else {
		return nil, err
		// fmt.Println(err)
		// fmt.Println(response.Error())
	}

}

func (in *Influx) SelectMetricsData() []client.Result {
	omcplog.V(4).Info("Func SelectMetricsData Called")
	q := client.NewQuery("select * from Nodes group by * order by desc limit 1", "Metrics", "")

	response, err := in.inClient.Query(q)

	if err == nil && response.Error() == nil {
		return response.Results
	}

	return nil
}

func (in *Influx) GetNetworkData(clusterName, nodeName string) []client.Result {
	omcplog.V(4).Info("Func GetNetworkData Called")
	query_str := "SELECT NetworkRxBytes, NetworkTxBytes FROM Nodes WHERE "
	query_str += "cluster='" + clusterName + "' "
	query_str += "AND node='" + nodeName + "' "
	query_str += "GROUP BY * ORDER BY DESC LIMIT 2"

	q := client.NewQuery(query_str, "Metrics", "")

	response, err := in.inClient.Query(q)

	if err == nil && response.Error() == nil {
		return response.Results
	} else {
		omcplog.V(0).Infof("Cannot get data from InfluxDB: ", err)
		return nil
	}

}

func (in *Influx) InsertClusterStatus(clusterName, time_t string, cpuScore, memScore, netScore, diskScore, latencyScore float64) {
	omcplog.V(4).Info("InsertClusterStatus Called")
	omcplog.V(2).Info("[Save InfluxDB] ClusterName: '", clusterName, "'")

	cpuWarning := false
	memWarning := false
	netWarning := false
	diskWarning := false
	latencyWarning := false

	if cpuScore >= 30 {
		cpuWarning = true
	}
	if memScore >= 30 {
		memWarning = true
	}
	if netScore >= 10 {
		netWarning = true
	}
	if diskScore >= 30 {
		diskWarning = true
	}
	if latencyScore >= 5 {
		latencyWarning = true
	}

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		//Precision:        "rfc3339", // yyyy-MM-ddTHH:mm:ss
		Database:         "Metrics",
		RetentionPolicy:  "",
		WriteConsistency: "",
	})

	tags := map[string]string{
		"cluster": clusterName,
	}

	fields := map[string]interface{}{
		"cpuScore":       cpuScore,
		"memScore":       memScore,
		"netScore":       netScore,
		"diskScore":      diskScore,
		"latencyScore":   latencyScore,
		"cpuWarning":     cpuWarning,
		"memWarning":     memWarning,
		"netWarning":     netWarning,
		"diskWarning":    diskWarning,
		"latencyWarning": latencyWarning,
	}
	t, err := time.Parse(time.RFC3339, time_t)
	if err != nil {
		fmt.Println("err!", err)
	}
	pt, err := client.NewPoint(
		"ClusterStatus",
		tags,
		fields,
		t,
	)

	if err != nil {
		fmt.Println("err!", err)
	}

	bp.AddPoint(pt)

	in.inClient.Write(bp)

}
