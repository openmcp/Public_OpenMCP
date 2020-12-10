package loadbalancing

import (
	"context"
	"errors"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing/clusterregistry"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing/countryregistry"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing/geo"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing/ingressregistry"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing/loadbalancingregistry"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing/serviceregistry"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/protobuf"
	"openmcp/openmcp/util/clusterManager"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	//"math"
)

var lock sync.RWMutex

var LoadbalancingRegistry = loadbalancingregistry.DefaultRegistry{}
var ClusterRegistry = clusterregistry.DefaultClusterInfo{}
var CountryRegistry = countryregistry.DefaultCountryInfo{}
var IngressRegistry = ingressregistry.DefaultRegistry{}
var ServiceRegistry = serviceregistry.DefaultRegistry{}

var (
	ErrInvalidService = errors.New("invalid service/version")
)

var ExtractPath = extractPath
var ExtractIP = extractIP

var RR = map[string]int{}

func extractPath(target *url.URL) (string, error) {
	omcplog.V(4).Info("Function Called ExtractPath")
	path := target.Path
	if len(path) > 1 && path[0] == '/' {
		path = path[1:]
	}
	if path == "favicon.ico" {
		return "", fmt.Errorf("Invalid path")
	}

	omcplog.V(5).Info("Path : " + path)
	return path, nil
}

func extractIP(target string) (string, error) {
	omcplog.V(4).Info("Function Called ExtractIP")
	tmp := strings.Split(target, ":")
	ip, _ := tmp[0], tmp[1]
	omcplog.V(5).Info("IP : " + ip)
	return ip, nil
}

var SERVER_IP = os.Getenv("GRPC_SERVER")
var SERVER_PORT = os.Getenv("GRPC_PORT")
var grpcClient = protobuf.NewGrpcClient(SERVER_IP, SERVER_PORT)

var TEST_IP = os.Getenv("TEST_IP")

var ResourceScore = map[string]float64{}

func RequestResourceScore(clusters []string, clientIP string) {
	for {
		omcplog.V(5).Info("Function Called RequestResourceScore")

		lbInfo := &protobuf.LBInfo{
			ClusterNameList: clusters,
			ClientIP:        clientIP,
		}

		response, err := grpcClient.SendLBAnalysis(context.TODO(), lbInfo)
		if err != nil {
			omcplog.V(0).Info(err)
		} else {
			ResourceScore = response.ScoreMap
		}
		omcplog.V(5).Info(ResourceScore)

		time.Sleep(time.Second * time.Duration(Policy["Period"]))
	}
}

func AnalyticRequestResourceScore(clusters []string, clientIP string) map[string]float64 {
	omcplog.V(5).Info("Function Called AnalyticRequestResourceScore")

	lbInfo := &protobuf.LBInfo{
		ClusterNameList: clusters,
		ClientIP:        clientIP,
	}

	response, err := grpcClient.SendLBAnalysis(context.TODO(), lbInfo)
	if err != nil {
		omcplog.V(0).Info(err)
	} else {
		ResourceScore = response.ScoreMap
	}
	omcplog.V(5).Info(ResourceScore)

	return ResourceScore
}

func Score(clusters []string, clientIP string, creg clusterregistry.Registry) map[string]float64 {

	omcplog.V(4).Info("Function Called Score")

	GeoScore := geoScore(clusters, creg, clientIP)
	omcplog.V(0).Info("Geo Score")
	omcplog.V(0).Info(GeoScore)

	omcplog.V(-1).Info("Geo Score")

	for cluster,_ := range GeoScore {
		omcplog.V(-1).Info(cluster," Geo Score : ", GeoScore[cluster])
	}

	var sumScore = map[string]float64{}

	isAnalytic := os.Getenv("isAnalytic")
	if isAnalytic == "yes" {
		AnalyticResourceScore := AnalyticRequestResourceScore(clusters, clientIP)

		for _, cluster := range clusters {
			sumScore[cluster] = GeoScore[cluster] + AnalyticResourceScore[cluster]
		}
	} else {
		for _, cluster := range clusters {
			sumScore[cluster] = GeoScore[cluster] + ResourceScore[cluster]
		}
	}

	omcplog.V(3).Info("Resource Score")
	omcplog.V(3).Info(ResourceScore)


	omcplog.V(-1).Info("Resource Score")

	for cluster,_ := range ResourceScore {
		omcplog.V(-1).Info(cluster," Resource Score : ", ResourceScore[cluster])
	}

	return sumScore
}

var score = map[string]float64{}

func scoring(clusters []string, tip string, creg clusterregistry.Registry) string {

	omcplog.V(4).Info("Function Called scoring")
	if len(clusters) == 1 {
		return clusters[0]
	}

	omcplog.V(4).Info(strconv.Itoa(len(clusters)) + " Clusters Scoring & Compare")

	score = Score(clusters, tip, creg)

	cluster := endpointCluster(score)
	return cluster
}

func endpointCluster(score map[string]float64) string {
	omcplog.V(4).Info("Function Called EnpointCluster")

	totalScore := 0.0
	endpoint := ""

	sumScore := map[string]float64{}
	for cluster, _ := range score {
		sumScore[cluster] = score[cluster]
		totalScore = totalScore + sumScore[cluster]
	}
        var clusterRatio = map[string]float64{}

        for cluster,_ := range  score {
		clusterRatio[cluster] = (sumScore[cluster] / totalScore) * 100
	}


	omcplog.V(-1).Info("Traffic Ratio")

	for cluster,_ := range clusterRatio {
		test := float64(int(clusterRatio[cluster] * 100)) / 100

		//vv := fmt.Sprint("%.2f", clusterRatio[cluster])
		omcplog.V(-1).Info(cluster," Traffic Ratio : ", test, "%")
	}

	//omcplog.V(5).Info("Traffic Ratio")
    //    omcplog.V(5).Info(clusterRatio)
	rand.Seed(time.Now().UnixNano())
	n := rand.Float64() * totalScore
	omcplog.V(4).Info("Random Num : ", n)

	checkScore := 0.0
	flag := true
	for cluster, _ := range sumScore {
		if flag == true {
			endpoint = cluster
			flag = false
		}
		checkScore = checkScore + sumScore[cluster]
		if n <= checkScore {
			endpoint = cluster
			return endpoint
		}
	}
	return endpoint
}

func proxy_lb(host, tip, network, path string, reg loadbalancingregistry.Registry, sreg serviceregistry.Registry, openmcpIP string, creg clusterregistry.Registry) (net.Conn, error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller] Apply Proxy Server")
	serviceName, err := reg.Lookup(host, path)
	endpoints, err := sreg.Lookup(serviceName)

	if err != nil {
		return nil, err
	}

	omcplog.V(5).Info("[OpenMCP Loadbalancing Controller] Apply Algorithm : Round Robin")

	var endpoint string
	rand.Seed(time.Now().UnixNano())
	round_index := rand.Intn(len(endpoints))
	endpoint = endpoints[round_index]

	omcplog.V(3).Info("[OpenMCP Loadbalancing Controller] Select Endpoint : " + endpoint)

	clusterIP, _ := creg.IngressIP(endpoint)

	conn, err := net.Dial(network, clusterIP+":80")
	if err != nil {
		fmt.Println(err)
	}
	return conn, nil

}

var GeoDB, GeoErr = geoip2.Open("/root/GeoLite2-City.mmdb")

func getCountry(clientIP string) string {
	omcplog.V(4).Info("Function Called getCountry")

	if GeoErr != nil {
		log.Fatal(GeoErr)
	}
	clientIP = TEST_IP
	//clientIP = "119.65.195.180"
	ip := net.ParseIP(clientIP)
	record, err := GeoDB.City(ip)
	if err != nil {
		log.Fatal(err)
	}

	//국가코드
	omcplog.V(5).Info("ISO country code: %v\n", record.Country.IsoCode)

//	fmt.Println("국가 코드")
//	fmt.Println(record.Country.IsoCode)
//	fmt.Println("국가/지역1")
//	fmt.Println(record.Subdivisions)

	if len(record.Subdivisions) > 0 {
		fmt.Println(record.Subdivisions[0].Names["en"])
	} else {
		fmt.Println("Not Exist")
	}

//	fmt.Println("국가/지역2")
//	fmt.Println(record.RepresentedCountry)

	return record.Country.IsoCode
}

func getGeo(clientIP string) (string, string) {
	omcplog.V(4).Info("Function Called getGeo")
	if GeoErr != nil {
		log.Fatal(GeoErr)
	}
	clientIP = TEST_IP
	//clientIP = "119.65.195.180"
	ip := net.ParseIP(clientIP)
	record, err := GeoDB.City(ip)
	if err != nil {
		log.Fatal(err)
	}

	//region 이 국가, zone 이 지역
	region := record.Country.IsoCode
	zone := ""

	if len(record.Subdivisions) > 0 {
		zone = record.Subdivisions[0].Names["en"]
	} else {
		zone = "Seoul"
	}
	return region, zone
}


func getContinent(country string) string {
	return Geo.Geo[country]
}

var Policy = map[string]float64{
	"Period" : 10.0,
}

func GetPolicy() {
	for {
		omcplog.V(4).Info("Function Called GetPolicy")
		cm := clusterManager.NewClusterManager()
		openmcpPolicyInstance, target_cluster_policy_err := cm.Crd_client.OpenMCPPolicy("openmcp").Get("loadbalancing-controller-policy", metav1.GetOptions{})

		if target_cluster_policy_err != nil {
			omcplog.V(4).Info(target_cluster_policy_err)
		} else {
			a := openmcpPolicyInstance.Spec.Template.Spec.Policies
			for _, b := range a {
				value, _ := strconv.ParseFloat(b.Value[0], 64)
				Policy[b.Type] = value
			}
			omcplog.V(4).Info("metricsWeight : ", Policy)
		}
		time.Sleep(time.Second * 3)
	}
}

func geoScore(clusters []string, creg clusterregistry.Registry, clientIP string) map[string]float64 {
	omcplog.V(4).Info("Function Called geoScore")

	baseScore := 100.0
	policy := Policy["GeoRate"]

	clientRegion, clientZone := getGeo(clientIP)

	//clientCountry := getCountry(clientIP)
	//clientContinent := getContinent(clientCountry)

	score := map[string]float64{}

	for _, cluster := range clusters{
		clusterRegion, err := creg.Region(cluster)
		if err != nil {
			omcplog.V(0).Info(cluster + " Not set Country")
		}
		clusterZone, err := creg.Zone(cluster)
		if err != nil {
			omcplog.V(0).Info(cluster + " Not set Region")
		}

		if clientZone == clusterZone {
			score[cluster] = baseScore + (baseScore * policy)
		} else if clientRegion == clusterRegion {
			score[cluster] = baseScore
		} else {
			score[cluster] = baseScore - (baseScore * policy)
		}
	}
	//for _, cluster := range clusters {
	//
	//	clustercountry, err := creg.Country(cluster)
	//	if err != nil {
	//		omcplog.V(0).Info(cluster + " Not set Country")
	//	}
	//
	//	clustercontinent, err := creg.Continent(cluster)
	//	if err != nil {
	//		omcplog.V(0).Info(cluster + " Not set Continent")
	//	}
	//
	//	if clientCountry == clustercountry {
	//		score[cluster] = baseScore + (baseScore * policy)
	//	} else if clientContinent == clustercontinent {
	//		score[cluster] = baseScore
	//	} else {
	//		score[cluster] = baseScore - (baseScore * policy)
	//	}
	//}
	omcplog.V(5).Info(score)
	return score
}

func loadbalancing(host, tip, path string, reg loadbalancingregistry.Registry, creg clusterregistry.Registry, countryreg countryregistry.Registry, sreg serviceregistry.Registry, openmcpIP string) (string, error) {
	omcplog.V(4).Info("Function Called loadbalancing")

	omcplog.V(-1).Info("Extract Host from Traffic : ", host)

	serviceName, err := reg.Lookup(host, path)
	endpoints, err := sreg.Lookup(serviceName)
	omcplog.V(-1).Info("Service Discovery, Endpoint(Cluster)")
	omcplog.V(-1).Info(endpoints)

	if err != nil {
		return "", err
	}

	lb := os.Getenv("LB")

	var endpoint string
	if lb == "RR" {
		omcplog.V(5).Info("Apply Algorithm : Round Robin")
		lock.Lock()
		index := RR[host+path] % len(endpoints)
		endpoint = endpoints[index]
		RR[host+path]++
		defer lock.Unlock()

	} else {
		omcplog.V(5).Info("Apply Algorithm : Geo, Resource Score")
		endpoint = scoring(endpoints, tip, creg)
	}
	omcplog.V(-1).Info("Select Endpoint : " + endpoint)
	omcplog.V(3).Info("Select Endpoint : " + endpoint)
	return endpoint, err
}

func NewMultipleHostReverseProxy(reg loadbalancingregistry.Registry, creg clusterregistry.Registry, countryreg countryregistry.Registry, sreg serviceregistry.Registry, openmcpIP string) http.HandlerFunc {
	omcplog.V(4).Info("NewMultipleHostReversProxy")

	return func(w http.ResponseWriter, req *http.Request) {
		host := req.Host
		fmt.Println(req)
		ip, _ := ExtractIP(req.RemoteAddr)
		//test:= getCountry(ip)
		//fmt.Println(test)
		path, err := ExtractPath(req.URL)
		omcplog.V(4).Info("Extract Host, IP, Path")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		endpoint, _ := loadbalancing(host, ip, path, reg, creg, countryreg, sreg, openmcpIP)

		if path == "/" {
			path = ""
		}
		omcplog.V(-1).Info("Exec Redirect (Code : 307)")
		omcplog.V(3).Info("Exec Redirect (Code : 307)")
		url := "http://" + endpoint + "." + host + "/" + path
		http.Redirect(w, req, url, 307)
	}
}

func NewMultipleHostReverseProxyRR(reg loadbalancingregistry.Registry, creg clusterregistry.Registry, countryreg countryregistry.Registry, sreg serviceregistry.Registry, openmcpIP string) http.HandlerFunc {
	omcplog.V(4).Info("NewMultipleHostReversProxyRR")

	return func(w http.ResponseWriter, req *http.Request) {

		host := req.Host
		ip, _ := ExtractIP(req.RemoteAddr)
		path, err := ExtractPath(req.URL)

		transport := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: func(network, addr string) (net.Conn, error) {
				return proxy_lb(host, ip, network, path, reg, sreg, openmcpIP, creg)
			},
			TLSHandshakeTimeout: 10 * time.Second,
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		(&httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = "http"
				req.URL.Host = host + "/" + path
			},
			Transport: transport,
		}).ServeHTTP(w, req)
	}
}
