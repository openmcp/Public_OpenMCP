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
	credential, err := json.Marshal(&cred{credType, privateKey, clientEmail, projectID})

	ts, err := google.CredentialsFromJSON(ctx, credential, container.CloudPlatformScope)

	if err != nil {
		fmt.Println(err)
	}
	client := oauth2.NewClient(ctx, ts.TokenSource)
	return client, ctx
}

func GKEChangeNodeCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	data := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지

	projectID := data["projectId"].(string)
	clientEmail := data["clientEmail"].(string)
	privateKey := data["privateKey"].(string)

	clusterName := data["cluster"].(string)
	nodePoolName := data["nodePool"].(string)
	nodeCount, err := strconv.ParseInt(data["desiredCnt"].(string), 10, 64)

	if err != nil {
		fmt.Println(err)
	}

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
	req := container.SetNodePoolSizeRequest{
		NodeCount: nodeCount,
	}
	task, err := svc.Projects.Zones.Clusters.NodePools.SetSize(projectID, zone, clusterName, nodePoolName, &req).Do()
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(task)
}

func GetGKEClusters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

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

	var clusters []GKEClusterInfo
	for _, v := range lists.Clusters {
		var Pools []GKENodePool
		for _, n := range v.NodePools {
			Pool := GKENodePool{n.Name, n.Config.MachineType, strconv.FormatInt(n.InitialNodeCount, 10)}
			Pools = append(Pools, Pool)
		}
		cluster := GKEClusterInfo{v.Name, v.Location, v.Zone, Pools, strconv.FormatInt(v.CurrentNodeCount, 10)}
		clusters = append(clusters, cluster)
	}
	json.NewEncoder(w).Encode(clusters)
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
