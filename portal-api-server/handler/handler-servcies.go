package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Services(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	clusterurl := "https://" + openmcpURL + "/apis/core.kubefed.io/v1beta1/kubefedclusters?clustername=openmcp"
	go CallAPI(token, clusterurl, ch)
	clusters := <-ch
	clusterData := clusters.data

	resServices := ServicesRes{}
	clusterNames := []string{}
	clusterNames = append(clusterNames, "openmcp")

	//get clusters Information
	for _, element := range clusterData["items"].([]interface{}) {
		clusterName := GetStringElement(element, []string{"metadata", "name"})
		clusterType := GetStringElement(element, []string{"status", "conditions", "type"})
		if clusterType == "Ready" {
			clusterNames = append(clusterNames, clusterName)
		}

	}

	for _, clusterName := range clusterNames {
		service := ServiceInfo{}
		serviceURL := "https://" + openmcpURL + "/api/v1/services?clustername=" + clusterName
		go CallAPI(token, serviceURL, ch)
		serviceResult := <-ch
		serviceData := serviceResult.data
		serviceItems := serviceData["items"].([]interface{})

		// get service Information
		for _, element := range serviceItems {
			name := GetStringElement(element, []string{"metadata", "name"})
			namespace := GetStringElement(element, []string{"metadata", "namespace"})
			serviceType := GetStringElement(element, []string{"spec", "type"})
			clusterIP := GetStringElement(element, []string{"spec", "clusterIP"})
			externalIP := GetStringElement(element, []string{"status", "loadBalancer", "ingress", "ip"})

			selector := ""
			selectorCheck := GetInterfaceElement(element, []string{"spec", "selector"})
			if selectorCheck != nil {
				i := 0
				for key, val := range selectorCheck.(map[string]interface{}) {
					i++
					value := fmt.Sprintf("%v", val)
					if i == len(selectorCheck.(map[string]interface{})) {
						selector = selector + key + " : " + value
					} else {
						selector = selector + key + " : " + value + "|"
					}
				}
			} else {
				selector = "-"
			}

			port := ""
			portCheck := GetArrayElement(element, []string{"spec", "ports"})
			if portCheck != nil {
				for i, item := range portCheck {
					j := 0
					for key, val := range item.(map[string]interface{}) {
						j++
						value := fmt.Sprintf("%v", val)
						if j == len(item.(map[string]interface{})) {
							port = port + "{ " + key + " : " + value + " }"
						} else {
							port = port + "{ " + key + " : " + value + " },  "
						}
					}
					if i < len(portCheck)-1 {
						port = port + "|"
					}
				}

			} else {
				port = "-"
			}
			createdTime := GetStringElement(element, []string{"metadata", "creationTimestamp"})

			service.Cluster = clusterName
			service.Name = name
			service.Project = namespace
			service.Type = serviceType
			service.Selector = selector
			service.Port = port
			service.CreatedTime = createdTime
			service.ClusterIP = clusterIP
			service.ExternalIP = externalIP

			resServices.Services = append(resServices.Services, service)
		}
	}
	json.NewEncoder(w).Encode(resServices.Services)
}

func GetServicesInProject(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	vars := mux.Vars(r)
	clusterName := vars["clusterName"]
	projectName := vars["projectName"]

	resServices := ServicesRes{}

	service := ServiceInfo{}
	serviceURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/services?clustername=" + clusterName
	go CallAPI(token, serviceURL, ch)
	serviceResult := <-ch
	serviceData := serviceResult.data
	serviceItems := serviceData["items"].([]interface{})

	// get service Information
	for _, element := range serviceItems {
		name := GetStringElement(element, []string{"metadata", "name"})
		namespace := GetStringElement(element, []string{"metadata", "namespace"})
		serviceType := GetStringElement(element, []string{"spec", "type"})
		clusterIP := GetStringElement(element, []string{"spec", "clusterIP"})
		externalIP := GetStringElement(element, []string{"status", "loadBalancer", "ingress", "ip"})

		selector := ""
		selectorCheck := GetInterfaceElement(element, []string{"spec", "selector"})
		if selectorCheck != nil {
			i := 0
			for key, val := range selectorCheck.(map[string]interface{}) {
				i++
				value := fmt.Sprintf("%v", val)
				if i == len(selectorCheck.(map[string]interface{})) {
					selector = selector + key + " : " + value
				} else {
					selector = selector + key + " : " + value + "|"
				}
			}
		} else {
			selector = "-"
		}

		port := ""
		portCheck := GetArrayElement(element, []string{"spec", "ports"})

		if portCheck != nil {
			for i, item := range portCheck {
				j := 0
				for key, val := range item.(map[string]interface{}) {
					j++
					value := fmt.Sprintf("%v", val)
					if j == len(item.(map[string]interface{})) {
						port = port + "{ " + key + " : " + value + " }"
					} else {
						port = port + "{ " + key + " : " + value + " },  "
					}
				}
				if i < len(portCheck)-1 {
					port = port + "|"
				}
			}

		} else {
			port = "-"
		}
		createdTime := GetStringElement(element, []string{"metadata", "creationTimestamp"})

		service.Cluster = clusterName
		service.Name = name
		service.Project = namespace
		service.Type = serviceType
		service.Selector = selector
		service.Port = port
		service.CreatedTime = createdTime
		service.ClusterIP = clusterIP
		service.ExternalIP = externalIP

		resServices.Services = append(resServices.Services, service)
	}
	json.NewEncoder(w).Encode(resServices.Services)
}

func GetServiceOverview(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	vars := mux.Vars(r)
	clusterName := vars["clusterName"]
	projectName := vars["projectName"]
	serviceName := vars["serviceName"]

	resServiceOverview := ServiceOverview{}

	service := ServiceBasicInfo{}
	serviceURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/services/" + serviceName + "?clustername=" + clusterName
	go CallAPI(token, serviceURL, ch)
	serviceResult := <-ch
	serviceData := serviceResult.data

	// get service Information
	name := GetStringElement(serviceData, []string{"metadata", "name"})
	namespace := GetStringElement(serviceData, []string{"metadata", "namespace"})
	serviceType := GetStringElement(serviceData, []string{"spec", "type"})
	clusterIP := GetStringElement(serviceData, []string{"spec", "clusterIP"})
	externalIP := GetStringElement(serviceData, []string{"status", "loadBalancer", "ingress", "ip"})
	createdTime := GetStringElement(serviceData, []string{"metadata", "creationTimestamp"})
	sessionAffinity := GetStringElement(serviceData, []string{"spec", "sessionAffinity"})
	selector := ""
	selectorObject := GetInterfaceElement(serviceData, []string{"spec", "selector"})
	if selectorObject != nil {
		i := 0
		for key, val := range selectorObject.(map[string]interface{}) {
			i++
			value := fmt.Sprintf("%v", val)
			if i == len(selectorObject.(map[string]interface{})) {
				selector = selector + key + " : " + value
			} else {
				selector = selector + key + " : " + value + ", "
			}
		}
	} else {
		selector = "-"
	}

	endPointURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/endpoints/" + serviceName + "?clustername=" + clusterName
	go CallAPI(token, endPointURL, ch)
	endPointResult := <-ch
	endPointData := endPointResult.data

	endPointIPs := []string{}
	endPointPorts := []string{}
	endPointPods := []string{}
	endPointNodes := []string{}
	endPointIPArrObj := GetArrayElement(endPointData, []string{"subsets", "addresses"})
	endPointPortArrObj := GetArrayElement(endPointData, []string{"subsets", "ports"})

	if endPointIPArrObj != nil {
		for _, element := range endPointIPArrObj {
			endPointIPs = append(endPointIPs, element.(map[string]interface{})["ip"].(string))
			targetRef := GetStringElement(element, []string{"targetRef", "kind"})
			if targetRef == "Pod" {
				endPointPods = append(endPointPods, GetStringElement(element, []string{"targetRef", "name"}))
			} else if targetRef == "Node" {
				endPointNodes = append(endPointNodes, GetStringElement(element, []string{"targetRef", "name"}))
			}
		}
	}

	if endPointPortArrObj != nil {
		for _, element := range endPointPortArrObj {
			port := fmt.Sprintf("%.0f", element.(map[string]interface{})["port"].(float64))
			endPointPorts = append(endPointPorts, port)
		}
	}

	endPoint := "" //endpoint에 service명으로 검색해서 결과확인
	for j, ip := range endPointIPs {
		for i, port := range endPointPorts {
			if i == len(endPointPorts)-1 {
				endPoint = endPoint + ip + ":" + port
			} else {
				endPoint = endPoint + ip + ":" + port + ", "
			}
		}

		if j != len(endPointIPs)-1 {
			endPoint = endPoint + ", "
		}
	}

	service.Name = name
	service.Project = namespace
	service.Type = serviceType
	service.Cluster = clusterName
	service.ClusterIP = clusterIP
	service.ExternalIP = externalIP
	service.SessionAffinity = sessionAffinity
	service.Selector = selector
	service.Endpoints = endPoint
	service.CreatedTime = createdTime

	resServiceOverview.Info = service
	resPod := PodRes{}
	if len(endPointPods) > 0 {
		podURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/pods?clustername=" + clusterName
		go CallAPI(token, podURL, ch)
		podResult := <-ch
		podData := podResult.data
		podItems := podData["items"].([]interface{})
		for _, element := range podItems {
			pod := PodInfo{}
			podName := GetStringElement(element, []string{"metadata", "name"})
			isExist := false

			for _, item := range endPointPods {
				if podName == item {
					isExist = true
					break
				}
			}

			if !isExist {
				continue
			}

			project := GetStringElement(element, []string{"metadata", "namespace"})
			status := GetStringElement(element, []string{"status", "phase"})
			podIP := "-"
			node := "-"
			nodeIP := "-"
			if status == "Running" {
				podIP = GetStringElement(element, []string{"status", "podIP"})
				node = GetStringElement(element, []string{"spec", "nodeName"})
				nodeIP = GetStringElement(element, []string{"status", "hostIP"})
			}

			cpu := "cpu"
			ram := "ram"
			createdTime := GetStringElement(element, []string{"metadata", "creationTimestamp"})

			pod.Name = podName
			pod.Status = status
			pod.Cluster = clusterName
			pod.Project = project
			pod.PodIP = podIP
			pod.Node = node
			pod.NodeIP = nodeIP
			pod.Cpu = cpu
			pod.Ram = ram
			pod.CreatedTime = createdTime

			resPod.Pods = append(resPod.Pods, pod)
		}
	} else {
		resPod.Pods = []PodInfo{}
	}

	//events
	eventURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/events?clustername=" + clusterName

	go CallAPI(token, eventURL, ch)
	eventResult := <-ch
	eventData := eventResult.data
	eventItems := eventData["items"].([]interface{})
	resServiceOverview.Events = []Event{}

	if len(eventItems) > 0 {
		event := Event{}
		for _, element := range eventItems {
			kind := GetStringElement(element, []string{"involvedObject", "kind"})
			objectName := GetStringElement(element, []string{"involvedObject", "name"})
			if kind == "Service" && objectName == serviceName {
				event.Typenm = GetStringElement(element, []string{"type"})
				event.Reason = GetStringElement(element, []string{"reason"})
				event.Message = GetStringElement(element, []string{"message"})
				event.Time = GetStringElement(element, []string{"lastTimestamp"})
				event.Object = kind
				event.Project = projectName

				resServiceOverview.Events = append(resServiceOverview.Events, event)
			}
		}
	}

	resServiceOverview.Pods = resPod.Pods

	json.NewEncoder(w).Encode(resServiceOverview)
}
