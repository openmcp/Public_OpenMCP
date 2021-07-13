package reverseProxy

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"
	"strings"

	"github.com/oschwald/geoip2-golang"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

//var GeoDB, GeoErr = geoip2.Open("/root/GeoLite2-City.mmdb")
var GeoDB, GeoErr = geoip2.Open("/root/dbip-city-lite-2021-07.mmdb")

func ExtractIP(target string) (string, error) {
	omcplog.V(4).Info("Function Called ExtractIP")
	tmp := strings.Split(target, ":")
	ip, _ := tmp[0], tmp[1]
	omcplog.V(5).Info("IP : " + ip)
	return ip, nil
}

func ReverseProxy() {
	cm := clusterManager.NewClusterManager()

	svc := &corev1.Service{}
	err := cm.Host_client.Get(context.TODO(), svc, "istio-system", "istio-ingressgateway")
	if err != nil && errors.IsNotFound(err) {
		fmt.Println("Error ! Not Found Service 'istio-ingressgateway' in 'istio-system'")
	} else {

	}
	//origin, _ := url.Parse("http://10.0.3.20:8812/")
	origin, _ := url.Parse("http://" + svc.Status.LoadBalancer.Ingress[0].IP)
	//origin, _ := url.Parse("http://10.0.6.147")

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host

		// fmt.Println(req)
		clientIP, _ := ExtractIP(req.RemoteAddr)
		//clientIP = "14.128.128.5"
		ip := net.ParseIP(clientIP)
		// fmt.Println(req.RemoteAddr, clientIP, ip)

		// record, err := GeoDB.Country(ip)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		record, err := GeoDB.City(ip)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(record)

		//region 이 국가, zone 이 지역
		region := "default"
		zone := "default"

		// IsoCode가 있는경우, 실제 OpenMCP 구성된 클러스터의 Region과 일치하면 Region을 할당
		// 그렇지 않으면 default Region 할당
		findRegion := false
		findZone := false
		if record.Country.IsoCode != "" {

			for _, cluster := range cm.Cluster_list.Items {
				nodeList := &corev1.NodeList{}
				err := cm.Cluster_genClients[cluster.Name].List(context.TODO(), nodeList, "default")
				if err != nil {
					fmt.Println("get NodeList Error")
					continue
				}
				for _, node := range nodeList.Items {

					nodeRegion := node.Labels["topology.kubernetes.io/region"]
					nodeZone := node.Labels["topology.kubernetes.io/zone"]

					if nodeRegion == record.Country.IsoCode {
						findRegion = true
						//region = record.Country.IsoCode
					}

					if len(record.Subdivisions) > 0 && nodeZone == record.Subdivisions[0].Names["en"] {
						findZone = true
						//zone = record.Subdivisions[0].Names["en"]
					}

					if findRegion && findZone {
						break
					}
				}
				if findRegion && findZone {
					break
				}
			}
		}
		if findRegion && findZone {
			region = record.Country.IsoCode
			zone = record.Subdivisions[0].Names["en"]

		}

		fmt.Println("Client IP: ", ip)
		//국가코드
		if record.Country.IsoCode != "" {
			fmt.Print("Client Region: ", record.Country.IsoCode)
		} else {
			fmt.Print("Client Region: ", "Not Found")
		}
		if findRegion && findZone {
			fmt.Print(" (Matched) Use '", region, "'")
		} else {
			fmt.Print(" (Not Matched) Use '", region, "'")
		}
		fmt.Println()

		if len(record.Subdivisions) > 0 {
			fmt.Print("Client Zone: ", record.Subdivisions[0].Names["en"])
		} else {
			fmt.Print("Client Zone: ", "Not Found")
		}
		if findRegion && findZone {
			fmt.Print(" (Matched) Use '", zone, "'")
		} else {
			fmt.Print(" (Not Matched) Use '", zone, "'")
		}
		fmt.Println()

		//req.Header.Add("Client-Zone", "usa")
		//req.Header.Add("Client-Zone", strings.ToLower(zone))
		req.Header.Add("Client-Region", region)
		req.Header.Add("Client-Zone", zone)

	}

	proxy := &httputil.ReverseProxy{Director: director}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("RemoteAddr: ", r.RemoteAddr)
		proxy.ServeHTTP(w, r)

	})

	log.Fatal(http.ListenAndServe(":80", nil))
}
