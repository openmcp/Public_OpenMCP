package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"portal-api-server/cloud"
	"portal-api-server/handler"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/eks"
)

var openmcpURL = handler.InitPortalConfig()

func Migration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	clusterurl := "https://" + openmcpURL + "/apis/openmcp.k8s.io/v1alpha1/namespaces/default/migrations?clustername=openmcp"
	resp, err := PostYaml(clusterurl, r.Body)
	defer r.Body.Close()
	if err != nil {
		errmsg := jsonErr{503, "failed", "request fail"}
		json.NewEncoder(w).Encode(errmsg)
	}

	var data map[string]interface{}
	json.Unmarshal([]byte(resp), &data)

	if data["kind"].(string) == "Status" {
		msg := jsonErr{501, "failed", data["message"].(string)}
		json.NewEncoder(w).Encode(msg)
	} else {
		msg := jsonErr{200, "success", "Migration Created"}
		json.NewEncoder(w).Encode(msg)
	}
}

func GetEKSClusterInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// http://192.168.0.51:4885/apis/geteksclusterinfo?region=ap-northeast-2
	// aws test(lkh1434@gmail.com)

	// region := r.URL.Query().Get("region")
	// akid := "AKIAJGFO6OXHRN2H6DSA"
	// secretkey := "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" //
	// // akid := "AKIAVJTB7UPJPEMHUAJR"
	// // secretkey := "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" //

	data := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지

	region := data["region"].(string)
	akid := data["accessKey"].(string)
	secretkey := data["secretKey"].(string)

	sess, err := session.NewSession(&aws.Config{
		// Region:      aws.String("	ap-northeast-2"), //
		Region:      aws.String(region), //
		Credentials: credentials.NewStaticCredentials(akid, secretkey, ""),
	})

	if err != nil {
		errmsg := jsonErr{503, "failed", "result fail"}
		json.NewEncoder(w).Encode(errmsg)
	}

	var clusters []EKSCluster

	svc := eks.New(sess)
	asSvc := autoscaling.New(sess)
	cls, _ := svc.ListClusters(&eks.ListClustersInput{})

	for _, v := range cls.Clusters {
		ngs, _ := svc.ListNodegroups(&eks.ListNodegroupsInput{
			ClusterName: aws.String(*v),
		})
		var nodegroups []EKSNodegroup
		for _, ng := range ngs.Nodegroups {
			fmt.Println(*ng)
			dng, _ := svc.DescribeNodegroup(&eks.DescribeNodegroupInput{
				ClusterName:   aws.String(*v),
				NodegroupName: aws.String(*ng),
			})
			desiredSize := dng.Nodegroup.ScalingConfig.DesiredSize
			maxSize := dng.Nodegroup.ScalingConfig.MaxSize
			minSize := dng.Nodegroup.ScalingConfig.MinSize
			instanceType := dng.Nodegroup.InstanceTypes[0]
			asgs := dng.Nodegroup.Resources.AutoScalingGroups
			asEKSInstances := make(map[string][]EKSInstance)
			var asgName string
			for _, asg := range asgs {
				instances, _ := asSvc.DescribeAutoScalingInstances(&autoscaling.DescribeAutoScalingInstancesInput{})
				var ints []EKSInstance
				for index, instance := range instances.AutoScalingInstances {
					if *asg.Name == *instance.AutoScalingGroupName {
						fmt.Println(instance)
						fmt.Println(index, *instance.AutoScalingGroupName, *instance.InstanceId)
						ints = append(ints, EKSInstance{*instance.InstanceId})
					}
				}
				asgName = *asg.Name
				asEKSInstances[*asg.Name] = ints
			}
			nodegroups = append(nodegroups, EKSNodegroup{
				*ng,
				*instanceType,
				*desiredSize,
				*maxSize,
				*minSize,
				asgName,
				asEKSInstances[asgName],
			})
		}
		clusters = append(clusters, EKSCluster{*v, nodegroups})
	}
	json.NewEncoder(w).Encode(clusters)
}

// add/remove eks node
func ChangeEKSnode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// Post로 변경
	body := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지
	region := body["region"].(string)
	cluster := body["cluster"].(string)
	nodegroup := body["nodePool"].(string)
	desiredSizeStr := body["desiredCnt"].(string)
	akid := body["accessKey"].(string)
	secretkey := body["secretKey"].(string)

	// http://192.168.0.51:4885/apis/changeeksnode?region=ap-northeast-2&cluster=eks-cluster1&nodegroup=ng-1&nodecount=3
	// region := r.URL.Query().Get("region")
	// cluster := r.URL.Query().Get("cluster")
	// nodegroup := r.URL.Query().Get("nodegroup")
	// desiredSizeStr := r.URL.Query().Get("nodecount")
	// akid := "AKIAJGFO6OXHRN2H6DSA"
	// secretkey := "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"
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(akid, secretkey, ""),
	})

	if err != nil {
		errmsg := jsonErr{503, "failed", "result fail"}
		json.NewEncoder(w).Encode(errmsg)
	}

	svc := eks.New(sess)

	desirecnt, err := strconv.ParseInt(desiredSizeStr, 10, 64)

	// // la := make(map[string]*string)
	// // namelabel := "newlabel01"
	// // la["newlabel01"] = &namelabel

	// labelinput := eks.UpdateLabelsPayload{la["newlabel01"]}

	addResult, err := svc.UpdateNodegroupConfig(&eks.UpdateNodegroupConfigInput{
		ClusterName:   aws.String(cluster), //
		NodegroupName: aws.String(nodegroup),
		// Labels:        &eks.UpdateLabelsPayload{AddOrUpdateLabels: la},
		ScalingConfig: &eks.NodegroupScalingConfig{
			DesiredSize: &desirecnt,
			MaxSize:     &desirecnt,
		},
	})

	if err != nil {
		errmsg := jsonErr{503, "failed", "result fail"}
		json.NewEncoder(w).Encode(errmsg)
	}

	successmsg := jsonErr{200, "success", addResult.String()}
	// fmt.Println(addResult)
	json.NewEncoder(w).Encode(successmsg)

}

func Addec2node(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	data := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지
	node := data["node"].(string)
	cluster := data["cluster"].(string)
	aKey := data["a_key"].(string)
	sKey := data["s_key"].(string)
	result := cloud.AddNode(node, aKey, sKey)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		errmsg := jsonErr{503, "failed", "result fail"}
		json.NewEncoder(w).Encode(errmsg)
	}
	if result.Result != "Could not create instance" {
		go cloud.GetNodeState(&result.InstanceID, node, cluster, aKey, sKey)
	}
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	// start := time.Now()
	ch := make(chan Resultmap)
	token := GetOpenMCPToken()

	var allUrls []string

	clusterurl := "https://" + openmcpURL + "/apis/core.kubefed.io/v1beta1/kubefedclusters?clustername=openmcp" //기존정보
	// clusterurl := "https://" + openmcpURL + "/apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcpclusters/cluster1?clustername=openmcp" //provider, joined 여부

	//https://192.168.0.152:30000/apis/core.kubefed.io/v1beta1/namespaces/kube-federation-system/kubefedclusters/cluster1?clustername=openmcp
	// https://192.168.0.152:30000/apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcpclusters?clustername=openmcp
	// https://192.168.0.152:30000/apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcpclusters/cluster1?clustername=openmcp

	go CallAPI(token, clusterurl, ch)
	clusters := <-ch
	clusterData := clusters.data

	resCluster := DashboardRes{}

	resJoinedClusters := JoinedClusters{}
	resJoinedClusters.Name = "OMCP-Master"

	var clusterlist = make(map[string]Region)
	var clusternames []string
	clusterHealthyCnt := 0
	clusterUnHealthyCnt := 0
	clusterUnknownCnt := 0
	for _, element := range clusterData["items"].([]interface{}) {
		// fmt.Println("element : ", element)
		region := GetStringElement(element, []string{"status", "zones"})
		zone := "Seoul" //todo Zone관련 데이터 필요 (openmcp)
		// region := element.(map[string]interface{})["status"].(map[string]interface{})["zones"].([]interface{})[0].(string)
		// if index > 0 {
		// 	region = "US"
		// }

		clustername := element.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
		// statusReason := element.(map[string]interface{})["status"].(map[string]interface{})["conditions"].([]interface{})[0].(map[string]interface{})["reason"].(string)
		statusReason := GetStringElement(element, []string{"status", "conditions", "reason"})
		// statusType := element.(map[string]interface{})["status"].(map[string]interface{})["conditions"].([]interface{})[0].(map[string]interface{})["type"].(string)
		statusType := GetStringElement(element, []string{"status", "conditions", "type"})
		// statusTF := element.(map[string]interface{})["status"].(map[string]interface{})["conditions"].([]interface{})[0].(map[string]interface{})["status"].(string)
		statusTF := GetStringElement(element, []string{"status", "conditions", "status"})
		clusterStatus := "Healthy"

		if statusReason == "ClusterNotReachable" && statusType == "Offline" && statusTF == "True" {
			clusterStatus = "Unhealthy"
			clusterUnHealthyCnt++
		} else if statusReason == "ClusterReady" && statusType == "Ready" && statusTF == "True" {
			clusterStatus = "Healthy"
			clusterHealthyCnt++
		} else {
			clusterStatus = "Unknown"
			clusterUnknownCnt++
		}

		clusterUrl := "https://" + openmcpURL + "/apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcpclusters/" + clustername + "?clustername=openmcp"
		go CallAPI(token, clusterUrl, ch)
		clusters := <-ch
		clusterData := clusters.data
		provider := GetStringElement(clusterData["spec"], []string{"clusterPlatformType"})
		joinStatus := GetStringElement(clusterData["spec"], []string{"joinStatus"})
		fmt.Println(joinStatus + " : " + provider)
		// fmt.Println("##############", clusterlist)
		// fmt.Println("##############", clusterlist[region])

		clusterlist[region] =
			Region{
				region,
				Attributes{clusterStatus, "", ""},
				append(clusterlist[region].Children, ChildNode{clustername, Attributes{clusterStatus, "", ""}})}

		resJoinedClusters.Children = append(resJoinedClusters.Children, ChildNode{clustername, Attributes{clusterStatus, zone, region}})
		clusternames = append(clusternames, clustername)
	}

	for _, outp := range clusterlist {
		resCluster.Regions = append(resCluster.Regions, outp)
	}

	for _, cluster := range clusternames {
		nodeurl := "https://" + openmcpURL + "/api/v1/nodes?clustername=" + cluster
		allUrls = append(allUrls, nodeurl)
		podurl := "https://" + openmcpURL + "/api/v1/pods?clustername=" + cluster
		allUrls = append(allUrls, podurl)
		projecturl := "https://" + openmcpURL + "/api/v1/namespaces?clustername=" + cluster
		allUrls = append(allUrls, projecturl)
	}

	for _, arg := range allUrls[0:] {
		go CallAPI(token, arg, ch)
	}

	var results = make(map[string]interface{})
	nsCnt := 0
	podCnt := 0
	nodeCnt := 0

	for range allUrls[0:] {
		result := <-ch
		results[result.url] = result.data
	}

	ruuningPodCnt := 0
	failedPodCnt := 0
	unknownPodCnt := 0
	pendingPodCnt := 0
	activeNSCnt := 0
	terminatingNSCnt := 0
	healthyNodeCnt := 0
	unhealthyNodeCnt := 0
	unknownNodeCnt := 0

	for _, result := range results {
		kind := result.(map[string]interface{})["kind"]

		if kind == "NamespaceList" {
			nsCnt = nsCnt + len(result.(map[string]interface{})["items"].([]interface{}))
			for _, element := range result.(map[string]interface{})["items"].([]interface{}) {
				phase := element.(map[string]interface{})["status"].(map[string]interface{})["phase"]
				if phase == "Active" {
					activeNSCnt++
				} else if phase == "Terminating" {
					terminatingNSCnt++
				}
			}
		} else if kind == "PodList" {
			podCnt = podCnt + len(result.(map[string]interface{})["items"].([]interface{}))
			for _, element := range result.(map[string]interface{})["items"].([]interface{}) {
				phase := element.(map[string]interface{})["status"].(map[string]interface{})["phase"]
				if phase == "Running" {
					ruuningPodCnt++
				} else if phase == "Pending" {
					pendingPodCnt++
				} else if phase == "Failed" {
					failedPodCnt++
				} else if phase == "Unknown" {
					unknownPodCnt++
				}
			}

		} else if kind == "NodeList" {
			nodeCnt = nodeCnt + len(result.(map[string]interface{})["items"].([]interface{}))
			for _, element := range result.(map[string]interface{})["items"].([]interface{}) {
				status := element.(map[string]interface{})["status"]
				var healthCheck = make(map[string]string)
				for _, elem := range status.(map[string]interface{})["conditions"].([]interface{}) {
					conType := elem.(map[string]interface{})["type"].(string)
					tf := elem.(map[string]interface{})["status"].(string)
					healthCheck[conType] = tf
				}

				if healthCheck["Ready"] == "True" && (healthCheck["NetworkUnavailable"] == "" || (healthCheck["NetworkUnavailable"] == "" || healthCheck["NetworkUnavailable"] == "False")) && healthCheck["MemoryPressure"] == "False" && healthCheck["DiskPressure"] == "False" && healthCheck["PIDPressure"] == "False" {
					healthyNodeCnt++
				} else {
					if healthCheck["Ready"] == "Unknown" || (healthCheck["NetworkUnavailable"] == "" || healthCheck["NetworkUnavailable"] == "Unknown") || healthCheck["MemoryPressure"] == "Unknown" || healthCheck["DiskPressure"] == "Unknown" || healthCheck["PIDPressure"] == "Unknown" {
						unknownNodeCnt++
					} else {
						unhealthyNodeCnt++
					}
				}
			}
		}
	}

	resCluster.Clusters.ClustersCnt = len(clusternames)
	resCluster.Nodes.NodesCnt = nodeCnt
	resCluster.Pods.PodsCnt = podCnt
	resCluster.Projects.ProjectsCnt = nsCnt
	resCluster.Projects.ProjectsStatus = append(resCluster.Projects.ProjectsStatus, NameVal{"Active", activeNSCnt})
	resCluster.Projects.ProjectsStatus = append(resCluster.Projects.ProjectsStatus, NameVal{"Terminating", terminatingNSCnt})
	resCluster.Pods.PodsStatus = append(resCluster.Pods.PodsStatus, NameVal{"Running", ruuningPodCnt})
	resCluster.Pods.PodsStatus = append(resCluster.Pods.PodsStatus, NameVal{"Pending", pendingPodCnt})
	resCluster.Pods.PodsStatus = append(resCluster.Pods.PodsStatus, NameVal{"Failed", failedPodCnt})
	resCluster.Pods.PodsStatus = append(resCluster.Pods.PodsStatus, NameVal{"Unknown", unknownPodCnt})
	resCluster.Nodes.NodesStatus = append(resCluster.Nodes.NodesStatus, NameVal{"Healthy", healthyNodeCnt})
	resCluster.Nodes.NodesStatus = append(resCluster.Nodes.NodesStatus, NameVal{"Unhealthy", unhealthyNodeCnt})
	resCluster.Nodes.NodesStatus = append(resCluster.Nodes.NodesStatus, NameVal{"Unknown", unknownNodeCnt})
	resCluster.Clusters.ClustersStatus = append(resCluster.Clusters.ClustersStatus, NameVal{"Healthy", clusterHealthyCnt})
	resCluster.Clusters.ClustersStatus = append(resCluster.Clusters.ClustersStatus, NameVal{"Unhealthy", clusterUnHealthyCnt})
	resCluster.Clusters.ClustersStatus = append(resCluster.Clusters.ClustersStatus, NameVal{"Unknown", clusterUnknownCnt})
	resCluster.JoinedClusters = resJoinedClusters
	json.NewEncoder(w).Encode(resCluster)
}

// func WorkloadsDeploymentsOverviewList(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	w.WriteHeader(http.StatusOK)
// 	if err := json.NewEncoder(w).Encode(resource.ListResource()); err != nil {
// 		panic(err)
// 	}

// }

// func WorkloadsPodsOverviewList(w http.ResponseWriter, r *http.Request) {
// 	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
// 	vars := mux.Vars(r)

// 	var client http.Client
// 	resp, err := client.Get("https://" + targetURL + "/seedcontainer/api/v1/clusters/" + vars["clusterName"] + "/daemonsets/list")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode == http.StatusOK {
// 		bodyBytes, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		w.WriteHeader(http.StatusOK)
// 		w.Write(bodyBytes)
// 	}
// }

// func getDeploymentList(w http.ResponseWriter, r *http.Request) {
// 	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
// 	vars := mux.Vars(r)

// 	var client http.Client
// 	resp, err := client.Get("https://" + targetURL + "/seedcontainer/api/v1/clusters/" + vars["clusterName"] + "/deployments/list")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode == http.StatusOK {
// 		bodyBytes, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		w.WriteHeader(http.StatusOK)
// 		w.Write(bodyBytes)
// 	}
// }

// func getDeploymentDetail(w http.ResponseWriter, r *http.Request) {
// 	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
// 	vars := mux.Vars(r)

// 	var client http.Client

// 	callUrl := "https://" + targetURL + "/seedcontainer/api/v1/clusters/" + vars["clusterName"] + "/namespaces/" + vars["namespaceName"] + "/deployments/" + vars["deploymentName"] + "/detail"
// 	//fmt.Print(callUrl)

// 	resp, err := client.Get(callUrl)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode == http.StatusOK {
// 		bodyBytes, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		w.WriteHeader(http.StatusOK)
// 		w.Write(bodyBytes)
// 	}
// }

// func getDeploymentYaml(w http.ResponseWriter, r *http.Request) {
// 	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
// 	vars := mux.Vars(r)

// 	var client http.Client

// 	callUrl := "https://" + targetURL + "/seedcontainer/api/v1/clusters/" + vars["clusterName"] + "/namespaces/" + vars["namespaceName"] + "/deployments/" + vars["deploymentName"] + "/yaml"
// 	//fmt.Print(callUrl)

// 	resp, err := client.Get(callUrl)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode == http.StatusOK {
// 		bodyBytes, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		w.WriteHeader(http.StatusOK)
// 		w.Write(bodyBytes)
// 	}
// }

// func getDeploymentEvent(w http.ResponseWriter, r *http.Request) {
// 	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
// 	vars := mux.Vars(r)

// 	var client http.Client

// 	callUrl := "https://" + targetURL + "/seedcontainer/api/v1/clusters/" + vars["clusterName"] + "/namespaces/" + vars["namespaceName"] + "/deployments/" + vars["deploymentName"] + "/events"
// 	//fmt.Print(callUrl)

// 	resp, err := client.Get(callUrl)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode == http.StatusOK {
// 		bodyBytes, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		w.WriteHeader(http.StatusOK)
// 		w.Write(bodyBytes)
// 	}
// }
