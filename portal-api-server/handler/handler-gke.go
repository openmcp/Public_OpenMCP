package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	container "google.golang.org/api/container/v1"
	"google.golang.org/api/option"
)

func GetGKEAuth(projectID string, clientEmail string, privateKey string) (*http.Client, context.Context) {
	ctx := context.Background()
	type cred struct {
		AuthType    string `json:"type"`
		PrivateKey  string `json:"private_key"`
		ClientEmail string `json:"client_email"`
		ProjectID   string `json:"project_id"`
	}
	credType := "service_account"
	// projectID = "just-advice-302807"
	// clientEmail = "gkeadmin@just-advice-302807.iam.gserviceaccount.com"
	// privateKey = "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDWdGxXcM+cRb39\nN6fbCBpibF+EIVFkKGjsuuuGJxEoTQIKp2dnl5FlBFKKSa0cSIz4duwgxc5+25KS\neR5cBB6MjSxBC62qK6VeyNUT2KzyIrQfp/zGmxkBVpXFZ13u0JopiwSH5Kvp4vU1\nOJn4wLA3aLs3QMzUC4rXl6IW0yuyMeClooJLFqxjW7ihry2Y0MjMLuSWeHpqCQCK\n0IntRpqhPoKEkWUjonJnQo7Lem5/iqp8rL80vMDPHuDTPLcQt3pI7Ak6z2qk7etm\ng5jkUS1cVU9Xne2jffEMOjTXPrEgozoHWxN0QLwzrA/7vW6zAt3nfOdO9C6wBzh9\n4GgUeTDDAgMBAAECggEAAlWPaFQ+A5bEE/bVyOM0W6Xk/uyDP50rpzKm+vV/O6UQ\nRKAV1rbQ9PyFuXjxKBb8vHzu4lxvfEn/imtEZ/6o0SF9kyesDZIetq1mRFUIwjSb\n0/cMH/fy3w+GNHkvjeM6ClcNBuhM8WVwWH1JOmqT1caPYxvoHta7/XoVCufkLd2q\nqpFcod8LISW3HN7wSgzB5lpDry+Zk8KoXtxn2bAJyRYeky7tkXQbkCwrE10oUkAs\nivgR27wGF0nowoSvs8KwxWME3zW836fVALyF+dGCBlYVtIMvx6T4cu868dI5JANj\nY6U4H3xjB98MQ/zp7uH6w4kj1/cMxvbfAT7jBTiqAQKBgQD9Q/9bEVxPBc+gKEMo\ncXYCJTCT5XsdAgdw/kXHdR463z70sUbLhvHjt/6xwlNCS5j2jkTbN7InLI6xIwY/\nzdfppXsoW4qyEqrgMHjG1af3AlslEA3GLnkLEIx/VM6zoDKlBlI3uz2PMf4wJiFK\nli3X/5tcpYlyc0pCkIJBQ+o2QQKBgQDYxSf8b2/WW87+L3l6/VlbyWMG9aw5RP2v\nitP0cIqoFj/LkD1pJWtJre0Lnlzgz8JJDcRsbrqDZFuIiWnTc8dy8YM1Pv1kz7xZ\nANvpJGEDr5cZjopOoq+w5zfNDrLf/SPB2g6u9/33Ukds3F0++14901b/f7SjHFN2\nH+OPFwMOAwKBgQDPugrird2Rbwm5qexTaqRI5Cnw1ELjKvvhgJzJGNV/ogXn+tM/\nMeKKTSqYr/NMJ+dBKrVtPERh/xjWTwzcHkBegfz+v/6FSexfT0Jwi2NlpMgPIRi7\nGPjsy1kBQxT6nYWMdx/OWEQIhA+hfFTH8V+OjzbliVyvw8H/0LkVQNgEQQKBgBJr\nhn9T9NvxR0CgRiFmX+6FyW1w+OaQ70G4eVRfL9kist8Yba9+p4RGTEtddKUB4o+U\npOlV63F42LJcguqd/wfMcArZRG0JngauJQHFvpyykhNw4l3WQzm0HDDHm/meqCgz\n4GWL2z/l9P3SJ/ZPI+37BHyHnJDzuj/ia9Lf8LmDAoGBAOm92Sp7qFkrwogzIBfp\nU9PtDc2GeiSj7WJctIakuxQ+bSWtOoPq6CPd8OAWmpgZA8SzCfkWMnBQJhB7A6RQ\nZOA50xvE07ybQ397NLkDKAB56zdQ9hDAYpgkzCFWL1AvIouM8OLU48LLIh3KJLxG\nSUwFrPzKIQz4RKj3em+M+iQP\n-----END PRIVATE KEY-----\n"

	credential, err := json.Marshal(&cred{credType, privateKey, clientEmail, projectID})
	// fmt.Println(string(credential))

	// rawdata, err := ioutil.ReadFile("gke.json")
	// if err != nil {
	// 	fmt.Println("11111")
	// 	fmt.Println(err)
	// }

	ts, err := google.CredentialsFromJSON(ctx, credential, container.CloudPlatformScope)

	// ts, err := google.CredentialsFromJSON(ctx, rawdata, container.CloudPlatformScope)

	if err != nil {
		fmt.Println(err)
	}
	client := oauth2.NewClient(ctx, ts.TokenSource)
	return client, ctx
	// svc, err := container.NewService(ctx, option.WithHTTPClient(client))

	// list, err := svc.Projects.Zones.Clusters.List(projectID, "asia-northeast3-a").Do()
	// if err != nil {
	// 	fmt.Println("errrrrr")
	// 	fmt.Println(err)
	// }
	// for _, v := range list.Clusters {
	// 	fmt.Println(v.Name)
	// }
	// clusterName := "c-66lrt"
	// pList, err := svc.Projects.Zones.Clusters.NodePools.List(projectID, "asia-northeast3-a", clusterName).Do()

	// for _, v := range pList.NodePools {
	// 	fmt.Println(v.Name, v.InstanceGroupUrls)
	// }
	// cSvc, err := compute.NewService(ctx, option.WithHTTPClient(client))
	// instanceLists, err := cSvc.Instances.List(projectID, "asia-northeast3-a").Do()
	// var vmNames []string
	// for _, v := range instanceLists.Items {
	// 	fmt.Println(v.Name)
	// 	vmNames = append(vmNames, v.Name)
	// 	for _, d := range v.Metadata.Items {
	// 		fmt.Println(d.Key, *d.Value)

	// 	}
	// }
	// fmt.Println("==============================")
	// fmt.Println(vmNames[0])

	// //vm auto restart option changed but no effected on managed instnace group vm
	// var boolfalse bool
	// boolfalse = false
	// ss := compute.Scheduling{
	// 	AutomaticRestart: &boolfalse,
	// }
	// task, err := cSvc.Instances.SetScheduling(projectID, "asia-northeast3-a", vmNames[0], &ss).Do()
	// json.NewEncoder(w).Encode(task)

	// // vm shutdown but restarted immediately
	// vmoff, err := cSvc.Instances.Stop(projectID, "asia-northeast3-a", vmNames[0]).Do()
	// fmt.Println(vmoff)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// // "add instances to unmanaged instance group"
	// var vmReqList []*compute.InstanceReference
	// vm := compute.InstanceReference{"https://www.googleapis.com/compute/v1/projects/just-advice-302807/zones/asia-northeast3-a/instances/gke-my-first-cluster-1-default-pool-78434536-7xj6", nil, nil}
	// vmReqList = append(vmReqList, &vm)
	// req := compute.InstanceGroupsAddInstancesRequest{vmReqList, nil, nil}
	// task, err := cSvc.InstanceGroups.AddInstances(projectID, "asia-northeast3-a", "gke-my-first-cluster-1-default-pool-78434536-grp", &req).Do()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// json.NewEncoder(w).Encode(instanceLists)
	// json.NewEncoder(w).Encode(task)
}

func GKEChangeNodeCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// http://192.168.0.89:4885/apis/gkechangenodecount?cluster=cluster-1&pool=default-pool&nodecnt=2

	data := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지

	projectID := data["projectId"].(string)
	clientEmail := data["clientEmail"].(string)
	privateKey := data["privateKey"].(string)

	clusterName := data["cluster"].(string)
	nodePoolName := data["nodePool"].(string)
	nodeCount, err := strconv.ParseInt(data["desiredCnt"].(string), 10, 64)

	// projectID := "just-advice-302807"
	// clientEmail := ""
	// privateKey := ""
	// clusterName := r.URL.Query().Get("cluster")
	// nodePoolName := r.URL.Query().Get("pool")
	// desireNodeCnt := r.URL.Query().Get("nodecnt")
	// desireNodeCnt = strings.TrimSpace(desireNodeCnt)
	// nodeCount, err := strconv.ParseInt(desireNodeCnt, 10, 64)

	// fmt.Println(projectID)
	// fmt.Println(clientEmail)
	// fmt.Println(privateKey)
	// fmt.Println(clusterName)
	// fmt.Println(nodePoolName)
	// fmt.Println(nodeCount,err)

	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(nodePoolName, nodeCount)
	client, ctx := GetGKEAuth(projectID, clientEmail, privateKey)
	svc, err := container.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Println(err)
	}

	lists, err := svc.Projects.Zones.Clusters.List(projectID, "-").Do()
	if err != nil {
		fmt.Println(err)
	}
	var zone string
	for _, v := range lists.Clusters {
		if v.Name == clusterName {
			zone = v.Zone
			break
		}
	}
	// fmt.Println(zone, clusterName)
	req := container.SetNodePoolSizeRequest{
		NodeCount: nodeCount,
	}
	// fmt.Println(req)
	task, err := svc.Projects.Zones.Clusters.NodePools.SetSize(projectID, zone, clusterName, nodePoolName, &req).Do()
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(task)
}

func GetGKEClusters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// http://192.168.0.89:4885/apis/getgkeclusters
	// projectID := "just-advice-302807"
	// clientEmail := ""
	// privateKey := ""

	data := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지

	projectID := data["projectId"].(string)
	clientEmail := data["clientEmail"].(string)
	privateKey := data["privateKey"].(string)

	client, ctx := GetGKEAuth(projectID, clientEmail, privateKey)
	svc, err := container.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Println(err)
	}

	lists, err := svc.Projects.Zones.Clusters.List(projectID, "-").Do()
	if err != nil {
		fmt.Println(err)
	}

	// json.NewEncoder(w).Encode(lists)
	var clusters []GKEClusterInfo
	for _, v := range lists.Clusters {
		// fmt.Println(v.Name)
		// fmt.Println(v.CurrentNodeCount)
		// fmt.Println(v.Zone, v.Location)
		var Pools []GKENodePool
		for _, n := range v.NodePools {
			// fmt.Println(n.Name, n.Config.MachineType, n.InitialNodeCount)
			Pool := GKENodePool{n.Name, n.Config.MachineType, strconv.FormatInt(n.InitialNodeCount, 10)}
			Pools = append(Pools, Pool)
		}
		cluster := GKEClusterInfo{v.Name, v.Location, v.Zone, Pools, strconv.FormatInt(v.CurrentNodeCount, 10)}
		clusters = append(clusters, cluster)
	}
	json.NewEncoder(w).Encode(clusters)

	// 	[
	//     {
	//         "clusterName": "cluster-1",
	//         "location": "asia-northeast3-c",
	//         "zone": "asia-northeast3-c",
	//         "nodePools": [
	//             {
	//                 "nodePoolName": "default-pool",
	//                 "machineType": "g1-small",
	//                 "initialNodeCount": "2"
	//             },
	//             {
	//                 "nodePoolName": "pool-1",
	//                 "machineType": "n1-standard-1",
	//                 "initialNodeCount": "1"
	//             }
	//         ],
	//         "nodeCount": "3"
	//     }
	// ]
}

type GKEClusterInfo struct {
	ClusterName string        `json:"clusterName"`
	Location    string        `json:"location"`
	Zone        string        `json:"zone"`
	NodePools   []GKENodePool `json:"nodePools"`
	NodeCount   string        `json:"nodeCount"`
}

type GKENodePool struct {
	Name        string `json:"nodePoolName"`
	MachineType string `json:"machineType"`
	NodeCount   string `json:"initialNodeCount"`
}
