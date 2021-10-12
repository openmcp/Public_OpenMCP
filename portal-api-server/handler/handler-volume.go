package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

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
	resPod := PodRes{}
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

	resVolumeOverview.MountedBy = resPod.Pods

	//Events
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
				event.Time = GetStringElement(element, []string{"lastTimestamp"})
				event.Object = kind
				event.Project = projectName

				resVolumeOverview.Events = append(resVolumeOverview.Events, event)
			}
		}
	}

	json.NewEncoder(w).Encode(resVolumeOverview)
}
