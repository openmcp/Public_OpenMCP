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
func GetConfigmaps(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	vars := mux.Vars(r)
	clusterName := vars["clusterName"]
	projectName := vars["projectName"]

	resConfigmap := ConfigmapRes{}
	configMap := ConfigmapInfo{}
	configMapURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/configmaps?clustername=" + clusterName
	// fmt.Println(configMapURL)

	go CallAPI(token, configMapURL, ch)

	configMapResult := <-ch
	configMapData := configMapResult.data
	configMapItems := configMapData["items"].([]interface{})

	if configMapItems != nil && len(configMapItems) > 0 {
		for _, element := range configMapItems {
			name := GetStringElement(element, []string{"metadata", "name"})
			namespace := GetStringElement(element, []string{"metadata", "namespace"})
			createdTime := GetStringElement(element, []string{"metadata", "creationTimestamp"})

			keys := ""
			keysObject := GetInterfaceElement(element, []string{"data"})
			if keysObject != nil {
				i := 0
				for key, _ := range keysObject.(map[string]interface{}) {
					i++
					if i == len(keysObject.(map[string]interface{})) {
						keys = keys + key
					} else {
						keys = keys + key + "|"
					}
				}
			} else {
				keys = "-"
			}
			configMap.Name = name
			configMap.Project = namespace
			configMap.Keys = keys
			configMap.CreatedTime = createdTime

			resConfigmap.Configmaps = append(resConfigmap.Configmaps, configMap)
		}
	}
	json.NewEncoder(w).Encode(resConfigmap.Configmaps)
}

func GetConfigmapOverView(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	vars := mux.Vars(r)
	clusterName := vars["clusterName"]
	projectName := vars["projectName"]
	configmapName := vars["configmapName"]

	resConfigmapOverView := ConfigmapOverView{}
	configMap := ConfigmapInfo{}
	configMapURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/configmaps/" + configmapName + "?clustername=" + clusterName
	// fmt.Println(configMapURL)

	go CallAPI(token, configMapURL, ch)

	configMapResult := <-ch
	configMapData := configMapResult.data

	if configMapData != nil {
		name := GetStringElement(configMapData, []string{"metadata", "name"})
		namespace := GetStringElement(configMapData, []string{"metadata", "namespace"})
		createdTime := GetStringElement(configMapData, []string{"metadata", "creationTimestamp"})
		configMap.Name = name
		configMap.Project = namespace
		configMap.CreatedTime = createdTime

		resConfigmapOverView.Info = configMap

		dataObject := GetInterfaceElement(configMapData, []string{"data"})
		if dataObject != nil {
			data := Data{}
			for key, value := range dataObject.(map[string]interface{}) {
				data.Key = key
				data.Value = value.(string)
				resConfigmapOverView.Data = append(resConfigmapOverView.Data, data)
			}
		}
	}
	json.NewEncoder(w).Encode(resConfigmapOverView)
}
