package influx

import (
	"github.com/influxdata/influxdb/client/v2"
	"log"
	//"openmcp-analytic-engine/pkg/protobuf"
)

type Influx struct{
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
		Addr:     "http://"+INFLUX_IP+":"+INFLUX_PORT,
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	return c
}

func (in *Influx) GetClusterMetricsData(clusterName string) []client.Result{
	q := client.NewQuery("SELECT * FROM Nodes WHERE cluster = '" + clusterName +"' GROUP BY * ORDER BY DESC LIMIT 2", "Metrics", "")

	response, err := in.inClient.Query(q)

	if err == nil && response.Error() == nil {
		//fmt.Println(response.Results)
		return response.Results
	}

	return nil

}
func (in *Influx) SelectMetricsData() []client.Result{
	q := client.NewQuery("select * from Nodes group by * order by desc limit 1", "Metrics", "")

	response, err := in.inClient.Query(q)

	if err == nil && response.Error() == nil {
		//fmt.Println(response.Results)
		return response.Results
	}

	return nil
}