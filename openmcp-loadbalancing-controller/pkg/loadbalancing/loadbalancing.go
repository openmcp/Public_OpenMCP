package loadbalancing

import (
	//"container/list"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http/httputil"
	"os"
	"sync"

	//"github.com/abh/geoip"
	//"log"
	//"log"
	"math/rand"
	//"net"
	"net/http"
	//"net/http/httputil"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing/clusterregistry"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing/countryregistry"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing/ingressregistry"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing/loadbalancingregistry"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/loadbalancing/serviceregistry"
	"openmcp/openmcp/openmcp-loadbalancing-controller/pkg/protobuf"
	//"github.com/oschwald/geoip2-golang"
	"net/url"
	"strings"
	"time"

	//"k8s.io/klog"
	"openmcp/openmcp/omcplog"
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


//type Queue struct {
//	items []string
//}
//
//func (q *Queue) Set(value string){
//	q.items = append( q.items, value)
//}
//
//func (q *Queue) Get() (string, error){
//	if len(q.items) == 0 {
//		return "", errors.New("No Data")
//	}
//	item, items := q.items[0], q.items[1:]
//	q.items = items
//	return item, nil
//}
//
//func PV(v Queue) *Queue{
//	return &v
//}

//var RR = map[string][]string{}
var RR = map[string]int{}

func extractPath(target *url.URL) (string, error) {
	omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] Function Called ExtractPath")
	path := target.Path
	if len(path) > 1 && path[0] == '/' {
		path = path[1:]
	}
	if path == "favicon.ico" {
		return "", fmt.Errorf("Invalid path")
	}

	omcplog.V(0).Info("Path : " + path)
	return path, nil
}

func extractIP(target string) (string, error) {
	omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] Function Called ExtractIP")
	tmp := strings.Split(target, ":")
	ip, _ := tmp[0], tmp[1]
	omcplog.V(0).Info("IP : " + ip)
        //ip = "202.131.30.11"
        ip = "10.0.6.245"
	return ip, nil
}

//func extractCountry(cip string) string {
//	fmt.Println("*****Extract Country*****")
//	db, err := geoip2.Open("/root/GeoLite2-City.mmdb")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer db.Close()
//	// If you are using strings that may be invalid, check that ip is not nil
//	ip := net.ParseIP("8.8.8.8")
//
//	record, err := db.City(ip)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
//
//	return record.Country.IsoCode
//}

//Traffic 의 ip로부터 국가와 대륙 추출
//func extractGeo(tip string, countryreg countryregistry.Registry) (string, string) {
//	//country := ""
//	//continent := ""
//	country := extractCountry(tip)
//	continent, _ := countryreg.Lookup(country)
//	return country, continent
//}

var SERVER_IP = os.Getenv("GRPC_SERVER")
var SERVER_PORT = os.Getenv("GRPC_PORT")
var grpcClient = protobuf.NewGrpcClient(SERVER_IP, SERVER_PORT)


func Score(clusters []string, tip string, openmcpIP string) map[string]float64 {
	omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] Function Called Score")
	//SERVER_IP := openmcpIP
	//fmt.Println(SERVER_IP2)
	//SERVER_IP :="10.0.3.30"
	//SERVER_PORT := os.Getenv("GRPC_PORT")
	//grpcClient := protobuf.NewGrpcClient(SERVER_IP, SERVER_PORT)

	lbInfo := &protobuf.LBInfo{
		ClusterNameList: clusters,
		ClientIP:        tip,
	}

	omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] Request Geo, Resource Score")
	response, err := grpcClient.SendLBAnalysis(context.TODO(), lbInfo)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(response)
	//fmt.Println(response.ScoreMap)
	omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] Response Geo, Resource Score")
	omcplog.V(0).Info(response.ScoreMap)

	//score := map[string]float64{}
	//for _, cluster := range clusters {
	//	cScore, _ := creg.ResourceScore(cluster)
	//		score[cluster] = cScore
	//	}
	//fmt.Println(score)
	return response.ScoreMap
}

//func geoScore(clusters []string, tcountry, tcontinent string, creg clusterregistry.Registry) map[string]float64 {
//	fmt.Println("*****Geo Score*****")
//
//	score := map[string]float64{}
//	for _, cluster := range clusters {
//		ccountry, _ := creg.Country(cluster)
//		ccontinent, _ := creg.Continent(cluster)
//		if tcountry == ccountry {
//			score[cluster] = 100.0
//		} else if tcontinent == ccontinent {
//			score[cluster] = 50.0
//		} else {
//			score[cluster] = 0.0
//		}
//	}
//	fmt.Println(score)
//	return score
//}

var test_score = map[string]float64 {}

func scoring(clusters []string, tip string, openmcpIP string) string {
	//fmt.Println("*****Scoring*****")
	omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] Function Called scoring")
	if len(clusters) == 1 {
		return clusters[0]
	}
	//gscore := geoScore(clusters, tcountry, tcontinent, creg)
	if len(test_score) == 0 {
		//score := Score(clusters, tip, openmcpIP)
		test_score = Score(clusters, tip, openmcpIP)
	}
	//score := Score(clusters, tip, openmcpIP)
	//cluster := endpointCluster(score)
	cluster := endpointCluster(test_score)
	return cluster
}

//geo score, resource score, hop score를 합쳐서 비율 계산
//난수를 생성하여 비율에 속하는 클러스터를 엔드포인트로 선정
func endpointCluster(score map[string]float64) string {
	omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] Function Called EnpointCluster")
	//geoPolicyWeight := 1.0
	//resourcePolicyWeight := 1.0
	totalScore := 0.0
	endpoint := ""

	sumScore := map[string]float64{}
	for cluster, _ := range score {
		//sumScore[cluster] = (gscore[cluster] * geoPolicyWeight) + (rscore[cluster] * resourcePolicyWeight)
		sumScore[cluster] = score[cluster]
		totalScore = totalScore + sumScore[cluster]
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Float64() * totalScore
	//fmt.Println("***********test*************")
	//fmt.Println(n)
	checkScore := 0.0
	flag := true
	for cluster, _ := range sumScore {
		if flag == true {
			endpoint = cluster
			flag = false
		}
		//checkScore = checkScore + (sumScore[cluster] / totalScore)
		checkScore = checkScore + sumScore[cluster]
		if n <= checkScore {
			endpoint = cluster
			return endpoint
		}
	}
	return endpoint
}


func proxy_lb(host, tip, network, path string, reg loadbalancingregistry.Registry, sreg serviceregistry.Registry, openmcpIP string , creg clusterregistry.Registry) (net.Conn, error) {
	omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] Apply Proxy Server")
	serviceName, err := reg.Lookup(host, path)
	endpoints, err := sreg.Lookup(serviceName)

	if err != nil {
		return nil, err
	}

	for {
		omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] Apply Algorithm : Round Robin")
		//fmt.Println("*****Round Robin*****")
		//fmt.Println(RR)
		var endpoint string
		lock.Lock()
		index := RR[host+path] % len(endpoints)
		endpoint = endpoints[index]
		RR[host+path]++
		defer lock.Unlock()

		omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] Select Endpoint : " + endpoint)

		clusterIP, _ := creg.IngressIP(endpoint)

		conn, err := net.Dial(network, clusterIP + ":80")
		if err != nil {
			fmt.Println(err)
		}
		return conn, nil
	}
	return nil, fmt.Errorf("Error*****")
}


func loadbalancing(host, tip, path string, reg loadbalancingregistry.Registry, creg clusterregistry.Registry, countryreg countryregistry.Registry, sreg serviceregistry.Registry, openmcpIP string) (string, error) {
	omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] Function Called loadbalancing")

	serviceName, err := reg.Lookup(host, path)
	endpoints, err := sreg.Lookup(serviceName)
	omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] Check Service, Endpoint(Cluster)")

	if err != nil {
		return "", err
	}
	//tcountry, tcontinent := extractGeo(tip, countryreg)

	lb := os.Getenv("LB")

	var endpoint string
	if lb == "RR" {
		omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] Apply Algorithm : Round Robin")
		//fmt.Println("*****Round Robin*****")
		//fmt.Println(RR)


		lock.Lock()
		index := RR[host+path] % len(endpoints)
		endpoint = endpoints[index]
		RR[host+path]++
		defer lock.Unlock()

	} else {
		//fmt.Println("*****Geo/Resource*****")
		omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] Apply Algorithm : Geo, Resource Score")
		endpoint = scoring(endpoints, tip, openmcpIP)
	}
	//fmt.Println("*****End Point*****")
	//fmt.Println(endpoint)
	omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] Select Endpoint : " + endpoint)
	return endpoint, err
}







func NewMultipleHostReverseProxy(reg loadbalancingregistry.Registry, creg clusterregistry.Registry, countryreg countryregistry.Registry, sreg serviceregistry.Registry, openmcpIP string) http.HandlerFunc {
	omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] NewMultipleHostReversProxy")

	return func(w http.ResponseWriter, req *http.Request) {
		host := req.Host
		ip, _ := ExtractIP(req.RemoteAddr)
		path, err := ExtractPath(req.URL)
		omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] Extract Host, IP, Path")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		endpoint, _ := loadbalancing(host, ip, path, reg, creg, countryreg, sreg, openmcpIP)

		if path == "/" {
			path = ""
		}
		omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] Exec Redirect (Code : 307)")
		url := "http://" + endpoint + "." + host + "/" + path
		http.Redirect(w, req, url, 307)
	}
}

func NewMultipleHostReverseProxyRR(reg loadbalancingregistry.Registry, creg clusterregistry.Registry, countryreg countryregistry.Registry, sreg serviceregistry.Registry, openmcpIP string) http.HandlerFunc {
	omcplog.V(0).Info("[OpenMCP Loadbalancing Controller] NewMultipleHostReversProxyRR")

	return func(w http.ResponseWriter, req *http.Request) {

		host := req.Host
		ip, _ := ExtractIP(req.RemoteAddr)
		path, err := ExtractPath(req.URL)

		transport := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: func(network, addr string) (net.Conn, error) {
				//addr = strings.Split(addr, ":")[0]
				//tmp := strings.Split(addr, "/")
				//if len(tmp) != 2 {
				//	return nil, ErrInvalidService
				//}
				return proxy_lb(host, ip, network, path, reg , sreg, openmcpIP , creg)
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