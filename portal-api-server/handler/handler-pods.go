package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func GetPods(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	clusterURL := "https://" + openmcpURL + "/apis/core.kubefed.io/v1beta1/kubefedclusters?clustername=openmcp"
	go CallAPI(token, clusterURL, ch)
	clusters := <-ch
	clusterData := clusters.data

	resPod := PodRes{}
	clusterNames := []string{}
	clusterNames = append(clusterNames, "openmcp")
	//get clusters Information
	for _, element := range clusterData["items"].([]interface{}) {
		clusterName := GetStringElement(element, []string{"metadata", "name"})
		//  element.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
		clusterType := GetStringElement(element, []string{"status", "conditions", "type"})
		if clusterType == "Ready" {
			clusterNames = append(clusterNames, clusterName)
		}
	}

	for _, clusterName := range clusterNames {
		podURL := "https://" + openmcpURL + "/api/v1/pods?clustername=" + clusterName
		go CallAPI(token, podURL, ch)
		podResult := <-ch
		podData := podResult.data
		podItems := podData["items"].([]interface{})

		// get podUsage counts by nodename groups
		for _, element := range podItems {
			pod := PodInfo{}
			podName := GetStringElement(element, []string{"metadata", "name"})
			// element.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
			project := GetStringElement(element, []string{"metadata", "namespace"})
			// element.(map[string]interface{})["metadata"].(map[string]interface{})["namespace"].(string)
			status := GetStringElement(element, []string{"status", "phase"})
			// element.(map[string]interface{})["status"].(map[string]interface{})["phase"].(string)
			podIP := "-"
			node := "-"
			nodeIP := "-"
			if status == "Running" {
				podIP = GetStringElement(element, []string{"status", "podIP"})
				// element.(map[string]interface{})["status"].(map[string]interface{})["podIP"].(string)
				node = GetStringElement(element, []string{"spec", "nodeName"})
				// element.(map[string]interface{})["spec"].(map[string]interface{})["nodeName"].(string)
				nodeIP = GetStringElement(element, []string{"status", "hostIP"})
				// element.(map[string]interface{})["status"].(map[string]interface{})["hostIP"].(string)
			}

			cpu := "cpu"
			ram := "ram"
			createdTime := GetStringElement(element, []string{"metadata", "creationTimestamp"})
			// element.(map[string]interface{})["metadata"].(map[string]interface{})["creationTimestamp"].(string)

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
	}

	json.NewEncoder(w).Encode(resPod.Pods)
}

func GetPodsInCluster(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	vars := mux.Vars(r)
	clusterName := vars["clusterName"]
	resPod := PodRes{}

	podURL := "https://" + openmcpURL + "/api/v1/pods?clustername=" + clusterName
	go CallAPI(token, podURL, ch)
	podResult := <-ch
	podData := podResult.data
	podItems := podData["items"].([]interface{})

	// get podUsage counts by nodename groups
	for _, element := range podItems {
		pod := PodInfo{}
		podName := GetStringElement(element, []string{"metadata", "name"})
		// element.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
		project := GetStringElement(element, []string{"metadata", "namespace"})
		// element.(map[string]interface{})["metadata"].(map[string]interface{})["namespace"].(string)
		status := GetStringElement(element, []string{"status", "phase"})
		// element.(map[string]interface{})["status"].(map[string]interface{})["phase"].(string)
		podIP := "-"
		node := "-"
		nodeIP := "-"
		if status == "Running" {
			podIP = GetStringElement(element, []string{"status", "podIP"})
			// element.(map[string]interface{})["status"].(map[string]interface{})["podIP"].(string)
			node = GetStringElement(element, []string{"spec", "nodeName"})
			// element.(map[string]interface{})["spec"].(map[string]interface{})["nodeName"].(string)
			nodeIP = GetStringElement(element, []string{"status", "hostIP"})
			// element.(map[string]interface{})["status"].(map[string]interface{})["hostIP"].(string)
		}
		cpu := "cpu"
		ram := "ram"
		createdTime := GetStringElement(element, []string{"metadata", "creationTimestamp"})
		// element.(map[string]interface{})["metadata"].(map[string]interface{})["creationTimestamp"].(string)

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

	json.NewEncoder(w).Encode(resPod.Pods)
}

func GetVPAs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	var allUrls []string

	clusterurl := "https://" + openmcpURL + "/apis/core.kubefed.io/v1beta1/kubefedclusters?clustername=openmcp"
	go CallAPI(token, clusterurl, ch)
	clusters := <-ch
	clusterData := clusters.data
	var clusternames []string
	for _, element := range clusterData["items"].([]interface{}) {
		clusterName := element.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
		clusternames = append(clusternames, clusterName)
	}

	for _, cluster := range clusternames {
		vpaURL := "https://" + openmcpURL + "/apis/autoscaling.k8s.io/v1beta2/verticalpodautoscalers?clustername=" + cluster
		allUrls = append(allUrls, vpaURL)
	}

	for _, arg := range allUrls[0:] {
		go CallAPI(token, arg, ch)
	}

	var results = make(map[string]interface{})
	for range allUrls[0:] {
		result := <-ch
		results[result.url] = result.data
	}
	var VPAResList []VPARes

	for key, result := range results {
		clusterName := string(key[strings.LastIndex(key, "=")+1:])
		items := result.(map[string]interface{})["items"].([]interface{})
		for _, item := range items {
			hpaName := item.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)

			namespace := item.(map[string]interface{})["metadata"].(map[string]interface{})["namespace"].(string)

			reference := item.(map[string]interface{})["spec"].(map[string]interface{})["targetRef"].(map[string]interface{})["kind"].(string) + "/" + item.(map[string]interface{})["spec"].(map[string]interface{})["targetRef"].(map[string]interface{})["name"].(string)

			updateMode := item.(map[string]interface{})["spec"].(map[string]interface{})["updatePolicy"].(map[string]interface{})["updateMode"].(string)

			res := VPARes{hpaName, namespace, clusterName, reference, updateMode}

			VPAResList = append(VPAResList, res)

		}
	}
	json.NewEncoder(w).Encode(VPAResList)
}

func GetHPAs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	var allUrls []string

	clusterurl := "https://" + openmcpURL + "/apis/core.kubefed.io/v1beta1/kubefedclusters?clustername=openmcp"
	go CallAPI(token, clusterurl, ch)
	clusters := <-ch
	clusterData := clusters.data
	var clusternames []string
	for _, element := range clusterData["items"].([]interface{}) {
		clusterName := element.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
		clusternames = append(clusternames, clusterName)
	}
	// add openmcp master cluster
	clusternames = append(clusternames, "openmcp")

	for _, cluster := range clusternames {
		hpaURL := "https://" + openmcpURL + "/apis/autoscaling/v1/horizontalpodautoscalers?clustername=" + cluster
		allUrls = append(allUrls, hpaURL)
	}

	for _, arg := range allUrls[0:] {
		go CallAPI(token, arg, ch)
	}

	var results = make(map[string]interface{})
	for range allUrls[0:] {
		result := <-ch
		results[result.url] = result.data
	}

	var HPAResList []HPARes

	for key, result := range results {
		clusterName := string(key[strings.LastIndex(key, "=")+1:])
		items := result.(map[string]interface{})["items"].([]interface{})
		for _, item := range items {
			hpaName := item.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)

			namespace := item.(map[string]interface{})["metadata"].(map[string]interface{})["namespace"].(string)

			reference := item.(map[string]interface{})["spec"].(map[string]interface{})["scaleTargetRef"].(map[string]interface{})["kind"].(string) + "/" + item.(map[string]interface{})["spec"].(map[string]interface{})["scaleTargetRef"].(map[string]interface{})["name"].(string)

			minRepl := item.(map[string]interface{})["spec"].(map[string]interface{})["minReplicas"].(float64)
			minReplStr := strconv.FormatFloat(minRepl, 'f', -1, 64)

			maxRepl := item.(map[string]interface{})["spec"].(map[string]interface{})["maxReplicas"].(float64)
			maxReplStr := strconv.FormatFloat(maxRepl, 'f', -1, 64)

			currentRepl := item.(map[string]interface{})["status"].(map[string]interface{})["currentReplicas"].(float64)
			currentRepllStr := strconv.FormatFloat(currentRepl, 'f', -1, 64)

			res := HPARes{hpaName, namespace, clusterName, reference, minReplStr, maxReplStr, currentRepllStr}

			HPAResList = append(HPAResList, res)

		}
	}
	json.NewEncoder(w).Encode(HPAResList)
}

func GetPodsInProject(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	vars := mux.Vars(r)
	clusterName := vars["clusterName"]
	projectName := vars["projectName"]
	resPod := PodRes{}

	// http: //192.168.0.152:31635/api/v1/namespaces/kube-system/pods?clustername=cluster2
	podURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/pods?clustername=" + clusterName
	go CallAPI(token, podURL, ch)
	podResult := <-ch
	podData := podResult.data
	podItems := podData["items"].([]interface{})

	// get podUsage counts by nodename groups
	for _, element := range podItems {
		pod := PodInfo{}
		podName := GetStringElement(element, []string{"metadata", "name"})
		// element.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
		project := GetStringElement(element, []string{"metadata", "namespace"})
		// element.(map[string]interface{})["metadata"].(map[string]interface{})["namespace"].(string)
		status := GetStringElement(element, []string{"status", "phase"})
		// element.(map[string]interface{})["status"].(map[string]interface{})["phase"].(string)
		podIP := "-"
		node := "-"
		nodeIP := "-"
		if status == "Running" {
			podIP = GetStringElement(element, []string{"status", "podIP"})
			// element.(map[string]interface{})["status"].(map[string]interface{})["podIP"].(string)
			node = GetStringElement(element, []string{"spec", "nodeName"})
			// element.(map[string]interface{})["spec"].(map[string]interface{})["nodeName"].(string)
			nodeIP = GetStringElement(element, []string{"status", "hostIP"})
			// element.(map[string]interface{})["status"].(map[string]interface{})["hostIP"].(string)
		}
		cpu := "cpu"
		ram := "ram"
		createdTime := GetStringElement(element, []string{"metadata", "creationTimestamp"})
		// element.(map[string]interface{})["metadata"].(map[string]interface{})["creationTimestamp"].(string)

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

	json.NewEncoder(w).Encode(resPod.Pods)
}

func GetPodOverview(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	clusterName := r.URL.Query().Get("cluster")
	podName := vars["podName"]
	// podName := r.URL.Query().Get("pod")
	projectName := r.URL.Query().Get("project")

	// fmt.Println("################", clusterName, podName, projectName)
	if clusterName == "" || podName == "" || projectName == "" {
		errorMSG := jsonErr{500, "failed", "need some params"}
		json.NewEncoder(w).Encode(errorMSG)
	} else {
		ch := make(chan Resultmap)
		token := GetOpenMCPToken()
		// http://192.168.0.152:31635/api/v1/namespaces/{namespace}/pods/{podname}?clustername={clustername}

		podURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/pods/" + podName + "?clustername=" + clusterName
		go CallAPI(token, podURL, ch)

		podResult := <-ch
		podData := podResult.data
		// fmt.Println(podData)
		if podData["spec"] != nil {
			podMetadata := podData["metadata"].(map[string]interface{})
			podSpec := podData["spec"].(map[string]interface{})
			podStatus := podData["status"].(map[string]interface{})
			podPhase := podStatus["phase"].(string)
			totalRestartCount := 0
			var containers []PodOverviewContainer
			if podStatus["containerStatuses"] != nil {
				for _, element := range podStatus["containerStatuses"].([]interface{}) {
					restartCount := int(element.(map[string]interface{})["restartCount"].(float64))
					totalRestartCount = totalRestartCount + restartCount
				}
			}
			for _, element := range podSpec["containers"].([]interface{}) {
				containerNm := element.(map[string]interface{})["name"].(string)
				containerImgae := element.(map[string]interface{})["image"].(string)
				containerPort := "-"
				if element.(map[string]interface{})["ports"] != nil {
					portInt := int(element.(map[string]interface{})["ports"].([]interface{})[0].(map[string]interface{})["containerPort"].(float64))
					containerPort = strconv.Itoa(portInt)
				} else {
					containerPort = "-"
				}
				restartCount := 0
				state := "-"
				if podStatus["containerStatuses"] != nil {
					for _, contStatus := range podStatus["containerStatuses"].([]interface{}) {
						if contStatus.(map[string]interface{})["name"].(string) == containerNm {
							restartCount = int(contStatus.(map[string]interface{})["restartCount"].(float64))
							for k := range contStatus.(map[string]interface{})["state"].(map[string]interface{}) {
								state = k
							}
							break
						}
					}
				}
				container := PodOverviewContainer{containerNm, state, restartCount, containerPort, containerImgae}
				containers = append(containers, container)
			}

			var podConditons []PodOverviewStatus

			for _, element := range podStatus["conditions"].([]interface{}) {
				conType := element.(map[string]interface{})["type"].(string)
				status := element.(map[string]interface{})["status"].(string)
				updateTime := element.(map[string]interface{})["lastTransitionTime"].(string)
				message := "-"
				reason := "-"
				if element.(map[string]interface{})["reason"] != nil {
					reason = element.(map[string]interface{})["reason"].(string)
				}
				if element.(map[string]interface{})["message"] != nil {
					message = element.(map[string]interface{})["message"].(string)
				}

				podConditons = append(podConditons, PodOverviewStatus{conType, status, updateTime, reason, message})
			}

			podHostIP := "-"
			podIP := "-"
			nodeName := "-"
			if podPhase != "Pending" {
				podHostIP = podStatus["hostIP"].(string)
				podIP = podStatus["podIP"].(string)
				nodeName = podSpec["nodeName"].(string)
			}
			podBasicInfo := PodOverviewInfo{
				podMetadata["name"].(string),
				podPhase,
				clusterName,
				podMetadata["namespace"].(string),
				podIP,
				nodeName,
				podHostIP,
				podMetadata["namespace"].(string),
				strconv.Itoa(totalRestartCount),
				podMetadata["creationTimestamp"].(string),
			}

			podMetric := GetInfluxPod10mMetric(clusterName, projectName, podName)

			// http://192.168.0.152:31635/api/v1/namespaces/{namespace}/events?clustername={clustername}
			podEventURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/events?clustername=" + clusterName

			// go CallAPI(token, podEventURL, ch)

			// eventsResult := <-ch
			// eventsData := eventsResult.data
			// eventsItems := eventsData["items"].([]interface{})
			// var events []Event
			// if eventsItems != nil {
			// 	for _, element := range eventsItems {
			// 		// project := element.(map[string]interface{})["metadata"].(map[string]interface{})["namespace"].(string)
			// 		typeNm := element.(map[string]interface{})["type"].(string)
			// 		reason := element.(map[string]interface{})["reason"].(string)
			// 		objectKind := element.(map[string]interface{})["involvedObject"].(map[string]interface{})["kind"].(string)
			// 		objectName := element.(map[string]interface{})["involvedObject"].(map[string]interface{})["name"].(string)
			// 		message := element.(map[string]interface{})["message"].(string)
			// 		time := "-"
			// 		if element.(map[string]interface{})["lastTimestamp"] != nil {
			// 			time = element.(map[string]interface{})["lastTimestamp"].(string)
			// 		}

			// 		if objectKind == "Pod" && objectName == podName {
			// 			events = append(events, Event{"", typeNm, reason, "", message, time})
			// 		}
			// 	}
			// }

			go CallAPI(token, podEventURL, ch)
			eventResult := <-ch
			eventData := eventResult.data
			eventItems := eventData["items"].([]interface{})
			events := []Event{}

			if len(eventItems) > 0 {
				event := Event{}
				for _, element := range eventItems {
					kind := GetStringElement(element, []string{"involvedObject", "kind"})
					objectName := GetStringElement(element, []string{"involvedObject", "name"})
					if kind == "Pod" && objectName == podName {
						event.Typenm = GetStringElement(element, []string{"type"})
						event.Reason = GetStringElement(element, []string{"reason"})
						event.Message = GetStringElement(element, []string{"message"})
						// event.Time = GetStringElement(element, []string{"metadata", "creationTimestamp"})
						event.Time = GetStringElement(element, []string{"lastTimestamp"})
						event.Object = kind
						event.Project = projectName
						events = append(events, event)
					}
				}
			}

			response := PodOverviewRes{podBasicInfo, containers, podConditons, podMetric, events}

			json.NewEncoder(w).Encode(response)
		}
	}
}

func GetPodPhysicalRes(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	clusterName := r.URL.Query().Get("cluster")
	podName := vars["podName"]
	projectName := r.URL.Query().Get("project")

	if clusterName == "" || podName == "" || projectName == "" {
		errorMSG := jsonErr{500, "failed", "need some params"}
		json.NewEncoder(w).Encode(errorMSG)
	} else {
		podMetric := GetInfluxPod10mMetric(clusterName, projectName, podName)
		json.NewEncoder(w).Encode(podMetric)
	}
}
