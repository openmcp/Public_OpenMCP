package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Persistent Volume Claim (PVC)
// name : items ("metadata","name") string
// status : items ("status","phase") string
// namespace : items ("metatdata","namespace") string
// capacity : items ("status","capacity" storage") string
// createdTime items ("metadata","creationTimestamp")

// persistent volume claim (PVC)
func GetVolumes(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	vars := mux.Vars(r)
	clusterName := vars["clusterName"]
	projectName := vars["projectName"]

	resVolume := VolumeRes{}
	volume := VolumeInfo{}
	volumeURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/persistentvolumeclaims?clustername=" + clusterName

	go CallAPI(token, volumeURL, ch)

	volumeResult := <-ch
	volumeData := volumeResult.data
	volumeItems := volumeData["items"].([]interface{})

	if volumeItems != nil {
		for _, element := range volumeItems {
			name := GetStringElement(element, []string{"metadata", "name"})
			namespace := GetStringElement(element, []string{"metadata", "namespace"})
			status := GetStringElement(element, []string{"status", "phase"})
			capacity := GetStringElement(element, []string{"status", "capacity", "storage"})
			createdTime := GetStringElement(element, []string{"metadata", "creationTimestamp"})

			volume.Name = name
			volume.Project = namespace
			volume.Status = status
			volume.Capacity = capacity
			volume.CreatedTime = createdTime
			volume.StorageClass = ""
			volume.AccessMode = ""
			resVolume.Volumes = append(resVolume.Volumes, volume)
		}
	}
	json.NewEncoder(w).Encode(resVolume.Volumes)
}

func GetVolumeOverview(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	vars := mux.Vars(r)
	clusterName := vars["clusterName"]
	projectName := vars["projectName"]
	volumeName := vars["volumeName"]

	resVolumeOverview := VolumeOverview{}
	volume := VolumeInfo{}
	volumeURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/persistentvolumeclaims/" + volumeName + "?clustername=" + clusterName

	// fmt.Println(volumeURL)
	go CallAPI(token, volumeURL, ch)

	volumeResult := <-ch
	volumeData := volumeResult.data

	if volumeData != nil {
		name := GetStringElement(volumeData, []string{"metadata", "name"})
		namespace := GetStringElement(volumeData, []string{"metadata", "namespace"})
		status := GetStringElement(volumeData, []string{"status", "phase"})
		capacity := GetStringElement(volumeData, []string{"status", "capacity", "storage"})
		createdTime := GetStringElement(volumeData, []string{"metadata", "creationTimestamp"})
		accessMode := ""
		accessModeArrStr := GetArrayElement(volumeData, []string{"spec", "accessModes"})
		if len(accessModeArrStr) > 0 {
			for i, item := range accessModeArrStr {
				if i == len(accessModeArrStr)-1 {
					accessMode = accessMode + item.(string)
				} else {
					accessMode = accessMode + item.(string) + ", "
				}
			}
		} else {
			accessMode = "-"
		}
		storageClass := GetStringElement(volumeData, []string{"spec", "storageClassName"})

		volume.Name = name
		volume.Project = namespace
		volume.Status = status
		volume.Capacity = capacity
		volume.AccessMode = accessMode
		volume.StorageClass = storageClass
		volume.CreatedTime = createdTime

		resVolumeOverview.Info = volume
	}

	//MountedBy
	// http://192.168.0.152:31635/api/v1/namespaces/webproject/pods?clustername=cluster1
	resPod := PodRes{}

	// http: //192.168.0.152:31635/api/v1/namespaces/kube-system/pods?clustername=cluster2
	podURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/pods?clustername=" + clusterName
	go CallAPI(token, podURL, ch)
	podResult := <-ch
	podData := podResult.data
	podItems := podData["items"].([]interface{})

	if len(podItems) > 0 {

	}

	// get podUsage counts by nodename groups
	for _, element := range podItems {

		isMatchPVC := false
		pvcArrObj := GetArrayElement(element, []string{"spec", "volumes"})
		if len(pvcArrObj) > 0 {
			for _, item := range pvcArrObj {
				claimName := GetStringElement(item, []string{"persistentVolumeClaim", "claimName"})
				if claimName == volumeName {
					isMatchPVC = true
				}
			}
		} else {
			continue
		}

		if !isMatchPVC {
			continue
		}

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

	resVolumeOverview.MountedBy = resPod.Pods

	//Events
	// http://192.168.0.152:31635/api/v1/namespaces/ingress-nginx/events?clustername=cluster1
	eventURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/events?clustername=" + clusterName

	go CallAPI(token, eventURL, ch)
	eventResult := <-ch
	eventData := eventResult.data
	eventItems := eventData["items"].([]interface{})
	resVolumeOverview.Events = []Event{}

	if len(eventItems) > 0 {
		event := Event{}
		for _, element := range eventItems {
			kind := GetStringElement(element, []string{"involvedObject", "kind"})
			objectName := GetStringElement(element, []string{"involvedObject", "name"})
			if kind == "Volume" && objectName == volumeName {
				event.Typenm = GetStringElement(element, []string{"type"})
				event.Reason = GetStringElement(element, []string{"reason"})
				event.Message = GetStringElement(element, []string{"message"})
				// event.Time = GetStringElement(element, []string{"metadata", "creationTimestamp"})
				event.Time = GetStringElement(element, []string{"lastTimestamp"})
				event.Object = kind
				event.Project = projectName

				resVolumeOverview.Events = append(resVolumeOverview.Events, event)
			}
		}
	}

	json.NewEncoder(w).Encode(resVolumeOverview)
}
