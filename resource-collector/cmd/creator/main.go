package main

import (
	"context"
	"fmt"
	"bytes"
	"sync"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	res "github.com/hth0919/resourcecollector"
	"net"
	"runtime"
	"sigs.k8s.io/kubefed/pkg/controller/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
	"time"
	"net/http"
	"strings"
)

const (
	port = ":50051"
)

type server struct{}
var ClusterList = make(map[string]*res.ClusterInfo)
var ClusterVersion = make(map[string]int)
var Result = &res.ReturnValue{
	Tick:                 0,
	ClusterName:          "",
}
var wg sync.WaitGroup
type BulkInsert struct {
	Measurementname string
	Clusters []string
	Nodes []string
	Pods []string
	Metricnames []string
	Values []string
	TimeStamps []int64
	PostLines []string
	HostString string
}
var manager *ClusterManager
type ClusterManager struct {
	Fed_namespace string
	Host_config *rest.Config
	Host_client genericclient.Client
	Cluster_list *fedv1b1.KubeFedClusterList
	Cluster_configs map[string]*rest.Config
	Cluster_clients map[string]genericclient.Client
	Kubeconfig map[string]*kubernetes.Clientset
}

func ListKubeFedClusters(client genericclient.Client, namespace string) *fedv1b1.KubeFedClusterList {
	clusterList := &fedv1b1.KubeFedClusterList{}
	err := client.List(context.TODO(), clusterList, namespace)
	if err != nil {
		fmt.Println("Error retrieving list of federated clusters: %+v", err)
	}
	if len(clusterList.Items) == 0 {
		fmt.Println("No federated clusters found")
	}
	return clusterList
}

func KubeFedClusterConfigs(clusterList *fedv1b1.KubeFedClusterList, client genericclient.Client, fedNamespace string) map[string]*rest.Config {
	clusterConfigs := make(map[string]*rest.Config)
	for _, cluster := range clusterList.Items {
		config, _ := util.BuildClusterConfig(&cluster, client, fedNamespace)
		clusterConfigs[cluster.Name] = config
	}
	return clusterConfigs
}
func KubeFedClusterClients(clusterList *fedv1b1.KubeFedClusterList, cluster_configs map[string]*rest.Config) map[string]genericclient.Client {

	cluster_clients := make(map[string]genericclient.Client)
	for _, cluster := range clusterList.Items {
		clusterName := cluster.Name
		cluster_config := cluster_configs[clusterName]
		cluster_client := genericclient.NewForConfigOrDie(cluster_config)
		cluster_clients[clusterName] = cluster_client
	}
	return cluster_clients
}

func Kubeconfigs(clusterList *fedv1b1.KubeFedClusterList, cluster_configs map[string]*rest.Config)  map[string]*kubernetes.Clientset {
	kube_clientes := make(map[string]*kubernetes.Clientset)
	for _, cluster := range clusterList.Items {
		clusterName := cluster.Name
		ClusterList[clusterName] = &res.ClusterInfo{
			MetricValue:      []string{},
			Clustername:      "",
			KubeConfig:       "",
			AdminToken:       "",
			NodeList:         []*res.NodeInfo{},
			ClusterMetricSum: map[string]float64{},
			Host:             "",
			Pods: 			  []string{},
		}
		ClusterVersion[clusterName] = 0
		client, err := kubernetes.NewForConfig(cluster_configs[clusterName])
		if err != nil {
			panic(err.Error())
		}
		kube_clientes[clusterName] = client
	}
	return kube_clientes
}

func NewClusterManager() *ClusterManager {
	fed_namespace := "kube-federation-system"
	host_config, _ := rest.InClusterConfig()
	host_client := genericclient.NewForConfigOrDie(host_config)
	cluster_list := ListKubeFedClusters(host_client, fed_namespace)
	cluster_configs := KubeFedClusterConfigs(cluster_list, host_client, fed_namespace)
	cluster_clients := KubeFedClusterClients(cluster_list, cluster_configs)
	kube_clientes := Kubeconfigs(cluster_list, cluster_configs)

	cm := &ClusterManager{
		Fed_namespace: fed_namespace,
		Host_config: host_config,
		Host_client: host_client,
		Cluster_list: cluster_list,
		Cluster_configs: cluster_configs,
		Cluster_clients: cluster_clients,
		Kubeconfig: kube_clientes,
	}
	return cm
}
func SetClusterName(pn string) string {
	// creates the in-cluster config
	cn := ""
	for _, cluster := range manager.Cluster_list.Items {
		clusterName := cluster.Name
		pods, err := manager.Kubeconfig[clusterName].CoreV1().Pods("rescollect").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		if len(pods.Items)>0 {
			if pods.Items[0].Name == pn {
				cn = clusterName
			}
		}
	}
	return cn
}
func (s *server) SendCluster(ctx context.Context, in *res.ClusterInfo) (*res.ReturnValue, error) {
	startTime := time.Now()
	log.Println("Recieved : ", in.Host, "::::", in.Clustername )
	timetick := 3600
	clustername := ""
	if len(in.Clustername)==0{
		clustername = SetClusterName(in.Pods[0])
		in.Clustername = clustername
		timetick = 3
	}
	clustername = in.Clustername
	ClusterList[clustername] = in
	ClusterVersion[clustername]++
	amain(in, clustername)

	elapsedTime := time.Since(startTime)
	fmt.Println(clustername, "::::", elapsedTime)
	return &res.ReturnValue{
		Tick:                 int64(timetick),
		ClusterName:          clustername,
	}, nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Println("grpcGiver start at port %s", port)

	manager = NewClusterManager()
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	res.RegisterSendClusterServer(s,&server{})
	if err := s.Serve(l); err!=nil {
		log.Fatalf("fail to serve: %v", err)
	}
}

func (b *BulkInsert) Init(Host string, DB string, Measurementname string) {
	b.Measurementname = Measurementname
	b.HostString = Host + "/write?db=" + DB
	b.Clusters = make([]string,0,0)
	b.Nodes = make([]string,0,0)
	b.Pods = make([]string,0,0)
	b.Metricnames = make([]string,0,0)
	b.Values = make([]string,0,0)
	b.TimeStamps = make([]int64,0,0)
	b.PostLines = make([]string,0,0)
}

func (b *BulkInsert) GenerateSampleData(Cluster string, Node string, Pod string, MetricName string, Value interface{}) {
	if MetricName=="scrape_error 0" {
		MetricName="scrape_error"
	}
	s := fmt.Sprint(Value)
	t := fmt.Sprint(time.Now().UnixNano())
	tempstring := b.Measurementname + ",Cluster=" + Cluster + ",Node=" + Node + ",Pod=" + Pod + " MetricName=\"" + MetricName + "\",Value=\"" + s + "\" " + t
	//tempstring := b.Measurementname + ",MetricName=" + MetricName +  " Cluster=\"" + Cluster + "\",Node=\"" + Node + "\",Pod=\"" + Pod +  "\",Value=\"" + s + "\" " + t
	b.PostLines = append(b.PostLines, tempstring)
	time.Sleep(time.Microsecond*1)
}

func (b *BulkInsert) InsertData(clustername string) {
	start := time.Now()
	wholestring := fmt.Sprint(strings.Join(b.PostLines,"\n"))
	fmt.Println(len(b.PostLines))
	reqBody := bytes.NewBufferString(wholestring)
	resp, err := http.Post(b.HostString, "text/plain", reqBody)
	if err != nil {
		panic(err)
	}
	end := time.Since(start)
	fmt.Println(clustername,"request done:::",end)
	resp.Body.Close()
	//wg.Done()
	end = time.Since(start)
	fmt.Println(clustername,"done:::",end)
}




func amain(input *res.ClusterInfo, name string) {

	start := time.Now()
	a := &BulkInsert{
		Measurementname: "Metric",
		Clusters:        []string{},
		Nodes:           []string{},
		Pods:            []string{},
		Metricnames:     []string{},
		Values:          []string{},
		TimeStamps:      []int64{},
		PostLines:       []string{},
		HostString:      "",
	}

	cluster := input.Clustername
	node := ""
	pod := "-"
	metricname := ""
	a.Init("http://10.0.3.20:31209", "mydb", "Metric")
	var value interface{}
	for i := 0; i < len(input.NodeList); i++ {
		node = input.NodeList[i].NodeName
		pod = "-"
		metricname = "Zone"
		value = input.NodeList[i].GeoInfo["Zone"]
		a.GenerateSampleData(cluster, node, pod, metricname, value)
		metricname = "Region"
		value = input.NodeList[i].GeoInfo["Region"]
		a.GenerateSampleData(cluster, node, pod, metricname, value)
		metricname = "CPUAllocatable"
		value = input.NodeList[i].NodeAllocatable["CPU"]
		a.GenerateSampleData(cluster, node, pod, metricname, value)
		metricname = "MemoryAllocatable"
		value = input.NodeList[i].NodeAllocatable["Memory"]
		a.GenerateSampleData(cluster, node, pod, metricname, value)
		metricname = "StorageAllocatable"
		value = input.NodeList[i].NodeAllocatable["EphemeralStorage"]
		a.GenerateSampleData(cluster, node, pod, metricname, value)
		for j := 0; j < len(input.NodeList[i].PodList); j++ {
			pod = input.NodeList[i].PodList[j].PodName
			for k := range input.NodeList[i].PodList[j].PodMetrics {
				metricname = k
				value = input.NodeList[i].PodList[j].PodMetrics[k]
				if pod == "" {
					pod = "-"
				}
				a.GenerateSampleData(cluster, node, pod, metricname, value)
			}

		}
	}
	a.InsertData(name)

	end := time.Since(start)
	fmt.Println("generate done:::", end)
}

func InsertGo(a map[string]*BulkInsert) {
	start := time.Now()

	for k := range a{
		wg.Add(1)
		go a[k].InsertData(k)
	}
	end := time.Since(start)
	fmt.Println("for done:::",end)
	wg.Wait()
	end = time.Since(start)
	fmt.Println("wait done:::",end)
}
