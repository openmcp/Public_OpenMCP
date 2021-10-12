package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Projects(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	clusterURL := "https://" + openmcpURL + "/apis/core.kubefed.io/v1beta1/kubefedclusters?clustername=openmcp"
	// clusterURL := "https://" + openmcpURL + "/apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcpclusters?clustername=openmcp"

	go CallAPI(token, clusterURL, ch)
	clusters := <-ch
	clusterData := clusters.data

	resProject := ProjectRes{}
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
		projectURL := "https://" + openmcpURL + "/api/v1/namespaces?clustername=" + clusterName
		go CallAPI(token, projectURL, ch)
		projectResult := <-ch
		projectData := projectResult.data
		projectItems := projectData["items"].([]interface{})

		// get podUsage counts by nodename groups
		for _, element := range projectItems {
			project := ProjectInfo{}
			projectName := GetStringElement(element, []string{"metadata", "name"})
			// element.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
			createdTime := GetStringElement(element, []string{"metadata", "creationTimestamp"})
			// element.(map[string]interface{})["metadata"].(map[string]interface{})["creationTimestamp"].(string)
			status := GetStringElement(element, []string{"status", "phase"})
			// element.(map[string]interface{})["status"].(map[string]interface{})["phase"].(string)

			labels := make(map[string]interface{})
			labelCheck := GetInterfaceElement(element, []string{"metadata", "labels"})
			// element.(map[string]interface{})["metadata"].(map[string]interface{})["labels"]
			if labelCheck == nil {
				//undefined lable
				labels = map[string]interface{}{}
			} else {
				for key, val := range labelCheck.(map[string]interface{}) {
					// fmt.Println(key, val)
					labels[key] = val
				}
			}
			project.Name = projectName
			project.Status = status
			project.Cluster = clusterName
			project.CreatedTime = createdTime
			project.Labels = labels

			resProject.Projects = append(resProject.Projects, project)

		}
	}

	json.NewEncoder(w).Encode(resProject.Projects)
}

func GetProjectOverview(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	vars := mux.Vars(r)
	clusterName := vars["clusterName"]
	projectName := vars["projectName"]

	resProjectOverview := ProjectOverview{}

	projectURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "?clustername=" + clusterName
	go CallAPI(token, projectURL, ch)
	projectResult := <-ch
	projectData := projectResult.data
	// projectItems := projectData["items"].([]interface{})

	// get podUsage counts by nodename groups
	project := ProjectInfo{}
	createdTime := GetStringElement(projectData, []string{"metadata", "creationTimestamp"})
	status := GetStringElement(projectData, []string{"status", "phase"})
	uid := GetStringElement(projectData, []string{"metadata", "uid"})
	labels := make(map[string]interface{})
	labelCheck := GetInterfaceElement(projectData, []string{"metadata", "labels"})
	if labelCheck == nil {
		//undefined lable
		labels = map[string]interface{}{}
	} else {
		for key, val := range labelCheck.(map[string]interface{}) {
			labels[key] = val
		}
	}
	project.Name = projectName
	project.Status = status
	project.Cluster = clusterName
	project.CreatedTime = createdTime
	project.Labels = labels
	project.UID = uid
	resProjectOverview.Info = project

	// Project Resources
	// Pods, Deployments, StatefulSets, DaemonSets, Jobs (count, UnhealthyCount)
	// 1.pod //////////////////////////////////////////////////////////
	podResources := ProjectResourceType{}
	podURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/pods?clustername=" + clusterName
	go CallAPI(token, podURL, ch)
	podResult := <-ch
	podData := podResult.data
	podItems := podData["items"].([]interface{})
	for _, element := range podItems {
		//Pending, Running, Succeeded, Failed, Unknown
		status := GetStringElement(element, []string{"status", "phase"})
		if status != "Running" || status != "Succeeded" {
			podResources.Abnormal++
		}
	}
	podResources.Name = "Pods"
	podResources.Total = len(podItems)
	resProjectOverview.ProjectResource = append(resProjectOverview.ProjectResource, podResources)

	// 2.deployments //////////////////////////////////////////////////////////
	deploymentResources := ProjectResourceType{}
	deploymentURL := "https://" + openmcpURL + "/apis/apps/v1/namespaces/" + projectName + "/deployments?clustername=" + clusterName
	go CallAPI(token, deploymentURL, ch)
	deploymentResult := <-ch
	deploymentData := deploymentResult.data
	deploymentItems := deploymentData["items"].([]interface{})

	for _, element := range deploymentItems {
		unavailableReplicas := GetInterfaceElement(element, []string{"status", "unavailableReplicas"})

		if unavailableReplicas != nil && unavailableReplicas.(float64) > 0 {
			deploymentResources.Abnormal++
		}
	}
	deploymentResources.Name = "Deployments"
	deploymentResources.Total = len(deploymentItems)
	resProjectOverview.ProjectResource = append(resProjectOverview.ProjectResource, deploymentResources)

	// 3.StatefulSets //////////////////////////////////////////////////////////
	stateFulSetResources := ProjectResourceType{}
	stateFulSetURL := "https://" + openmcpURL + "/apis/apps/v1/namespaces/" + projectName + "/statefulsets?clustername=" + clusterName
	go CallAPI(token, stateFulSetURL, ch)
	stateFulSetResult := <-ch
	stateFulSetData := stateFulSetResult.data
	stateFulSetItems := stateFulSetData["items"].([]interface{})

	for _, element := range stateFulSetItems {
		replicas := GetFloat64Element(element, []string{"status", "replicas"})
		readyReplicas := GetFloat64Element(element, []string{"status", "readyReplicas"})

		abnormals := replicas - readyReplicas

		if replicas > readyReplicas || abnormals > 0 {
			stateFulSetResources.Abnormal++
		}
	}
	stateFulSetResources.Name = "StatefulSets"
	stateFulSetResources.Total = len(stateFulSetItems)
	resProjectOverview.ProjectResource = append(resProjectOverview.ProjectResource, stateFulSetResources)

	// 4.DaemonSets //////////////////////////////////////////////////////////
	daemonSetResources := ProjectResourceType{}
	daemonSetURL := "https://" + openmcpURL + "/apis/apps/v1/namespaces/" + projectName + "/deployments?clustername=" + clusterName
	go CallAPI(token, daemonSetURL, ch)
	daemonSetResult := <-ch
	daemonSetData := daemonSetResult.data
	daemonSetItems := daemonSetData["items"].([]interface{})

	for _, element := range daemonSetItems {
		numberUnavailable := GetInterfaceElement(element, []string{"status", "numberUnavailable"})
		if numberUnavailable != nil && numberUnavailable.(int) > 0 {
			daemonSetResources.Abnormal++
		}
	}
	daemonSetResources.Name = "DaemonSets"
	daemonSetResources.Total = len(daemonSetItems)
	resProjectOverview.ProjectResource = append(resProjectOverview.ProjectResource, daemonSetResources)

	// 5.Jobs //////////////////////////////////////////////////////////
	jobResources := ProjectResourceType{}
	jobURL := "https://" + openmcpURL + "/apis/apps/v1/namespaces/" + projectName + "/deployments?clustername=" + clusterName
	go CallAPI(token, jobURL, ch)
	jobResult := <-ch
	jobData := jobResult.data
	jobItems := jobData["items"].([]interface{})

	for _, element := range jobItems {
		//Complete, Failed
		statusType := GetStringElement(element, []string{"status", "type"})
		if statusType == "Failed" {
			jobResources.Abnormal++
		}
	}
	jobResources.Name = "Jobs"
	jobResources.Total = len(jobItems)
	resProjectOverview.ProjectResource = append(resProjectOverview.ProjectResource, jobResources)

	// Usage Top5
	// workload (Deployment, Replicaset, statefulSet, DemonSet, Job, CronJob)
	usageResult := GetInfluxPodTop5(clusterName, projectName)
	resProjectOverview.UsageTop5 = usageResult

	results := GetInfluxDBPod10mMetric(clusterName, projectName)
	resProjectOverview.PhysicalResources = results

	json.NewEncoder(w).Encode(resProjectOverview)
}

func AddProject(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// w.WriteHeader(http.StatusOK)

	fmt.Println("AddProject")

	data := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지

	clusters := data["clusters"]
	project := data["project"].(string)

	type Metadata struct {
		Name string `json:"name"`
	}

	type jsonBody struct {
		ApiVersion string   `json:"apiVersion"`
		Kind       string   `json:"kind"`
		Metadata   Metadata `json:"metadata"`
	}
	// 	{
	//     "apiVersion": "v1",
	//     "kind": "Namespace",
	//     "metadata": {
	//         "name": "test-namespace1"
	//     }
	// }

	body := jsonBody{}
	metaInfo := Metadata{}
	metaInfo.Name = project
	body.ApiVersion = "v1"
	body.Kind = "Namespace"
	body.Metadata = metaInfo

	var jsonErrs []jsonErr

	for _, element := range clusters.([]interface{}) {

		clusterName := element.(map[string]interface{})["name"].(string)
		//https://192.168.0.152:30000/api/v1/namespaces?clustername=cluster1
		projectURL := "https://" + openmcpURL + "/api/v1/namespaces?clustername=" + clusterName
		// fmt.Println(clusterName)
		// fmt.Println(projectURL)
		// fmt.Println(project)

		resp, err := CallPostAPI(projectURL, "application/json", body)
		var msg jsonErr

		if err != nil {
			msg = jsonErr{503, "failed", "request fail | " + body.Kind + " | " + metaInfo.Name}
			// json.NewEncoder(w).Encode(msg)
		}

		var data map[string]interface{}
		json.Unmarshal([]byte(resp), &data)
		if data != nil {
			if data["kind"].(string) == "Status" {
				msg = jsonErr{501, "failed", data["message"].(string) + body.Kind + " | " + metaInfo.Name}
				// json.NewEncoder(w).Encode(msg)
			} else {
				msg = jsonErr{200, "success", "Resource Created" + body.Kind + " | " + metaInfo.Name}
				// json.NewEncoder(w).Encode(msg)
			}
		}

		jsonErrs = append(jsonErrs, msg)

	}
	json.NewEncoder(w).Encode(jsonErrs)

}
