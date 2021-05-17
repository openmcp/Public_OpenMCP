/*
Copyright 2018 The Multicluster-Controller Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	//"sync"
	//"time"

	"fmt"
	"github.com/oschwald/geoip2-golang"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"openmcp/openmcp/omcplog"
	"strings"
)

var GeoDB, GeoErr = geoip2.Open("/root/GeoLite2-City.mmdb")

func ExtractIP(target string) (string, error) {
	omcplog.V(4).Info("Function Called ExtractIP")
	tmp := strings.Split(target, ":")
	ip, _ := tmp[0], tmp[1]
	omcplog.V(5).Info("IP : " + ip)
	return ip, nil
}

func main() {
	//origin, _ := url.Parse("http://10.0.3.20:8812/")
	origin, _ := url.Parse("http://10.0.3.196")

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host

		fmt.Println(req)

		clientIP, _ := ExtractIP(req.RemoteAddr)
		//clientIP = "119.65.195.180"
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
			zone = "Seoul"
		}

		//국가코드
		fmt.Println("ISO country code(region): ", region)
		fmt.Println("ISO country zone: ", zone)

		req.Header.Add("Client-Zone", "usa")

	}


	proxy := &httputil.ReverseProxy{Director: director}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
		fmt.Println(r.RemoteAddr)
		proxy.ServeHTTP(w, r)

	})

	log.Fatal(http.ListenAndServe(":80", nil))

}
