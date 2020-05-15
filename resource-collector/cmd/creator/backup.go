package backup

import (
	"context"
	"fmt"
	crm "github.com/hth0919/resourcecollector"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"time"
	"encoding/json"
)

const (
	port = ":50051"
)
var cs []crm.ClusterInfo
var host string
var curcon =0
var sleeptime = time.Second
var con []*grpc.ClientConn
var c []crm.GetClusterClient
type YamlConfig  struct {
	APIVersion string `yaml:"apiVersion"`
	Clusters   []struct {
		Cluster struct {
			CertificateAuthorityData string `yaml:"certificate-authority-data"`
			Server                   string `yaml:"server"`
		} `yaml:"cluster"`
		Name string `yaml:"name"`
	} `yaml:"clusters"`
	Contexts []struct {
		Context struct {
			Cluster string `yaml:"cluster"`
			User    string `yaml:"user"`
		} `yaml:"context"`
		Name string `yaml:"name"`
	} `yaml:"contexts"`
	CurrentContext string `yaml:"current-context"`
	Kind           string `yaml:"kind"`
	Preferences    struct {
	} `yaml:"preferences"`
	Users []struct {
		Name string `yaml:"name"`
		User struct {
			ClientCertificateData string `yaml:"client-certificate-data"`
			ClientKeyData         string `yaml:"client-key-data"`
		} `yaml:"user"`
	} `yaml:"users"`
}

func main() {
	var yc YamlConfig
	hostset := []string{"cluster-resource-collector.openmcp.hth-domain.svc.china.asia.hthtest.org"+port, "cluster-resource-collector.openmcp.hth-domain.svc.northamerica.hthtest.org"+port, "cluster-resource-collector.openmcp.hth-domain.svc.europe.hthtest.org"+port}

	yc = yc.yamlunmarshal()
	for i:= 0 ;i<len(hostset);i++ {
		host := hostset[i]
		//con = append(con, connecthost(host))
		con = append(con, connecthost(host))
		fmt.Println("connect host",i,":::",con)

	}

	for i:= 0 ;i<len(con);i++ {
		defer con[i].Close()
		c = append(c, crm.NewGetClusterClient(con[i]))
		fmt.Println("connect client",i,":::",c)
	}

	for i:= 0 ;i<10;i++ {
		fmt.Println(i,"::::")
		for j := 0; j < len(c); j++ {
			fmt.Println(j)
			var ci *crm.ClusterInfo
			now := time.Now()
			fmt.Println(now)
			ci = &crm.ClusterInfo{
				MetricValue:      nil,
				Clustername:      yc.Clusters[j].Name,
				KubeConfig:       "",
				AdminToken:       "",
				NodeList:         nil,
				ClusterMetricSum: nil,
				Host:             host,
			}
			go fromcluster(j, ci)
		}
		time.Sleep(time.Second*5)
	}
}

func fromcluster(index int, ci *crm.ClusterInfo) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	r, err := c[index].GetCluster(ctx, ci)
	if err != nil {
		fmt.Printf("could not connect : %v", err)
	}
	data, err := json.MarshalIndent(r, "", "	")
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(r.Clustername+".txt", data, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Clustername)
	fmt.Println("Cluster", index+1, ":::::::::::: DONE")
}


func gethost(server string) string{
	temp := server
	hs := strings.Split(temp,"/")
	host := hs[2]
	t := strings.Split(host,":")
	host = t[0] + port
	return host
}


func connecthost(host string) *grpc.ClientConn {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("did not connect: %v", err)
	}
	return conn
}

func (yc *YamlConfig)yamlunmarshal() YamlConfig{
	fmt.Println("Parsing YAML file")
	var yamlConfig YamlConfig

	var kubeconfig string
	kubeconfig = homeDir()

	yamlFile, err := ioutil.ReadFile(kubeconfig)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return yamlConfig
	}

	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
		return yamlConfig
	}
	fmt.Println("parsing done")
	return yamlConfig
}

func homeDir() string {
	return "/usr/local/bin/config"
}
