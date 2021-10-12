package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Secrets
// name : items ("metadata","name") string
// namespace : items ("metatdata","namespace") string
// type : items ("status","capacity" storage") string
// createdTime items ("metadata","creationTimestamp")

// persistent volume claim (PVC)
func GetSecrets(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	vars := mux.Vars(r)
	clusterName := vars["clusterName"]
	projectName := vars["projectName"]

	resSecret := SecretRes{}
	secret := SecretInfo{}
	// http://192.168.0.152:31635/api/v1/namespaces/openmcp/persistentvolumeclaims?clustername=openmcp
	secretURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/secrets?clustername=" + clusterName

	go CallAPI(token, secretURL, ch)

	secretResult := <-ch
	secretData := secretResult.data
	secretItems := secretData["items"].([]interface{})

	if secretItems != nil {
		for _, element := range secretItems {
			name := GetStringElement(element, []string{"metadata", "name"})
			namespace := GetStringElement(element, []string{"metadata", "namespace"})
			scretType := GetStringElement(element, []string{"type"})
			createdTime := GetStringElement(element, []string{"metadata", "creationTimestamp"})

			secret.Name = name
			secret.Project = namespace
			secret.Type = scretType
			secret.CreatedTime = createdTime

			resSecret.Secrets = append(resSecret.Secrets, secret)
		}
	}
	json.NewEncoder(w).Encode(resSecret.Secrets)
}

func GetSecretOverView(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	vars := mux.Vars(r)
	clusterName := vars["clusterName"]
	projectName := vars["projectName"]
	secretName := vars["secretName"]

	resSecretOverview := SecretOverView{}
	secret := SecretInfo{}

	// http://192.168.0.152:31635/api/v1/namespaces/openmcp/secrets/default-token-fwpzk?clustername=openmcp

	// http://192.168.0.152:31635/api/v1/namespaces/openmcp/secrets/default-token-fwpzk?clustername=openmcp
	secretURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/secrets/" + secretName + "?clustername=" + clusterName
	fmt.Println(secretURL)
	go CallAPI(token, secretURL, ch)

	secretResult := <-ch
	secretData := secretResult.data
	// secretItems := secretData["items"].([]interface{})

	if secretData != nil {

		name := GetStringElement(secretData, []string{"metadata", "name"})
		namespace := GetStringElement(secretData, []string{"metadata", "namespace"})
		scretType := GetStringElement(secretData, []string{"type"})
		createdTime := GetStringElement(secretData, []string{"metadata", "creationTimestamp"})

		secret.Name = name
		secret.Project = namespace
		secret.Type = scretType
		secret.CreatedTime = createdTime

		resSecretOverview.Info = secret

		dataObject := GetInterfaceElement(secretData, []string{"data"})
		if dataObject != nil {
			data := Data{}
			for key, value := range dataObject.(map[string]interface{}) {
				data.Key = key
				data.Value = value.(string)
				resSecretOverview.Data = append(resSecretOverview.Data, data)
			}
		}

	}
	json.NewEncoder(w).Encode(resSecretOverview)
}
