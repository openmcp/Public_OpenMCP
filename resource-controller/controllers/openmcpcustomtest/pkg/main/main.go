package main

import (
	//	"bytes"
	//	"crypto/tls"
	"fmt"
	res "github.com/hth0919/resourcecollector"
	//	"net/http"
	//	"os"
	//	"strconv"
)

func main() {
	clientCluster := res.ClusterInfo{}
	clientCluster.NewClusterClient("")
	clientCluster.NodeListInit()
	/*
		apiserver := clientCluster.Host
		baselink := "/api/v1/namespaces/custom-metrics/services/custom-metrics-apiserver:http/proxy/"
		basepath := "write-metrics"
		resourceKind := "pods"
	*/
	resourceNamespace := res.PodInfo{}.PodNamespace
	resourceName := res.PodInfo{}.PodName
	podList := res.PodInfo{}.PodMetrics

	fmt.Println(resourceNamespace)
	fmt.Println(resourceName)
	fmt.Println(podList)
	/*
		var resourceMetricName string
		var resourceMetricValue float64

		for key, value := range podList {
			resourceMetricName = key
			resourceMetricValue = value

			valueString := strconv.FormatFloat(resourceMetricValue, 'e', 4, 64)

			url := "" + apiserver + baselink + basepath + "/namespaces/" + resourceNamespace + "/" + resourceKind + "/" + resourceName + "/" + resourceMetricName
			buff := bytes.NewBufferString(valueString)

			fmt.Println(url)

			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{Transport: tr}

			req, err := http.NewRequest("POST", os.ExpandEnv(url), buff)
			if err != nil {
				// handle err
			}

			adminToken := clientCluster.AdminToken

			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", os.ExpandEnv("Bearer"+adminToken))

			resp, err := client.Do(req)
			if err != nil {
				// handle err
			}
			defer resp.Body.Close()
		}
	*/
}
