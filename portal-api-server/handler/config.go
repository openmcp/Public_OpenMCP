package handler

import (
	"fmt"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/jinzhu/configor"
)

var portalConfig = struct {
	Portal struct {
		OpenmcpURL         string
		Port               string
		Kubeconfig         string
		OpenmcpClusterName string
	}
}{}

func InitPortalConfig() string {
	configor.Load(&portalConfig, "portalConfig.yml")
	return portalConfig.Portal.OpenmcpURL + ":" + portalConfig.Portal.Port
}

var openmcpURL = InitPortalConfig()
var kubeConfigFile = portalConfig.Portal.Kubeconfig
var openmcpAddress = portalConfig.Portal.OpenmcpURL
var openmcpClusterName = portalConfig.Portal.OpenmcpClusterName

var InfluxConfig = struct {
	Influx struct {
		Ip       string
		Port     string
		Username string
		Password string
	}
}{}

//Influx Configration
func InitInfluxConfig() {
	configor.Load(&InfluxConfig, "dbconfig.yml")
}

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
		// InsecureSkipVerify: true,
	})
	if err != nil {
		fmt.Println(err)
	}
	return c
}

//Influx Configration
