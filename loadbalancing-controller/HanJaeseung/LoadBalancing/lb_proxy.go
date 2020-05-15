package lb_proxy

import (
	"errors"
	"fmt"
	//"github.com/abh/geoip"
	//"log"
	"log"
	"math/rand"
	"net"
	"net/http"
	//"net/http/httputil"
	"github.com/HanJaeseung/LoadBalancing/clusterregistry"
	"github.com/HanJaeseung/LoadBalancing/countryregistry"
	"github.com/HanJaeseung/LoadBalancing/ingressregistry"
	"github.com/oschwald/geoip2-golang"
	"net/url"
	"strings"
	"time"
	//"github.com/umahmood/haversine"
)

var IngressRegistry = ingressregistry.DefaultRegistry{}
var ClusterRegistry = clusterregistry.DefaultClusterInfo{}
var CountryRegistry = countryregistry.DefaultCountryInfo{}

var (
	ErrInvalidService = errors.New("invalid service/version")
)

var ExtractPath = extractPath
var LoadBalance = loadBalance
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
	return ip, nil
}

func extractCountry(cip string) string {
	fmt.Println("*****Extract Country*****")
	db, err := geoip2.Open("/root/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP("8.8.8.8")

	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)

	return record.Country.IsoCode
}

//Traffic 의 ip로부터 국가와 대륙 추출
func extractGeo(tip string, countryreg countryregistry.Registry) (string, string) {
	//country := ""
	//continent := ""
	country := extractCountry(tip)
	continent, _ := countryreg.Lookup(country)
	return country, continent
}

func resourceScore(clusters []string, creg clusterregistry.Registry) map[string]float64 {
	fmt.Println("*****Resource Score*****")
	score := map[string]float64{}
	for _, cluster := range clusters {
		cScore, _ := creg.ResourceScore(cluster)
		score[cluster] = cScore
	}
	fmt.Println(score)
	return score
}

func geoScore(clusters []string, tcountry, tcontinent string, creg clusterregistry.Registry) map[string]float64 {
	fmt.Println("*****Geo Score*****")

	score := map[string]float64{}
	for _, cluster := range clusters {
		ccountry, _ := creg.Country(cluster)
		ccontinent, _ := creg.Continent(cluster)
		if tcountry == ccountry {
			score[cluster] = 100.0
		} else if tcontinent == ccontinent {
			score[cluster] = 50.0
		} else {
			score[cluster] = 0.0
		}
	}
	fmt.Println(score)
	return score
}

func scoring(clusters []string, tcountry, tcontinent string, creg clusterregistry.Registry) string {
	fmt.Println("*****Scoring*****")
	if len(clusters) == 1 {
		return clusters[0]
	}
	gscore := geoScore(clusters, tcountry, tcontinent, creg)
	rscore := resourceScore(clusters, creg)
	cluster := endpointCluster(gscore, rscore)
	return cluster
}

//geo score, resource score, hop score를 합쳐서 비율 계산
//난수를 생성하여 비율에 속하는 클러스터를 엔드포인트로 선정
func endpointCluster(gscore map[string]float64, rscore map[string]float64) string {
	fmt.Println("*****Endpoint Cluster*****")
	geoPolicyWeight := 1.0
	resourcePolicyWeight := 1.0
	totalScore := 0.0
	endpoint := ""

	sumScore := map[string]float64{}
	for cluster, _ := range gscore {
		sumScore[cluster] = (gscore[cluster] * geoPolicyWeight) + (rscore[cluster] * resourcePolicyWeight)
		totalScore = totalScore + sumScore[cluster]
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Float64()
	checkScore := 0.0
	for cluster, _ := range sumScore {
		checkScore = checkScore + (sumScore[cluster] / totalScore)
		if n <= checkScore {
			endpoint = cluster
			return endpoint
		}
	}
	return endpoint
}

func loadBalance(host, tip, network, servicePath string, reg ingressregistry.Registry, creg clusterregistry.Registry, countryreg countryregistry.Registry) (net.Conn, error) {
	fmt.Println("*****LoadBalance*****")

	endpoints, err := reg.Lookup(host, servicePath)
	if err != nil {
		return nil, err
	}
	for {
		//Traffic IP로부터 국가,대륙 추출
		tcountry, tcontinent := extractGeo(tip, countryreg)
		//스코어링 진행하여 엔드포인트 결정
		endpoint := scoring(endpoints, tcountry, tcontinent, creg)
		fmt.Println(endpoint)
		conn, err := net.Dial(network, endpoint)

		if err != nil {
			reg.Failure(host, servicePath, endpoint, err)
			continue
		}
		return conn, nil
	}
	return nil, fmt.Errorf("No endpoint available for %s", servicePath)
}

func loadbalancing(host, tip, servicePath string, reg ingressregistry.Registry, creg clusterregistry.Registry, countryreg countryregistry.Registry) (string, error) {
	fmt.Println("*****Loadbalancing*****")

	endpoints, err := reg.Lookup(host, servicePath)
	if err != nil {
		return "", err
	}
	tcountry, tcontinent := extractGeo(tip, countryreg)
	endpoint := scoring(endpoints, tcountry, tcontinent, creg)
	fmt.Println("*****End Point*****")
	fmt.Println(endpoint)
	return endpoint, err
}

func NewMultipleHostReverseProxy(reg ingressregistry.Registry, creg clusterregistry.Registry, countryreg countryregistry.Registry) http.HandlerFunc {
	fmt.Println("*****NewMultipleHostReversProxy*****")

	return func(w http.ResponseWriter, req *http.Request) {
		host := req.Host
		ip, _ := ExtractIP(req.RemoteAddr)
		path, err := ExtractPath(req.URL)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		endpoint, _ := loadbalancing(host, ip, path, reg, creg, countryreg)
		url := "http://" + endpoint + "." + host + "/" + path
		http.Redirect(w, req, url, 307)
	}
}
