package main

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"
	"strings"
)

func ExtractIP(target string) (string, error) {
	omcplog.V(4).Info("Function Called ExtractIP")
	tmp := strings.Split(target, ":")
	ip, _ := tmp[0], tmp[1]
	omcplog.V(5).Info("IP : " + ip)
	return ip, nil
}

func reverseProxy() {
	cm := clusterManager.NewClusterManager()
	svc := &corev1.Service{}
	err := cm.Host_client.Get(context.TODO(), svc, "istio-system", "istio-ingressgateway")
	if err != nil && errors.IsNotFound(err) {
		fmt.Println("Error ! Not Found Service 'istio-ingressgateway' in 'istio-system'")
	} else {

	}
	//origin, _ := url.Parse("http://10.0.3.20:8812/")
	origin, _ := url.Parse("http://"+svc.Status.LoadBalancer.Ingress[0].IP)
	//origin, _ := url.Parse("http://10.0.3.195")

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host


		clientIP, _ := ExtractIP(req.RemoteAddr)
		//clientIP = "71.67.12.248"
		ip := net.ParseIP(clientIP)
		fmt.Println(req.RemoteAddr, clientIP, ip)
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
			//zone = "Seoul"
		}

		//국가코드
		fmt.Println("ISO country code(region): ", region)
		fmt.Println("ISO country zone: ", zone)

		//req.Header.Add("Client-Zone", "usa")
		//req.Header.Add("Client-Zone", strings.ToLower(zone))
		req.Header.Add("Client-Zone", strings.ToLower(zone))

	}


	proxy := &httputil.ReverseProxy{Director: director}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("RemoteAddr: ", r.RemoteAddr)
		proxy.ServeHTTP(w, r)

	})


	log.Fatal(http.ListenAndServe(":80", nil))
}