package httphandler

import (
	"bytes"
	"context"
	"encoding/json"
	"k8s.io/client-go/rest"
	"log"
	"net/http"
	"sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
)

type JoinableClusterList struct {
	Items []JoinableCluster `json:"items"`
}

type JoinableCluster struct {
	Name      		string      `json:"name"`
	ENDPOINT  	 	string      `json:"endpoint"`

	PlatForm  		string		`json:"platform"`
	Region 			string 		`json:"region"`
	Zone  			string		`json:"zone"`


}

func (h *HttpManager)JoinableHandler(w http.ResponseWriter, r *http.Request) {
	// GET http://10.0.3.20:31635/joinable1


	host_config, _ := rest.InClusterConfig()
	host_client := genericclient.NewForConfigOrDie(host_config)

	kfcList := &v1beta1.KubeFedClusterList{}
	err := host_client.List(context.TODO(), kfcList, "kube-federation-system")
	if err != nil {
		log.Fatal(err)
	}

	defineClusterNames := []string{"cluster1", "cluster2", "cluster3", "gke-cluster1", "eks-cluster1"}
	joinableClusterNames := []string{}

	for _, defineClusterName := range defineClusterNames {
		find := false
		for _, kfc := range kfcList.Items {
			if kfc.Name == defineClusterName {
				find = true
				break
			}
		}
		if !find {
			joinableClusterNames = append(joinableClusterNames, defineClusterName)
		}
	}

	jcl := JoinableClusterList{}
	for _, joinableClusterName := range joinableClusterNames {
		jc := JoinableCluster{}
		jc.Name = joinableClusterName

		if jc.Name == "cluster1"{
			jc.ENDPOINT = "https://10.0.3.50:6443"
			jc.PlatForm = ""
			jc.Region = "KR"
			jc.Zone = "Seoul"
		} else if jc.Name == "cluster2"{
			jc.ENDPOINT = "https://10.0.3.70:6443"
			jc.PlatForm = ""
			jc.Region = "KR"
			jc.Zone = "Gyeonggi-do"
		} else if jc.Name == "cluster3"{
			jc.ENDPOINT = "https://10.0.6.1420:6443"
			jc.PlatForm = ""
			jc.Region = "KR"
			jc.Zone = "Seoul"
		} else if jc.Name == "eks-cluster1"{
			jc.ENDPOINT = "https://a6ea1a0793ef54cefb26e05961dde23d.gr7.eu-west-2.eks.amazonaws.com"
			jc.PlatForm = "eks"
			jc.Region = "eu-west-2"
			jc.Zone = "eu-west-2b,eu-west-2c"
		} else if jc.Name == "gke-cluster1"{
			jc.ENDPOINT = "https://34.67.246.107"
			jc.PlatForm = "gke"
			jc.Region = "us-central1"
			jc.Zone = "us-central1-c"
		}
		jcl.Items = append(jcl.Items, jc)
	}

	bytesJson, _ := json.Marshal(jcl)
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, bytesJson, "", "\t")
	if err != nil {
		panic(err.Error())
	}





	w.Write(prettyJSON.Bytes())
}
