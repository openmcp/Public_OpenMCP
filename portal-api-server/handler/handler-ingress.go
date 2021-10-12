package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Ingress(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	clusterurl := "https://" + openmcpURL + "/apis/core.kubefed.io/v1beta1/kubefedclusters?clustername=openmcp"
	go CallAPI(token, clusterurl, ch)
	clusters := <-ch
	clusterData := clusters.data

	resIngress := IngerssRes{}
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
		ingress := IngerssInfo{}
		ingressURL := "https://" + openmcpURL + "/apis/networking.k8s.io/v1beta1/ingresses?clustername=" + clusterName
		go CallAPI(token, ingressURL, ch)
		ingressResult := <-ch
		ingressData := ingressResult.data
		ingressItems := ingressData["items"].([]interface{})

		if len(ingressItems) > 0 {
			for _, element := range ingressItems {
				name := GetStringElement(element, []string{"metadata", "name"})
				namespace := GetStringElement(element, []string{"metadata", "namespace"})
				ipAddress := GetStringElement(element, []string{"status", "loadBalancer", "ingress", "ip"})
				createTime := GetStringElement(element, []string{"metadata", "creationTimestamp"})

				ingress.Name = name
				ingress.Project = namespace
				ingress.Address = ipAddress
				ingress.CreatedTime = createTime
				ingress.Cluster = clusterName

				resIngress.Ingress = append(resIngress.Ingress, ingress)
			}
		}
	}
	json.NewEncoder(w).Encode(resIngress.Ingress)
}

func GetIngressInProject(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	vars := mux.Vars(r)
	clusterName := vars["clusterName"]
	projectName := vars["projectName"]

	resIngress := IngerssRes{}
	ingress := IngerssInfo{}
	// "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/services?
	ingressURL := "https://" + openmcpURL + "/apis/networking.k8s.io/v1beta1/namespaces/" + projectName + "/ingresses?clustername=" + clusterName

	go CallAPI(token, ingressURL, ch)

	ingressResult := <-ch
	ingressData := ingressResult.data
	ingressItems := ingressData["items"].([]interface{})

	if ingressItems != nil {
		for _, element := range ingressItems {
			name := GetStringElement(element, []string{"metadata", "name"})
			namespace := GetStringElement(element, []string{"metadata", "namespace"})
			ipAddress := GetStringElement(element, []string{"status", "loadBalancer", "ingress", "ip"})
			createdTime := GetStringElement(element, []string{"metadata", "creationTimestamp"})

			ingress.Name = name
			ingress.Project = namespace
			ingress.Address = ipAddress
			ingress.CreatedTime = createdTime
			ingress.Cluster = clusterName

			resIngress.Ingress = append(resIngress.Ingress, ingress)
		}
	}
	json.NewEncoder(w).Encode(resIngress.Ingress)
}

func GetIngressOverview(w http.ResponseWriter, r *http.Request) {
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	vars := mux.Vars(r)
	clusterName := vars["clusterName"]
	projectName := vars["projectName"]
	ingressName := vars["ingressName"]

	resIngressOverView := IngressOverView{}
	ingress := IngerssInfo{}

	// "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/services?
	ingressURL := "https://" + openmcpURL + "/apis/networking.k8s.io/v1beta1/namespaces/" + projectName + "/ingresses/" + ingressName + "?clustername=" + clusterName

	go CallAPI(token, ingressURL, ch)

	ingressResult := <-ch
	ingressData := ingressResult.data

	resIngressOverView.Info = IngerssInfo{}
	resIngressOverView.Rules = []Rules{}

	if ingressData != nil {
		name := GetStringElement(ingressData, []string{"metadata", "name"})
		namespace := GetStringElement(ingressData, []string{"metadata", "namespace"})
		ipAddress := GetStringElement(ingressData, []string{"status", "loadBalancer", "ingress", "ip"})
		createdTime := GetStringElement(ingressData, []string{"metadata", "creationTimestamp"})

		rules := Rules{}
		tlsArrObj := GetArrayElement(ingressData, []string{"tls"})
		if len(tlsArrObj) > 0 {
			for i, item := range tlsArrObj {
				if len(tlsArrObj)-1 == i {
					rules.Secret = rules.Secret + GetStringElement(item, []string{"secretName"})
				} else {
					rules.Secret = rules.Secret + GetStringElement(item, []string{"secretName"}) + "|"
				}
			}
			// 	hosts :[https-example.foo.com]
			// 	secretName: testsecret-tls
			// ]
		} else {
			rules.Secret = "-"
		}

		ruleArrObj := GetArrayElement(ingressData, []string{"spec", "rules"})
		if len(ruleArrObj) > 0 {
			for _, element := range ruleArrObj {
				rules.Domain = GetStringElement(element, []string{"host"})
				rules.Protocol = "http"

				pathArrObj := GetArrayElement(element, []string{"http", "paths"})
				for _, item := range pathArrObj {
					rules.Path = rules.Path + GetStringElement(item, []string{"path"})
					rules.Services = rules.Services + GetStringElement(item, []string{"backend", "serviceName"})

					port := GetFloat64Element(item, []string{"backend", "servicePort"})
					rules.Port = rules.Port + strconv.FormatFloat(port, 'f', -1, 64)
				}

				resIngressOverView.Rules = append(resIngressOverView.Rules, rules)
			}
		}

		ingress.Name = name
		ingress.Project = namespace
		ingress.Address = ipAddress
		ingress.CreatedTime = createdTime
		ingress.Cluster = clusterName

		resIngressOverView.Info = ingress
	}

	//events
	// http://192.168.0.152:31635/api/v1/namespaces/ingress-nginx/events?clustername=cluster1
	eventURL := "https://" + openmcpURL + "/api/v1/namespaces/" + projectName + "/events?clustername=" + clusterName

	go CallAPI(token, eventURL, ch)
	eventResult := <-ch
	eventData := eventResult.data
	eventItems := eventData["items"].([]interface{})
	resIngressOverView.Events = []Event{}

	if len(eventItems) > 0 {
		event := Event{}
		for _, element := range eventItems {
			kind := GetStringElement(element, []string{"involvedObject", "kind"})
			objectName := GetStringElement(element, []string{"involvedObject", "name"})
			if kind == "Ingress" && objectName == ingressName {
				event.Typenm = GetStringElement(element, []string{"type"})
				event.Reason = GetStringElement(element, []string{"reason"})
				event.Message = GetStringElement(element, []string{"message"})
				// event.Time = GetStringElement(element, []string{"metadata", "creationTimestamp"})
				event.Time = GetStringElement(element, []string{"lastTimestamp"})
				event.Object = kind
				event.Project = projectName

				resIngressOverView.Events = append(resIngressOverView.Events, event)
			}
		}
	}

	json.NewEncoder(w).Encode(resIngressOverView)
}
