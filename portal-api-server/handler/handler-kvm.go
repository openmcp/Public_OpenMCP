package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetKVMNodes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// http://192.168.0.89:4885/apis/getkvmnodes?agenturl=192.168.0.96
	agentURL := r.URL.Query().Get("agenturl")

	var client http.Client
	resp, err := client.Get("https://" + agentURL + ":10000/getkvmlists")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	defer resp.Body.Close()

	var data interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	json.NewEncoder(w).Encode(&data)
}

func StartKVMNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// http://192.168.0.89:4885/apis/startkvmnode?agenturl=192.168.0.96&node=rancher
	// agentURL := r.URL.Query().Get("agentURL")
	// nodeName := r.URL.Query().Get("node")

	//POST
	body := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지

	agentURL := body["agentURL"].(string)
	nodeName := body["node"].(string)

	// fmt.Println(agentURL)
	// fmt.Println(nodeName)

	var client http.Client
	resp, err := client.Get("https://" + agentURL + ":10000/kvmstartnode?node=" + nodeName)
	if err != nil {

		errorJson := jsonErr{500, "agent connect fail", err.Error()}
		json.NewEncoder(w).Encode(errorJson)
	}

	defer resp.Body.Close()

	var data interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	json.NewEncoder(w).Encode(&data)
}

func StopKVMNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// http://192.168.0.89:4885/apis/stopkvmnode?agenturl=192.168.0.96&node=rancher
	// agentURL := r.URL.Query().Get("agentURL")
	// nodeName := r.URL.Query().Get("node")

	//POST
	body := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지

	agentURL := body["agentURL"].(string)
	nodeName := body["node"].(string)

	// fmt.Println(agentURL)
	// fmt.Println(nodeName)

	var client http.Client
	resp, err := client.Get("https://" + agentURL + ":10000/kvmstopnode?node=" + nodeName)
	if err != nil {
		errorJson := jsonErr{500, "agent connect fail", err.Error()}
		json.NewEncoder(w).Encode(errorJson)
	}

	defer resp.Body.Close()

	var data interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	json.NewEncoder(w).Encode(&data)
}

func ChangeKVMNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// // http://192.168.0.89:4885/apis/changekvmnode?agenturl=192.168.0.96&node=rancher&cpu=4&mem=8088
	// agentURL := r.URL.Query().Get("agenturl")
	// nodeName := r.URL.Query().Get("node")
	// vCpu := r.URL.Query().Get("cpu")
	// memory := r.URL.Query().Get("mem")

	body := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지

	agentURL := body["agentURL"].(string)
	nodeName := body["node"].(string)
	vCpu := body["cpu"].(string)
	memory := body["memory"].(string)

	var client http.Client
	resp, err := client.Get("https://" + agentURL + ":10000/changekvmnode?node=" + nodeName + "&cpu=" + vCpu + "&mem=" + memory)

	if err != nil {
		errorJson := jsonErr{500, "agent connect fail", err.Error()}
		json.NewEncoder(w).Encode(errorJson)
	}

	defer resp.Body.Close()

	var data interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	json.NewEncoder(w).Encode(&data)
}

func CreateKVMNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	//GET
	// http://192.168.0.89:4885/apis/createkvmnode?agenturl=192.168.0.96&template=ubuntu16.04-clean&newvm=newvmvmvmvmvmvm

	//Post
	body := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지

	agentURL := body["agentURL"].(string)
	newvm := body["newvm"].(string)
	template := body["template"].(string)
	// cluster := body["cluster"].(string)

	clusterMaster := body["master"].(string)
	mPass := body["mpass"].(string)
	wPass := body["wpass"].(string)

	// fmt.Println(agentURL)
	// fmt.Println(newvm)
	// fmt.Println(template)
	// fmt.Println(cluster)

	// agentURL := r.URL.Query().Get("agenturl")
	// newvm := r.URL.Query().Get("newvm")
	// template := r.URL.Query().Get("template")

	var client http.Client
	resp, err := client.Get("https://" + agentURL + ":10000/createkvmnode?template=" + template + "&newvm=" + newvm + "&master=" + clusterMaster + "&mpass=" + mPass + "&wpass=" + wPass)

	if err != nil {
		errorJson := jsonErr{500, "Agent Connect Fail", err.Error()}
		fmt.Println("err", errorJson)
		json.NewEncoder(w).Encode(errorJson)
	} else {
		defer resp.Body.Close()
		var data interface{}
		json.NewDecoder(resp.Body).Decode(&data)
		json.NewEncoder(w).Encode(&data)
	}
}

func DeleteKVMNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// http://192.168.0.89:4885/apis/deletekvmnode?agenturl=192.168.0.96&targetvm=newnode-1
	// agentURL := r.URL.Query().Get("agenturl")
	// targetvm := r.URL.Query().Get("targetvm")

	body := GetJsonBody(r.Body)
	agentURL := body["agentURL"].(string)
	targetvm := body["targetvm"].(string)
	mastervm := body["mastervm"].(string)
	mastervmpwd := body["mastervmpwd"].(string)

	// fmt.Println(agentURL)
	// fmt.Println(targetvm)

	var client http.Client
	resp, err := client.Get("https://" + agentURL + ":10000/deletekvmnode?node=" + targetvm + "&mastervm=" + mastervm + "&mastervmpwd=" + mastervmpwd)

	if err != nil {
		errorJson := jsonErr{500, "agent connect fail", err.Error()}
		json.NewEncoder(w).Encode(errorJson)
	}

	defer resp.Body.Close()

	var data interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	json.NewEncoder(w).Encode(&data)
}
