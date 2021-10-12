package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetStatefulsets(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	clusterurl := "https://" + openmcpURL + "/apis/core.kubefed.io/v1beta1/kubefedclusters?clustername=openmcp"
	go CallAPI(token, clusterurl, ch)
	clusters := <-ch
	clusterData := clusters.data

	resStatefulset := StatefulsetRes{}
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
		statefulset := StatefulsetInfo{}
		// get node names, cpu(capacity)
		statefulsetURL := "https://" + openmcpURL + "/apis/apps/v1/statefulsets?clustername=" + clusterName
		go CallAPI(token, statefulsetURL, ch)
		statefulsetResult := <-ch
		statefulsetData := statefulsetResult.data
		statefulsetItems := statefulsetData["items"].([]interface{})

		// get deployement Information
		for _, element := range statefulsetItems {
			name := GetStringElement(element, []string{"metadata", "name"})
			namespace := GetStringElement(element, []string{"metadata", "namespace"})

			status := "-"
			availableReplicas := GetInterfaceElement(element, []string{"status", "availableReplicas"})
			readyReplicas := GetInterfaceElement(element, []string{"status", "readyReplicas"})
			replicas := GetFloat64Element(element, []string{"status", "replicas"})

			replS := fmt.Sprintf("%.0f", replicas)

			if readyReplicas != nil {
				readyReplS := fmt.Sprintf("%.0f", readyReplicas)
				status = readyReplS + "/" + replS
			} else if availableReplicas == nil {
				status = "0/" + replS
			} else {
				status = "0/0"
			}

			image := GetStringElement(element, []string{"spec", "template", "spec", "containers", "image"})
			created_time := GetStringElement(element, []string{"metadata", "creationTimestamp"})

			statefulset.Name = name
			statefulset.Status = status
			statefulset.Cluster = clusterName
			statefulset.Project = namespace
			statefulset.Image = image
			statefulset.CreatedTime = created_time
			statefulset.Uid = ""
			statefulset.Labels = make(map[string]interface{})

			resStatefulset.Statefulsets = append(resStatefulset.Statefulsets, statefulset)
		}
	}
	json.NewEncoder(w).Encode(resStatefulset.Statefulsets)
}

//get Statefulset-project list handler
func GetStatefulsetsInProject(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	vars := mux.Vars(r)
	clusterName := vars["clusterName"]
	projectName := vars["projectName"]

	resStatefulset := StatefulsetRes{}
	statefulset := StatefulsetInfo{}
	statefulsetURL := "https://" + openmcpURL + "/apis/apps/v1/namespaces/" + projectName + "/statefulsets?clustername=" + clusterName
	go CallAPI(token, statefulsetURL, ch)
	statefulsetResult := <-ch
	statefulsetData := statefulsetResult.data
	statefulsetItems := statefulsetData["items"].([]interface{})

	// get deployement Information
	for _, element := range statefulsetItems {
		name := GetStringElement(element, []string{"metadata", "name"})
		namespace := GetStringElement(element, []string{"metadata", "namespace"})

		status := "-"
		availableReplicas := GetInterfaceElement(element, []string{"status", "availableReplicas"})
		readyReplicas := GetInterfaceElement(element, []string{"status", "readyReplicas"})

		replicas := GetFloat64Element(element, []string{"status", "replicas"})

		replS := fmt.Sprintf("%.0f", replicas)

		if readyReplicas != nil {
			readyReplS := fmt.Sprintf("%.0f", readyReplicas.(float64))
			status = readyReplS + "/" + replS
		} else if availableReplicas == nil {
			status = "0/" + replS
		} else {
			status = "0/0"
		}

		image := GetStringElement(element, []string{"spec", "template", "spec", "containers", "image"})
		created_time := GetStringElement(element, []string{"metadata", "creationTimestamp"})

		statefulset.Name = name
		statefulset.Status = status
		statefulset.Cluster = clusterName
		statefulset.Project = namespace
		statefulset.Image = image
		statefulset.CreatedTime = created_time
		statefulset.Uid = ""
		statefulset.Labels = make(map[string]interface{})

		resStatefulset.Statefulsets = append(resStatefulset.Statefulsets, statefulset)
	}
	json.NewEncoder(w).Encode(resStatefulset.Statefulsets)
}

//get statefulset-overview
func GetStatefulsetOverview(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	vars := mux.Vars(r)
	clusterName := vars["clusterName"]
	projectName := vars["projectName"]
	statefulsetName := vars["statefulsetName"]

	resStatefulsetOverview := StatefulsetOverview{}
	statefulset := StatefulsetInfo{}
	statefulsetURL := "https://" + openmcpURL + "/apis/apps/v1/namespaces/" + projectName + "/statefulsets/" + statefulsetName + "?clustername=" + clusterName
	go CallAPI(token, statefulsetURL, ch)
	statefulsetResult := <-ch
	statefulsetData := statefulsetResult.data

	// get deployement Information
	name := GetStringElement(statefulsetData, []string{"metadata", "name"})
	namespace := GetStringElement(statefulsetData, []string{"metadata", "namespace"})
	uid := GetStringElement(statefulsetData, []string{"metadata", "uid"})

	status := "-"
	availableReplicas := GetInterfaceElement(statefulsetData, []string{"status", "availableReplicas"})
	readyReplicas := GetInterfaceElement(statefulsetData, []string{"status", "readyReplicas"})
	replicas := GetFloat64Element(statefulsetData, []string{"status", "replicas"})

	replS := fmt.Sprintf("%.0f", replicas)

	if readyReplicas != nil {
		readyReplS := fmt.Sprintf("%.0f", readyReplicas)
		status = readyReplS + "/" + replS
	} else if availableReplicas == nil {
		status = "0/" + replS
	} else {
		status = "0/0"
	}

	image := GetStringElement(statefulsetData, []string{"spec", "template", "spec", "containers", "image"})
	created_time := GetStringElement(statefulsetData, []string{"metadata", "creationTimestamp"})

	labels := make(map[string]interface{})
	labelCheck := GetInterfaceElement(statefulsetData, []string{"metadata", "labels"})
	if labelCheck == nil {
		labels = map[string]interface{}{}
	} else {
		for key, val := range labelCheck.(map[string]interface{}) {
			labels[key] = val
		}
	}

	statefulset.Name = name
	statefulset.Status = status
	statefulset.Cluster = clusterName
	statefulset.Project = namespace
	statefulset.Image = image
	statefulset.CreatedTime = created_time
	statefulset.Uid = uid
	statefulset.Labels = labels

	resStatefulsetOverview.Info = statefulset
	podURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/pods?clustername=" + clusterName
	go CallAPI(token, podURL, ch)
	podResult := <-ch
	podData := podResult.data
	podItems := podData["items"].([]interface{})

	for _, element := range podItems {
		kind := GetStringElement(element, []string{"metadata", "ownerReferences", "kind"})
		name := GetStringElement(element, []string{"metadata", "ownerReferences", "name"})
		if kind == "StatefulSet" && name == statefulsetName {
			//Get pod info
			pod := PodInfo{}
			podName := GetStringElement(element, []string{"metadata", "name"})
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

			cpu := "-"
			ram := "-"
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

			resStatefulsetOverview.Pods = append(resStatefulsetOverview.Pods, pod)
		}
	}

	//ports
	port := PortInfo{}

	containers := GetArrayElement(statefulsetData, []string{"spec", "template", "spec", "containers"})
	for _, element := range containers {
		ports := GetArrayElement(element, []string{"ports"})

		cNames := ""
		cPorts := ""
		cProtocols := ""

		for i, items := range ports {
			if len(ports)-1 == i {
				cNames = cNames + GetStringElement(items, []string{"name"})
				cPorts = cPorts + strconv.FormatFloat(GetFloat64Element(items, []string{"containerPort"}), 'f', -1, 64)
				cProtocols = cProtocols + GetStringElement(items, []string{"protocol"})
			} else {
				cNames = cNames + GetStringElement(items, []string{"name"}) + "|"
				cPorts = cPorts + strconv.FormatFloat(GetFloat64Element(items, []string{"containerPort"}), 'f', -1, 64) + "|"
				cProtocols = cProtocols + GetStringElement(items, []string{"protocol"}) + "|"
			}
		}
		port.Name = cNames
		port.Port = cPorts
		port.Protocol = cProtocols
		if port.Name != "" && port.Port != "" && port.Protocol != "" {
			resStatefulsetOverview.Ports = append(resStatefulsetOverview.Ports, port)
		}
	}
	if len(resStatefulsetOverview.Ports) == 0 {
		resStatefulsetOverview.Ports = []PortInfo{}
	}

	//events
	eventURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/events?clustername=" + clusterName
	go CallAPI(token, eventURL, ch)
	eventResult := <-ch
	eventData := eventResult.data
	eventItems := eventData["items"].([]interface{})
	resStatefulsetOverview.Events = []Event{}

	if len(eventItems) > 0 {
		event := Event{}
		for _, element := range eventItems {
			kind := GetStringElement(element, []string{"involvedObject", "kind"})
			objectName := GetStringElement(element, []string{"involvedObject", "name"})
			if kind == "statefulset" && objectName == statefulsetName {
				event.Typenm = GetStringElement(element, []string{"type"})
				event.Reason = GetStringElement(element, []string{"reason"})
				event.Message = GetStringElement(element, []string{"message"})
				event.Time = GetStringElement(element, []string{"lastTimestamp"})
				event.Object = kind
				event.Project = projectName

				resStatefulsetOverview.Events = append(resStatefulsetOverview.Events, event)
			}
		}
	}

	json.NewEncoder(w).Encode(resStatefulsetOverview)
}
