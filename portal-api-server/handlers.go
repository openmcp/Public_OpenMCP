package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"portal-api-server/cloud"
	"portal-api-server/resource"

	"github.com/gorilla/mux"
)

var targetURL = "172.17.1.241:7070"

func Test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	node, ok1 := r.URL.Query()["node"]
	cluster, ok2 := r.URL.Query()["cluster"]
	if !ok1 || !ok2 || len(cluster[0]) < 1 || len(node[0]) < 1 {
		log.Println("Url Params are missing")
		w.Write([]byte("Url Param are missing"))
	} else {
		result := cloud.AddNode(node[0])
		if err := json.NewEncoder(w).Encode(result); err != nil {
			panic(err)
		}
		go cloud.GetNodeState(&result.InstanceID, node[0], cluster[0])

		// id := "i-09ce908be9488f77c"
		// cloud.GetNodeState(&id)
	}

}

func WorkloadsDeploymentsOverviewList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resource.ListResource()); err != nil {
		panic(err)
	}

}

func WorkloadsPodsOverviewList(w http.ResponseWriter, r *http.Request) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	vars := mux.Vars(r)

	var client http.Client
	resp, err := client.Get("https://" + targetURL + "/seedcontainer/api/v1/clusters/" + vars["clusterName"] + "/daemonsets/list")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(bodyBytes)
	}
}

func getDeploymentList(w http.ResponseWriter, r *http.Request) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	vars := mux.Vars(r)

	var client http.Client
	resp, err := client.Get("https://" + targetURL + "/seedcontainer/api/v1/clusters/" + vars["clusterName"] + "/deployments/list")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(bodyBytes)
	}
}

func getDeploymentDetail(w http.ResponseWriter, r *http.Request) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	vars := mux.Vars(r)

	var client http.Client

	callUrl := "https://" + targetURL + "/seedcontainer/api/v1/clusters/" + vars["clusterName"] + "/namespaces/" + vars["namespaceName"] + "/deployments/" + vars["deploymentName"] + "/detail"
	//fmt.Print(callUrl)

	resp, err := client.Get(callUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(bodyBytes)
	}
}

func getDeploymentYaml(w http.ResponseWriter, r *http.Request) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	vars := mux.Vars(r)

	var client http.Client

	callUrl := "https://" + targetURL + "/seedcontainer/api/v1/clusters/" + vars["clusterName"] + "/namespaces/" + vars["namespaceName"] + "/deployments/" + vars["deploymentName"] + "/yaml"
	//fmt.Print(callUrl)

	resp, err := client.Get(callUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(bodyBytes)
	}
}

func getDeploymentEvent(w http.ResponseWriter, r *http.Request) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	vars := mux.Vars(r)

	var client http.Client

	callUrl := "https://" + targetURL + "/seedcontainer/api/v1/clusters/" + vars["clusterName"] + "/namespaces/" + vars["namespaceName"] + "/deployments/" + vars["deploymentName"] + "/events"
	//fmt.Print(callUrl)

	resp, err := client.Get(callUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(bodyBytes)
	}
}
