package util

import (
	"encoding/json"
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	"openmcp/openmcp/util"
	"openmcp/openmcp/util/clusterManager"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
	"strings"
)

func CheckAlreadyJoinClusterWithIP(memberIP string) (bool, error){
	//openmcpIP := GetOutboundIP()
	openmcpIP := GetEndpointIP()
	kc := GetKubeConfig("/mnt/openmcp/"+openmcpIP+"/members/unjoin/"+memberIP+"/config/config")

	clusterName := ""
	for _, cluster := range kc.Clusters {
		if strings.Contains(cluster.Cluster.Server, memberIP) {
			clusterName = cluster.Name
			break
		}
	}

	kubeconfig, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		return false, err
	}
	genClient := genericclient.NewForConfigOrDie(kubeconfig)
	clusterList := clusterManager.ListKubeFedClusters(genClient, "kube-federation-system")


	for _, cluster := range clusterList.Items {
		if strings.Contains(cluster.Spec.APIEndpoint, memberIP) {
			fmt.Println("Already Joined Cluster IP : ", memberIP)
			return true, nil
		}
		if clusterName == cluster.Name {
			fmt.Println("Already Joined Cluster Name : ", clusterName)
			return true, nil
		}
	}
	return false, nil

}
func CheckAlreadyJoinClusterWithPublicClusterName(clusterName, platform string) (bool, error){
	kubeconfig, _ := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	genClient := genericclient.NewForConfigOrDie(kubeconfig)
	clusterList := clusterManager.ListKubeFedClusters(genClient, "kube-federation-system")

	apiEndpoint := ""
	if platform == "gke" {
		s, err := util.CmdExec("gcloud container clusters list")
		if err != nil {
			fmt.Println(err)
		}
		gkeClusterInfo := strings.Split(s, "\n")

		for i := 1; i< len(gkeClusterInfo)-1; i++ {

			ss := strings.Fields(gkeClusterInfo[i])

			if len(ss) < 8 {
				continue
			}
			if ss[0] == clusterName {
				apiEndpoint = "https://"+ss[3]
				break
			}

		}
	} else if platform == "eks"{

		ss, err := util.CmdExec("aws eks describe-cluster --name "+clusterName+" | cat")

		if err != nil {
			return false, err
		}
		jsonData := make(map[string]interface{})
		err = json.Unmarshal([]byte(ss), &jsonData)
		if err != nil {
			return false, err
		}
		clusterInfo := jsonData["cluster"].(map[string]interface{})

		if _, ok := clusterInfo["endpoint"]; ok {
			apiEndpoint = clusterInfo["endpoint"].(string)
		} else {
			return false, fmt.Errorf("eks have not cluster")
		}

		apiEndpoint = strings.ToLower(apiEndpoint)


	} else {

	}

	for _, cluster := range clusterList.Items {

		if strings.Contains(cluster.Spec.APIEndpoint, apiEndpoint) {
			fmt.Println("Already Joined Cluster IP : ", apiEndpoint)
			return true, nil
		}

		if cluster.Name == clusterName {
			fmt.Println("Already Joined Cluster Name : ", clusterName)
			return true, nil
		}

	}
	fmt.Println("Not Joined Cluster IP : ", apiEndpoint, ", Name : ", clusterName)
	return false, nil

}
