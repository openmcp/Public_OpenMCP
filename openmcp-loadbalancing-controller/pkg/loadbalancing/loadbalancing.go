package loadbalancing

import (
	"context"
	"errors"
	"fmt"
	"os"

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
)

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

func extractPath(target *url.URL) (string, error) {
	fmt.Println("----Extract Path----")
	path := target.Path
	if len(path) > 1 && path[0] == '/' {
		path = path[1:]
	}
	if path == "favicon.ico" {
		return "", fmt.Errorf("Invalid path")
	}
	fmt.Println("Path : " + path)
	return path, nil
}

func extractIP(target string) (string, error) {
	fmt.Println("----Extract IP----")
	tmp := strings.Split(target, ":")
	ip, _ := tmp[0], tmp[1]
	fmt.Println("IP : " + ip)
        ip = "202.131.30.11"
	fmt.Println(ip)
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
	fmt.Println("*****Resource Score*****")
//	SERVER_IP := openmcpIP
	//fmt.Println(SERVER_IP2)
	//SERVER_IP :="10.0.3.20"
//	SERVER_PORT := os.Getenv("GRPC_PORT")
//	grpcClient := protobuf.NewGrpcClient(SERVER_IP, SERVER_PORT)

	lbInfo := &protobuf.LBInfo{
		ClusterNameList: clusters,
		ClientIP:        tip,
	}

	response, err := grpcClient.SendLBAnalysis(context.TODO(), lbInfo)
	if err != nil {
		fmt.Println(err)
                return nil
	}
	fmt.Println(response)
	fmt.Println(response.ScoreMap)
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

func scoring(clusters []string, tip string, openmcpIP string) string {
	fmt.Println("*****Scoring*****")
	if len(clusters) == 1 {
		return clusters[0]
	}
	//gscore := geoScore(clusters, tcountry, tcontinent, creg)
	score := Score(clusters, tip, openmcpIP)
        var cluster string
        if score == nil {
             cluster = clusters[0]
        } else {
	     cluster = endpointCluster(score)
        }
	return cluster
}

//geo score, resource score, hop score를 합쳐서 비율 계산
//난수를 생성하여 비율에 속하는 클러스터를 엔드포인트로 선정
func endpointCluster(score map[string]float64) string {
	fmt.Println("asfdsadfsadf")
	fmt.Println("*****Endpoint Cluster*****")
	fmt.Println("Endpoint")
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
	fmt.Println("*******test***************")
	fmt.Println(n)
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

func loadbalancing(host, tip, path string, reg loadbalancingregistry.Registry, creg clusterregistry.Registry, countryreg countryregistry.Registry, sreg serviceregistry.Registry, openmcpIP string) (string, error) {
	fmt.Println("*****Loadbalancing*****")

	serviceName, err := reg.Lookup(host, path)
	endpoints, err := sreg.Lookup(serviceName)

	if err != nil {
		return "", err
	}
	//tcountry, tcontinent := extractGeo(tip, countryreg)
	endpoint := scoring(endpoints, tip, openmcpIP)
	fmt.Println("*****End Point*****")
	fmt.Println(endpoint)
	return endpoint, err
}

func NewMultipleHostReverseProxy(reg loadbalancingregistry.Registry, creg clusterregistry.Registry, countryreg countryregistry.Registry, sreg serviceregistry.Registry, openmcpIP string) http.HandlerFunc {
	fmt.Println("*****NewMultipleHostReversProxy*****")

	return func(w http.ResponseWriter, req *http.Request) {
		host := req.Host
		ip, _ := ExtractIP(req.RemoteAddr)
		path, err := ExtractPath(req.URL)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		endpoint, _ := loadbalancing(host, ip, path, reg, creg, countryreg, sreg, openmcpIP)

		if path == "/" {
			path = ""
		}
		url := "http://" + endpoint + "." + host + "/" + path
		http.Redirect(w, req, url, 307)
	}
}
