package influx

import (
	"github.com/influxdata/influxdb/client/v2"
	"openmcp/openmcp/omcplog"
	"strings"
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

func (in *Influx) GetPodData(ns, name, clusterName, timeStart, timeEnd string) []client.Result {
	omcplog.V(4).Info("Func GetPodData Called")
	omcplog.V(5).Info("timeStart=",timeStart, ", timeEnd=",timeEnd)

	q := client.Query{}

	if timeStart == "" && timeEnd == "" {
		q = client.NewQuery("SELECT * FROM Pods WHERE pod = '"+name+"' AND namespace = '"+ns+"' AND cluster = '"+clusterName+"' ORDER BY DESC LIMIT 1", "Metrics", "")
	} else if timeStart != "" && timeEnd != ""{
		timeStarts := strings.Split(timeStart, "_")
		timeEnds := strings.Split(timeEnd, "_")
		q = client.NewQuery("SELECT * FROM Pods WHERE pod = '"+name+"' AND namespace = '"+ns+"' AND cluster = '"+clusterName+"' AND time >= '"+timeStarts[0] + " " + timeStarts[1] + "' AND time < '"+timeEnds[0] + " " + timeEnds[1]+"'", "Metrics", "")
		omcplog.V(5).Info("SELECT * FROM Pods WHERE pod = '"+name+"' AND namespace = '"+ns+"' AND cluster = '"+clusterName+"' AND time >= '"+timeStarts[0] + " " + timeStarts[1] + "' AND time < '"+timeEnds[0] + " " + timeEnds[1]+"'")
	} else if timeStart != ""{
		timeStarts := strings.Split(timeStart, "_")
		q = client.NewQuery("SELECT * FROM Pods WHERE pod = '"+name+"' AND namespace = '"+ns+"' AND cluster = '"+clusterName+"' AND time >= '"+timeStarts[0] + " " + timeStarts[1] + "'", "Metrics", "")
		omcplog.V(5).Info("SELECT * FROM Pods WHERE pod = '"+name+"' AND namespace = '"+ns+"' AND cluster = '"+clusterName+"' AND time >= '"+timeStarts[0] + " " + timeStarts[1]+"'" + "'", "Metrics")
	} else if timeEnd != "" {
		timeEnds := strings.Split(timeEnd, "_")
		//select * from Pods WHERE pod='nginx-deployment-55fbd9fd6d-h7d8t' AND namespace='default' AND cluster='cluster1' AND time >= '2020-09-03 00:00:00' and time < '2020-09-03 10:00:05';
		q = client.NewQuery("SELECT * FROM Pods WHERE pod = '"+name+"' AND namespace = '"+ns+"' AND cluster = '"+clusterName+"' AND time < '"+timeEnds[0] + " " + timeEnds[1]+"'", "Metrics", "")
		omcplog.V(5).Info("SELECT * FROM Pods WHERE pod = '"+name+"' AND namespace = '"+ns+"' AND cluster = '"+clusterName+"' AND time < '"+timeEnds[0] + " " + timeEnds[1]+"'")
	}


	response, err := in.inClient.Query(q)

	if err == nil && response.Error() == nil {

		return response.Results
	}


	return nil

}
func (in *Influx) GetNodeData(name, clusterName, timeStart, timeEnd string) []client.Result {
	omcplog.V(4).Info("Func GetNodeData Called")
	omcplog.V(5).Info("timeStart=",timeStart, ", timeEnd=",timeEnd)

	q := client.Query{}

	if timeStart == "" && timeEnd == "" {
		q = client.NewQuery("SELECT * FROM Nodes WHERE node = '"+name+"' AND cluster = '"+clusterName+"' ORDER BY DESC LIMIT 1", "Metrics", "")
	} else if timeStart != "" && timeEnd != ""{
		timeStarts := strings.Split(timeStart, "_")
		timeEnds := strings.Split(timeEnd, "_")
		q = client.NewQuery("SELECT * FROM Nodes WHERE pod = '"+name+"' AND cluster = '"+clusterName+"' AND time >= '"+timeStarts[0] + " " +timeStarts[1]+"' AND time < '"+timeEnds[0] + " " + timeEnds[1]+"'", "Metrics", "")
		omcplog.V(5).Info("SELECT * FROM Nodes WHERE pod = '"+name+"' AND cluster = '"+clusterName+"' AND time >= '"+timeStarts[0] + " " +timeStarts[1]+"' AND time < '"+timeEnds[0] + " " + timeEnds[1]+"'")
	} else if timeStart != ""{
		timeStarts := strings.Split(timeStart, "_")
		q = client.NewQuery("SELECT * FROM Nodes WHERE pod = '"+name+"' AND cluster = '"+clusterName+"' AND time >= '"+timeStarts[0] + " " + timeStarts[1]+"'", "Metrics", "")
		omcplog.V(5).Info("SELECT * FROM Nodes WHERE pod = '"+name+"' AND cluster = '"+clusterName+"' AND time >= '"+timeStarts[0] + " " + timeStarts[1]+"'")
	} else if timeEnd != "" {
		timeEnds := strings.Split(timeEnd, "_")
		//select * from Pods WHERE pod='nginx-deployment-55fbd9fd6d-h7d8t' AND namespace='default' AND cluster='cluster1' AND time >= '2020-09-03 00:00:00' and time < '2020-09-03 10:00:05';
		q = client.NewQuery("SELECT * FROM Nodes WHERE pod = '"+name+"' AND cluster = '"+clusterName+"' AND time < '"+timeEnds[0] + " " + timeEnds[1]+"'", "Metrics", "")
		omcplog.V(5).Info("SELECT * FROM Nodes WHERE pod = '"+name+"' AND cluster = '"+clusterName+"' AND time < '"+timeEnds[0] + " " + timeEnds[1]+"'")
	}
	response, err := in.inClient.Query(q)

	if err == nil && response.Error() == nil {
		return response.Results
	}

	return nil

}
func (in *Influx) GetStatusData(clusterName, timeStart, timeEnd string) []client.Result {
	omcplog.V(4).Info("Func GetNodeData Called")
	omcplog.V(5).Info("timeStart=",timeStart, ", timeEnd=",timeEnd)

	q := client.Query{}

	if timeStart == "" && timeEnd == "" {
		q = client.NewQuery("SELECT * FROM ClusterStatus WHERE cluster = '"+clusterName+"' ORDER BY DESC LIMIT 1", "Metrics", "")
	} else if timeStart != "" && timeEnd != ""{
		timeStarts := strings.Split(timeStart, "_")
		timeEnds := strings.Split(timeEnd, "_")
		q = client.NewQuery("SELECT * FROM ClusterStatus WHERE cluster = '"+clusterName+"' AND time >= '"+timeStarts[0] + " " +timeStarts[1]+"' AND time < '"+timeEnds[0] + " " + timeEnds[1]+"'", "Metrics", "")
		omcplog.V(5).Info("SELECT * FROM ClusterStatus WHERE cluster = '"+clusterName+"' AND time >= '"+timeStarts[0] + " " +timeStarts[1]+"' AND time < '"+timeEnds[0] + " " + timeEnds[1]+"'")
	} else if timeStart != ""{
		timeStarts := strings.Split(timeStart, "_")
		q = client.NewQuery("SELECT * FROM ClusterStatus WHERE cluster = '"+clusterName+"' AND time >= '"+timeStarts[0] + " " + timeStarts[1]+"'", "Metrics", "")
		omcplog.V(5).Info("SELECT * FROM ClusterStatus WHERE cluster = '"+clusterName+"' AND time >= '"+timeStarts[0] + " " + timeStarts[1]+"'")
	} else if timeEnd != "" {
		timeEnds := strings.Split(timeEnd, "_")
		//select * from Pods WHERE pod='nginx-deployment-55fbd9fd6d-h7d8t' AND namespace='default' AND cluster='cluster1' AND time >= '2020-09-03 00:00:00' and time < '2020-09-03 10:00:05';
		q = client.NewQuery("SELECT * FROM ClusterStatus WHERE cluster = '"+clusterName+"' AND time < '"+timeEnds[0] + " " + timeEnds[1]+"'", "Metrics", "")
		omcplog.V(5).Info("SELECT * FROM ClusterStatus WHERE cluster = '"+clusterName+"' AND time < '"+timeEnds[0] + " " + timeEnds[1]+"'")
	}
	response, err := in.inClient.Query(q)

	if err == nil && response.Error() == nil {
		return response.Results
	}

	return nil

}












//
//
//func (in *Influx) GetClusterMetricsData(clusterName string) []client.Result {
//	omcplog.V(4).Info("Func GetClusterMetricsData Called")
//	q := client.NewQuery("SELECT * FROM Nodes WHERE cluster = '"+clusterName+"' GROUP BY * ORDER BY DESC LIMIT 2", "Metrics", "")
//
//	response, err := in.inClient.Query(q)
//
//	if err == nil && response.Error() == nil {
//		return response.Results
//	}
//
//	return nil
//
//}
//
//func (in *Influx) SelectMetricsData() []client.Result {
//	omcplog.V(4).Info("Func SelectMetricsData Called")
//	q := client.NewQuery("select * from Nodes group by * order by desc limit 1", "Metrics", "")
//
//	response, err := in.inClient.Query(q)
//
//	if err == nil && response.Error() == nil {
//		return response.Results
//	}
//
//	return nil
//}
//
//func (in *Influx) GetNetworkData(clusterName, nodeName string) []client.Result {
//	omcplog.V(4).Info("Func GetNetworkData Called")
//	query_str := "SELECT NetworkRxBytes, NetworkTxBytes FROM Nodes WHERE "
//	query_str += "cluster='" + clusterName + "' "
//	query_str += "AND node='" + nodeName + "' "
//	query_str += "GROUP BY * ORDER BY DESC LIMIT 2"
//
//	q := client.NewQuery(query_str, "Metrics", "")
//
//	response, err := in.inClient.Query(q)
//
//	if err == nil && response.Error() == nil {
//		return response.Results
//	}else {
//		omcplog.V(0).Infof("Cannot get data from InfluxDB: ", err)
//		return nil
//	}
//
//	return nil
//}