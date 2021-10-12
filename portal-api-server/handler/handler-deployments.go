package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetDeployments(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	// vars := mux.Vars(r)
	// clusterName := vars["clusterName"]
	// projectName := vars["projectName"]

	// fmt.Println(clustrName, projectName)
	clusterurl := "https://" + openmcpURL + "/apis/core.kubefed.io/v1beta1/kubefedclusters?clustername=openmcp"
	go CallAPI(token, clusterurl, ch)
	clusters := <-ch
	clusterData := clusters.data

	resDeployment := DeploymentRes{}
	clusterNames := []string{}
	clusterNames = append(clusterNames, "openmcp")

	//get clusters Information
	for _, element := range clusterData["items"].([]interface{}) {
		clusterName := GetStringElement(element, []string{"metadata", "name"})
		// element.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
		clusterType := GetStringElement(element, []string{"status", "conditions", "type"})
		if clusterType == "Ready" {
			clusterNames = append(clusterNames, clusterName)
		}
	}

	for _, clusterName := range clusterNames {
		deployment := DeploymentInfo{}
		// get node names, cpu(capacity)
		deploymentURL := "https://" + openmcpURL + "/apis/apps/v1/deployments?clustername=" + clusterName
		go CallAPI(token, deploymentURL, ch)
		deploymentResult := <-ch
		// fmt.Println(deploymentResult)
		deploymentData := deploymentResult.data
		deploymentItems := deploymentData["items"].([]interface{})

		// get deployement Information
		for _, element := range deploymentItems {
			name := GetStringElement(element, []string{"metadata", "name"})
			// element.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
			namespace := GetStringElement(element, []string{"metadata", "namespace"})
			// element.(map[string]interface{})["metadata"].(map[string]interface{})["namespace"].(string)

			status := "-"
			availableReplicas := GetInterfaceElement(element, []string{"status", "availableReplicas"})
			// element.(map[string]interface{})["status"].(map[string]interface{})["availableReplicas"]
			readyReplicas := GetInterfaceElement(element, []string{"status", "readyReplicas"})
			// element.(map[string]interface{})["status"].(map[string]interface{})["readyReplicas"]
			replicas := GetFloat64Element(element, []string{"status", "replicas"})
			// element.(map[string]interface{})["status"].(map[string]interface{})["replicas"].(float64)

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
			// element.(map[string]interface{})["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})[0].(map[string]interface{})["image"].(string)
			created_time := GetStringElement(element, []string{"metadata", "creationTimestamp"})
			// element.(map[string]interface{})["metadata"].(map[string]interface{})["creationTimestamp"].(string)

			deployment.Name = name
			deployment.Status = status
			deployment.Cluster = clusterName
			deployment.Project = namespace
			deployment.Image = image
			deployment.CreatedTime = created_time
			deployment.Uid = ""
			deployment.Labels = make(map[string]interface{})

			resDeployment.Deployments = append(resDeployment.Deployments, deployment)
		}
	}
	json.NewEncoder(w).Encode(resDeployment.Deployments)
}

//get deployment-project list handler
func GetDeploymentsInProject(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	// fmt.Println("GetDeploymentsInProject")

	vars := mux.Vars(r)
	clusterName := vars["clusterName"]
	projectName := vars["projectName"]

	resDeployment := DeploymentRes{}
	deployment := DeploymentInfo{}
	// get node names, cpu(capacity)
	// http: //192.168.0.152:31635/apis/apps/v1/namespaces/kube-system/deployments?clustername=cluster1
	deploymentURL := "https://" + openmcpURL + "/apis/apps/v1/namespaces/" + projectName + "/deployments?clustername=" + clusterName
	go CallAPI(token, deploymentURL, ch)
	deploymentResult := <-ch
	// fmt.Println(deploymentResult)
	deploymentData := deploymentResult.data
	deploymentItems := deploymentData["items"].([]interface{})

	// get deployement Information
	for _, element := range deploymentItems {
		name := GetStringElement(element, []string{"metadata", "name"})
		// element.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
		namespace := GetStringElement(element, []string{"metadata", "namespace"})
		// element.(map[string]interface{})["metadata"].(map[string]interface{})["namespace"].(string)

		status := "-"
		availableReplicas := GetInterfaceElement(element, []string{"status", "availableReplicas"})
		// element.(map[string]interface{})["status"].(map[string]interface{})["availableReplicas"]
		readyReplicas := GetInterfaceElement(element, []string{"status", "readyReplicas"})

		// element.(map[string]interface{})["status"].(map[string]interface{})["readyReplicas"]

		replicas := GetFloat64Element(element, []string{"status", "replicas"})
		// element.(map[string]interface{})["status"].(map[string]interface{})["replicas"].(float64)

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
		// element.(map[string]interface{})["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})[0].(map[string]interface{})["image"].(string)
		created_time := GetStringElement(element, []string{"metadata", "creationTimestamp"})
		// element.(map[string]interface{})["metadata"].(map[string]interface{})["creationTimestamp"].(string)

		deployment.Name = name
		deployment.Status = status
		deployment.Cluster = clusterName
		deployment.Project = namespace
		deployment.Image = image
		deployment.CreatedTime = created_time
		deployment.Uid = ""
		deployment.Labels = make(map[string]interface{})

		resDeployment.Deployments = append(resDeployment.Deployments, deployment)
	}
	json.NewEncoder(w).Encode(resDeployment.Deployments)
}

//get deployment-overview
func GetDeploymentOverview(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	// fmt.Println("GetDeploymentsInProject")

	vars := mux.Vars(r)
	clusterName := vars["clusterName"]
	projectName := vars["projectName"]
	deploymentName := vars["deploymentName"]

	resDeploymentOverview := DeploymentOverview{}
	deployment := DeploymentInfo{}
	// get node names, cpu(capacity)
	// http: //192.168.0.152:31635/apis/apps/v1/namespaces/kube-system/deployments?clustername=cluster1
	deploymentURL := "https://" + openmcpURL + "/apis/apps/v1/namespaces/" + projectName + "/deployments/" + deploymentName + "?clustername=" + clusterName
	go CallAPI(token, deploymentURL, ch)
	deploymentResult := <-ch
	// fmt.Println(deploymentResult)
	deploymentData := deploymentResult.data

	// get deployement Information
	name := GetStringElement(deploymentData, []string{"metadata", "name"})
	namespace := GetStringElement(deploymentData, []string{"metadata", "namespace"})
	uid := GetStringElement(deploymentData, []string{"metadata", "uid"})

	status := "-"
	availableReplicas := GetInterfaceElement(deploymentData, []string{"status", "availableReplicas"})
	readyReplicas := GetInterfaceElement(deploymentData, []string{"status", "readyReplicas"})
	replicas := GetFloat64Element(deploymentData, []string{"status", "replicas"})

	replS := fmt.Sprintf("%.0f", replicas)

	if readyReplicas != nil {
		readyReplS := fmt.Sprintf("%.0f", readyReplicas)
		status = readyReplS + "/" + replS
	} else if availableReplicas == nil {
		status = "0/" + replS
	} else {
		status = "0/0"
	}

	image := GetStringElement(deploymentData, []string{"spec", "template", "spec", "containers", "image"})
	created_time := GetStringElement(deploymentData, []string{"metadata", "creationTimestamp"})

	labels := make(map[string]interface{})
	labelCheck := GetInterfaceElement(deploymentData, []string{"metadata", "labels"})
	if labelCheck == nil {
		labels = map[string]interface{}{}
	} else {
		for key, val := range labelCheck.(map[string]interface{}) {
			labels[key] = val
		}
	}

	deployment.Name = name
	deployment.Status = status
	deployment.Cluster = clusterName
	deployment.Project = namespace
	deployment.Image = image
	deployment.CreatedTime = created_time
	deployment.Uid = uid
	deployment.Labels = labels

	resDeploymentOverview.Info = deployment

	//pods
	// pod > ownerReferences[] > kind:"RepllicaSet", name,
	// Replicaset에서
	// > Deployement 검색 (이름/Uid)
	// Pod에서
	// > ownerreferences[{kind:"Deployment",name}] >

	// replicasets
	// http://192.168.0.152:31635/apis/apps/v1/namespaces/kube-system/replicasets?clustername=cluster2
	replURL := "https://" + openmcpURL + "/apis/apps/v1/namespaces/" + projectName + "/replicasets?clustername=" + clusterName
	go CallAPI(token, replURL, ch)
	replResult := <-ch
	// fmt.Println(deploymentResult)
	replData := replResult.data
	replItems := replData["items"].([]interface{})

	// find deployements within replicasets
	replUIDs := []string{}
	for _, element := range replItems {
		kind := GetStringElement(element, []string{"metadata", "ownerReferences", "kind"})
		name := GetStringElement(element, []string{"metadata", "ownerReferences", "name"})
		if kind == "Deployment" && name == deploymentName {
			uid := GetStringElement(element, []string{"metadata", "uid"})
			replUIDs = append(replUIDs, uid)
		}
	}

	//openmcp-apiserver-b84bf5cc7
	//ab2a2995-8dca-41ce-aead-8e112d75e3fe

	// find pods within deployments
	// replicasets
	// http://192.168.0.152:31635/apis/apps/v1/namespaces/kube-system/replicasets?clustername=cluster2
	podURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/pods?clustername=" + clusterName
	go CallAPI(token, podURL, ch)
	podResult := <-ch
	podData := podResult.data
	podItems := podData["items"].([]interface{})

	// fmt.Println("replUIDs : ", replUIDs)
	for _, element := range podItems {
		kind := GetStringElement(element, []string{"metadata", "ownerReferences", "kind"})
		if kind == "ReplicaSet" {
			uid := GetStringElement(element, []string{"metadata", "ownerReferences", "uid"})
			for _, item := range replUIDs {
				if item == uid {
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

					resDeploymentOverview.Pods = append(resDeploymentOverview.Pods, pod)
				}
			}
		}
	}

	//ports
	port := PortInfo{}

	containers := GetArrayElement(deploymentData, []string{"spec", "template", "spec", "containers"})
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
			resDeploymentOverview.Ports = append(resDeploymentOverview.Ports, port)
		}
	}
	if len(resDeploymentOverview.Ports) == 0 {
		resDeploymentOverview.Ports = []PortInfo{}
	}

	//events
	eventURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/events?clustername=" + clusterName
	go CallAPI(token, eventURL, ch)
	eventResult := <-ch
	eventData := eventResult.data
	eventItems := eventData["items"].([]interface{})
	resDeploymentOverview.Events = []Event{}

	if len(eventItems) > 0 {
		event := Event{}
		for _, element := range eventItems {
			kind := GetStringElement(element, []string{"involvedObject", "kind"})
			objectName := GetStringElement(element, []string{"involvedObject", "name"})
			if kind == "Deployment" && objectName == deploymentName {
				event.Typenm = GetStringElement(element, []string{"type"})
				event.Reason = GetStringElement(element, []string{"reason"})
				event.Message = GetStringElement(element, []string{"message"})
				// event.Time = GetStringElement(element, []string{"metadata", "creationTimestamp"})
				event.Time = GetStringElement(element, []string{"lastTimestamp"})
				event.Object = kind
				event.Project = projectName

				resDeploymentOverview.Events = append(resDeploymentOverview.Events, event)
			}
		}
	}

	json.NewEncoder(w).Encode(resDeploymentOverview)
}

func GetDeploymentReplicaStatus(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	vars := mux.Vars(r)
	cluster := vars["clusterName"]
	projectName := vars["projectName"]
	deploymentName := vars["deploymentName"]

	resReplicaStatus := ReplicaStatus{}
	// http://192.168.0.152:31635/apis/apps/v1/namespaces/openmcp/deployments/openmcp-deployment3?clustername=cluster1
	deploymentURL := "https://" + openmcpURL + "/apis/apps/v1/namespaces/" + projectName + "/deployments/" + deploymentName + "?clustername=" + cluster

	fmt.Println(deploymentURL)
	go CallAPI(token, deploymentURL, ch)
	deploymentResult := <-ch
	deploymentData := deploymentResult.data

	// get deployement Information
	namespace := GetStringElement(deploymentData, []string{"metadata", "namespace"})

	// unavailableReplicas := GetFloat64Element(deploymentData, []string{"status", "unavailableReplicas"})
	readyReplicas := GetFloat64Element(deploymentData, []string{"status", "readyReplicas"})
	replicas := GetFloat64Element(deploymentData, []string{"status", "replicas"})

	resReplicaStatus.Cluster = cluster
	resReplicaStatus.Project = namespace
	resReplicaStatus.Deployment = deploymentName
	resReplicaStatus.Replicas = int(replicas)
	resReplicaStatus.ReadyReplicas = int(readyReplicas)

	json.NewEncoder(w).Encode(resReplicaStatus)
}
