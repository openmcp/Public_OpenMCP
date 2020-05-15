package lb_proxy

import (
	"errors"
	"fmt"
	//"github.com/abh/geoip"
	//"log"
	"math/rand"
	"net"
	"log"
	"net/http"
	//"net/http/httputil"
	"net/url"
	"strings"
	"time"
	"github.com/HanJaeseung/LoadBalancing/ingressregistry"
	"github.com/HanJaeseung/LoadBalancing/clusterregistry"
	"github.com/HanJaeseung/LoadBalancing/countryregistry"
	"github.com/oschwald/geoip2-golang"
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

//ip로부터 국가 추출
//func extractCountry(tip string) string {
//	fmt.Println(country)
//	return "KR"
//}

func extractCountry(cip string) string{
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
func extractGeo(tip string, countryreg countryregistry.Registry) (string, string){
	//country := ""
	//continent := ""

	country := extractCountry(tip)
	continent,_ := countryreg.Lookup(country)
	return country, continent
}

func resourceScore(clusters []string, creg clusterregistry.Registry) map[string]float64 {
	fmt.Println("*****Resource Score*****")
	score := map[string]float64{}
	for _, cluster := range clusters {
		cScore,_ := creg.ResourceScore(cluster)
		score[cluster] = cScore
	}
	fmt.Println(score)
	return score
}


func geoScore(clusters []string, tcountry, tcontinent string, creg clusterregistry.Registry) map[string]float64 {
	fmt.Println("*****Geo Score*****")

	score := map[string]float64{}
	for _, cluster := range clusters {
		ccountry,_ := creg.Country(cluster)
		ccontinent,_ := creg.Continent(cluster)
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

func HopScore(clusters []string, creg clusterregistry.Registry) map[string]float64 {
	fmt.Println("*****Hop Score*****")
	score := map[string]float64{}
	for _, cluster := range clusters {
		hScore,_ := creg.HopScore(cluster)
		score[cluster] = hScore
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
	hscore := HopScore(clusters, creg)
	cluster := endpointCluster(gscore, rscore, hscore)
	////endpoint,_ := creg.IngressIP(cluster)
	//endpoint = endpoint + ":80"
	//fmt.Println(endpoint)
	return cluster
}


//스코어링 진행
//func scoring(clusters []string, tcountry, tcontinent string, creg clusterregistry.Registry) string {
//	fmt.Println("*****Scoring*****")
//	if len(clusters) == 1 {
//		endpoint,_ := creg.IngressIP(clusters[0])
//		endpoint = endpoint + ":80"
//		return endpoint
//	}
//	gscore := geoScore(clusters, tcountry, tcontinent, creg)
//	rscore := resourceScore(clusters, creg)
//	hscore := HopScore(clusters, creg)
//	cluster := endpointCluster(gscore, rscore, hscore)
//	endpoint,_ := creg.IngressIP(cluster)
//	endpoint = endpoint + ":80"
//	fmt.Println(endpoint)
//	return endpoint
//}

//geo score, resource score, hop score를 합쳐서 비율 계산
//난수를 생성하여 비율에 속하는 클러스터를 엔드포인트로 선정
func endpointCluster(gscore map[string]float64, rscore map[string]float64, hscore map[string]float64) string {
	fmt.Println("*****Endpoint Cluster*****")
	geoPolicyWeight := 1.0
	resourcePolicyWeight := 1.0
	hopPolicyWeight := 1.0
	totalScore := 0.0
	endpoint := ""

	sumScore := map[string]float64{}
	for cluster,_ := range gscore {
		sumScore[cluster] = (gscore[cluster] * geoPolicyWeight) + (rscore[cluster] * resourcePolicyWeight)  + (hscore[cluster] * hopPolicyWeight)
		totalScore = totalScore + sumScore[cluster]
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Float64()
	checkScore := 0.0
	for cluster,_ := range sumScore {
		checkScore = checkScore + (sumScore[cluster] / totalScore)
		if n <= checkScore {
			endpoint = cluster
			return endpoint
		}
	}
	fmt.Println("End Point Cluster's Ingress IP")
	fmt.Println(endpoint)
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
		fmt.Println("**********************")
		fmt.Println(conn)
		return conn, nil
	}
	return nil, fmt.Errorf("No endpoint available for %s", servicePath)
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
		endpoints, err := reg.Lookup(host, path)
		if err != nil {
			fmt.Println(err)
		}
		//Traffic IP로부터 국가,대륙 추출
		tcountry, tcontinent := extractGeo(ip, countryreg)
		//스코어링 진행하여 엔드포인트 결정
		endpoint := scoring(endpoints, tcountry, tcontinent, creg)
		fmt.Println(endpoint)
		url := "http://" + endpoint + "." + host + "/" + path
		http.Redirect(w, req, url, 307)
		//http.Redirect(w,req, "http://cluster1.lb.org.kchtest.org/test", 307)
	}
}

//redirect 방식
//func NewMultipleHostReverseProxy(reg ingressregistry.Registry, creg clusterregistry.Registry, countryreg countryregistry.Registry) http.HandlerFunc {
//	fmt.Println("*****NewMultipleHostReversProxy*****")
//
//	return func(w http.ResponseWriter, req *http.Request) {
//		host := req.Host
//		ip, _ := ExtractIP(req.RemoteAddr)
//		path, err := ExtractPath(req.URL)
//
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		endpoints, err := reg.Lookup(host, path)
//		if err != nil {
//			fmt.Println(err)
//		}
//		//Traffic IP로부터 국가,대륙 추출
//		tcountry, tcontinent := extractGeo(ip, countryreg)
//		//스코어링 진행하여 엔드포인트 결정
//		endpoint := scoring(endpoints, tcountry, tcontinent, creg)
//		fmt.Println(endpoint)
//
//		http.Redirect(w,req, "http://cluster1.lb.org.kchtest.org/test", 307)
//	}
//}



//func NewMultipleHostReverseProxy(reg ingressregistry.Registry, creg clusterregistry.Registry, countryreg countryregistry.Registry) http.HandlerFunc {
//	fmt.Println("*****NewMultipleHostReversProxy*****")
//
//	return func(w http.ResponseWriter, req *http.Request) {
//		host := req.Host
//		ip, _ := ExtractIP(req.RemoteAddr)
//		path, err := ExtractPath(req.URL)
//
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		(&httputil.ReverseProxy{
//			Director: func(req *http.Request) {
//				req.URL.Scheme = "http"
//				req.URL.Host = path
//			},
//			Transport: &http.Transport{
//				Proxy: http.ProxyFromEnvironment,
//				Dial: func(network, addr string) (net.Conn, error) {
//					addr = strings.Split(addr, ":")[0]
//					return LoadBalance(host, ip,  network, addr, reg, creg, countryreg)
//				},
//				TLSHandshakeTimeout: 10 * time.Second,
//			},
//		}).ServeHTTP(w, req)
//	}
//}


//backup code
//*********************************************************************************************
//*********************************************************************************************

//
//func backup_endpointCluster(dscore map[string]float64, rscore map[string]float64) string {
//	fmt.Println("----Select Cluster----")
//	distancePolicyWeight := 1.0
//	resourcePolicyWeight := 1.0
//	//maxScore := 0.0
//	//maxCluster := ""
//	totalScore := 0.0
//	endpoint := ""
//
//	sumScore := map[string]float64{}
//	for cluster,_ := range dscore {
//		sumScore[cluster] = (dscore[cluster] * distancePolicyWeight) + (rscore[cluster] * resourcePolicyWeight)
//		totalScore = totalScore + sumScore[cluster]
//	}
//
//	rand.Seed(time.Now().UnixNano())
//	n := rand.Float64()
//	checkScore := 0.0
//	for cluster,_ := range sumScore {
//		checkScore = checkScore + (sumScore[cluster] / totalScore)
//		if n <= checkScore {
//			endpoint = cluster
//			return endpoint
//		}
//	}
//
//
//	////print 용
//	//
//	//fmt.Println("")
//	//fmt.Println("Geo Score")
//	//fmt.Println(dscore)
//	//fmt.Println("")
//	//fmt.Println("Resource Score")
//	//fmt.Println(rscore)
//	//fmt.Println("")
//	//fmt.Println("Traffic Ratio")
//	//trafficRatio := map[string]float64{}
//	//trafficRatio["cluster1"] = 6
//	//trafficRatio["cluster2"] = 3
//	//trafficRatio["cluster3"] = 1
//	////for cluster,_ := range sumScore {
//	////	trafficRatio[cluster] = (sumScore[cluster] / totalScore) * 100
//	////}
//	//fmt.Println(trafficRatio)
//
//	//testrand := rand.Intn(100)
//	//if testrand >= 0 && testrand < 60 {
//	//	endpoint = "cluster1"
//	//}else if testrand >= 61 && testrand <90 {
//	//	endpoint = "cluster2"
//	//} else {
//	//	endpoint = "cluster3"
//	//}
//	//
//
//
//	//for cluster,_ := range dscore {
//	//	sumScore := (dscore[cluster] * distancePolicyWeight) + (rscore[cluster] * resourcePolicyWeight)
//	//	if maxScore <= sumScore {
//	//		maxScore = sumScore
//	//		maxCluster = cluster
//	//	}
//	//}
//
//	fmt.Println("End Point Cluster")
//	fmt.Println(endpoint)
//	return endpoint
//}


//func backup_scoring(clusters []string, tcountry string, tlat, tlon float64, creg clusterregistry.Registry) string {
//	fmt.Println("----Scoring----")
//
//	if len(clusters) == 1 {
//		endpoint,_ := creg.IngressIP(clusters[0])
//		endpoint = endpoint + ":80"
//		return endpoint
//	}
//	dscore := distanceScore(clusters, tcountry, tlat, tlon, creg)
//	rscore := resourceScore(clusters, creg)
//	cluster := backup_endpointCluster(dscore, rscore)
//	endpoint,_ := creg.IngressIP(cluster)
//	endpoint = endpoint + ":80"
//	return endpoint
//}


//func extractGeo(cip string) (string, float64, float64){
//	fmt.Println("----Extract Geo----")
//	db, err := geoip2.Open("GeoLite2-City.mmdb")
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
//	fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["pt-BR"])
//	if len(record.Subdivisions) > 0 {
//		fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["en"])
//	}
//	fmt.Printf("Russian country name: %v\n", record.Country.Names["ru"])
//	//fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
//	//fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
//
//	fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
//	return record.Country.IsoCode, record.Location.Latitude, record.Location.Longitude
//}


//func distanceScore(clusters []string, tcountry string, tlat, tlon float64, creg clusterregistry.Registry) map[string]float64 {
//	fmt.Println("----Distance Score----")
//	score := map[string]float64{}
//
//	var policyDistance = []float64{10.0, 100.0, 1000.0, 1000000}
//
//	for _,cluster := range clusters {
//		//ccountry,_ := creg.Country(cluster)
//		//ccontinent,_ := creg.Continent(cluster)
//		clat,_ := creg.Latitude(cluster)
//		clon,_ := creg.Longitude(cluster)
//		distance := calcDistance(tlat, tlon, clat, clon)
//
//		score[cluster] = 100.0
//		for i := range policyDistance {
//
//			if distance >= policyDistance[i] {
//				score[cluster] = score[cluster] - (100.0 / float64(len(policyDistance)))
//			}
//		}
//	}
//	return score
//}

