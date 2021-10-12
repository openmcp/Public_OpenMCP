package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 1. get dns
// /apis/openmcp.k8s.io/v1alpha1/namespaces/default/openmcpdnsendpoints/
// name : items > metatdata > name
// namespace : items > metadata > namespace (project)
// type : items > spec > type(string)
// selector : items > spec > selector > 여러개 나옴 (key:value 형태로 가져오기)
// port : items > spec > ports[여러개 나옴] > 안에 있는것 모두 나열

func Dns(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	// clusterurl := "https://" + openmcpURL + "/apis/core.kubefed.io/v1beta1/kubefedclusters?clustername=openmcp"

	// go CallAPI(token, clusterurl, ch)
	// clusters := <-ch
	// clusterData := clusters.data

	// resServices := ServicesRes{}
	// clusterNames := []string{}
	// clusterNames = append(clusterNames, "openmcp")

	// //get clusters Information
	// for _, element := range clusterData["items"].([]interface{}) {
	// 	clusterName := GetStringElement(element, []string{"metadata", "name"})
	// 	// element.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
	// 	clusterType := GetStringElement(element, []string{"status", "conditions", "type"})
	// 	if clusterType == "Ready" {
	// 		clusterNames = append(clusterNames, clusterName)
	// 	}

	// }

	// for _, clusterName := range clusterNames {
	resDns := DNSRes{}
	dnsInfo := DNSInfo{}
	// dnsURL := "https://" + openmcpURL + "/apis/openmcp.k8s.io/v1alpha1/openmcpdnsendpoints?clustername=" + clusterName
	// https://192.168.0.152:30000/apis/openmcp.k8s.io/v1alpha1/openmcpdnsendpoints?clustername=openmcp
	dnsURL := "https://" + openmcpURL + "/apis/openmcp.k8s.io/v1alpha1/openmcpdnsendpoints?clustername=" + openmcpClusterName
	go CallAPI(token, dnsURL, ch)
	dnsResult := <-ch
	dnsData := dnsResult.data
	dnsItems := dnsData["items"].([]interface{})

	// get service Information
	for _, element := range dnsItems {
		name := GetStringElement(element, []string{"metadata", "name"})
		namespace := GetStringElement(element, []string{"metadata", "namespace"})

		// Name    string   `json:"name"`
		// Project string   `json:"project"`
		// DnsName string   `json:"dns_name"`
		// IP      []string `json:"ip"`

		dnsName := ""
		ip := ""
		endpoints := GetArrayElement(element, []string{"spec", "endpoints"})
		if endpoints != nil {
			for _, item := range endpoints {
				dnsName = GetStringElement(item, []string{"dnsName"})
				ips := GetArrayElement(item, []string{"targets"})
				if ips != nil {
					for j, ipItem := range ips {
						ipString := fmt.Sprintf("%v", ipItem)
						if j+1 == len(ips) {
							ip = ip + ipString
						} else {
							ip = ip + ipString + ", "
						}
					}
				}

				dnsInfo.Name = name
				dnsInfo.Project = namespace
				dnsInfo.DnsName = dnsName
				dnsInfo.IP = ip

				resDns.DNS = append(resDns.DNS, dnsInfo)

			}
		}

		// port := ""
		// portCheck := GetArrayElement(element, []string{"spec", "ports"})
		// // element.(map[string]interface{})["spec"].(map[string]interface{})["ports"].([]interface{})
		// if portCheck != nil {
		// 	for i, item := range portCheck {
		// 		j := 0
		// 		for key, val := range item.(map[string]interface{}) {
		// 			j++
		// 			value := fmt.Sprintf("%v", val)
		// 			if j == len(item.(map[string]interface{})) {
		// 				port = port + "{ " + key + " : " + value + " }"
		// 			} else {
		// 				port = port + "{ " + key + " : " + value + " },  "
		// 			}
		// 		}
		// 		if i < len(portCheck)-1 {
		// 			port = port + "|"
		// 		}
		// 	}
		// } else {
		// 	port = "-"
		// }
		// createdTime := GetStringElement(element, []string{"metadata", "creationTimestamp"})
		// // element.(map[string]interface{})["metadata"].(map[string]interface{})["creationTimestamp"].(string)

	}
	// }
	json.NewEncoder(w).Encode(resDns.DNS)
}

//get cluster-overview list handler

//get cluster-node list handler

//get cluster-pods list handler
