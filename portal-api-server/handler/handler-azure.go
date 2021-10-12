package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2020-11-01/containerservice"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
)

func AKSAuthorizer(clientID string, clientSec string, tenantID string) (*autorest.BearerAuthorizer, context.Context, error) {
	// clientID = "1edadbd7-d466-43b1-ad73-15a2ee9080ff"
	// clientSec = "07.Tx2r7GobBf.Suq7quNRhO_642z-p~6a"
	// tenantID = "bc231a1b-ab45-4865-bdba-7724c2893f1c"

	authBaseURL := azure.PublicCloud.ActiveDirectoryEndpoint
	resourceURL := azure.PublicCloud.ResourceManagerEndpoint
	oauthConfig, err := adal.NewOAuthConfig(authBaseURL, tenantID)
	// fmt.Println(clientID)
	// fmt.Println(clientSec)
	// fmt.Println(tenantID)
	token, err := adal.NewServicePrincipalToken(*oauthConfig, clientID, clientSec, resourceURL)
	if err != nil {
		fmt.Println("tokenError")
		fmt.Println(err)
	}

	authorizer := autorest.NewBearerAuthorizer(token)
	ctx := context.Background()

	return authorizer, ctx, err
}

func AKSClusterInfo(authorizer autorest.Authorizer, ctx context.Context, subID string) []ManagedCluster {
	resourceURL := azure.PublicCloud.ResourceManagerEndpoint

	aksClient := containerservice.NewManagedClustersClientWithBaseURI(resourceURL, subID)
	aksClient.Authorizer = authorizer
	vmssClient := compute.NewVirtualMachineScaleSetsClientWithBaseURI(resourceURL, subID)
	vmssClient.Authorizer = authorizer

	var lists []ManagedCluster
	// lll, err := aksClient.ListComplete(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(lll)
	// }

	for list, err := aksClient.ListComplete(ctx); list.NotDone(); err = list.Next() {
		if err != nil {
			fmt.Println("got error while traverising Cluster list: ", err)
		}

		clusters := list.Value()

		aPools := *clusters.AgentPoolProfiles
		ap := make(map[string]AgentPool)

		var poolNames []string
		for _, pool := range aPools {
			poolName := *pool.Name
			poolCount := *pool.Count
			poolNames = append(poolNames, poolName)
			ap[poolName] = AgentPool{poolName, "", poolCount}
		}

		lis := strings.Split(*clusters.ID, "/")
		rgNum := 4
		for index, s := range lis {
			if s == "resourcegroups" {
				rgNum = index + 1
			}
		}
		rg := lis[rgNum]
		nodeRG := *clusters.NodeResourceGroup
		var aplist []AgentPool
		for list, err := vmssClient.ListComplete(ctx, nodeRG); list.NotDone(); err = list.Next() {
			if err != nil {
				fmt.Println("got error while traverising vms list: ", err)
			}
			i := list.Value()
			// fmt.Println(*i.Name)
			poolName := ap[*i.Tags["poolName"]].Name
			poolCount := ap[*i.Tags["poolName"]].Count
			vmssName := *i.Name
			aplist = append(aplist, AgentPool{poolName, vmssName, poolCount})
		}
		lists = append(lists, ManagedCluster{*clusters.Name, rg, nodeRG, aplist, *clusters.Location, *clusters.ProvisioningState})
	}

	return lists
}

// func AKSGetAllResources() []ManagedCluster {
func AKSGetAllResources(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	data := GetJsonBody(r.Body)

	// clientID := ""
	// clientSec := ""
	// tenantID := ""
	// subID := "dc80d3cf-4e1a-4b9a-8785-65c4b739e8d2"
	clientID := data["clientId"].(string)
	clientSec := data["clientSec"].(string)
	tenantID := data["tenantId"].(string)
	subID := data["subId"].(string)

	authorizer, ctx, err := AKSAuthorizer(clientID, clientSec, tenantID)
	if err != nil {
		fmt.Println("AKSAuth failed", err)
	}

	mc := AKSClusterInfo(authorizer, ctx, subID)
	json.NewEncoder(w).Encode(mc)
}

func StopAKSNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// // http://192.168.0.89:4885/apis/aksnodepower?cluster=azure-cluster-1&node=aks-agentpool-17101166-vmss_0
	// clientID := ""
	// clientSec := ""
	// tenantID := ""
	// subID := "dc80d3cf-4e1a-4b9a-8785-65c4b739e8d2"

	// clusterName := r.URL.Query().Get("cluster")
	// vmName := r.URL.Query().Get("node")

	//Post
	body := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지

	clientID := body["clientId"].(string)
	clientSec := body["clientSec"].(string)
	tenantID := body["tenantId"].(string)
	subID := body["subId"].(string)

	clusterName := body["cluster"].(string)
	vmName := body["node"].(string)

	// fmt.Println(clientID)
	// fmt.Println(clientSec)
	// fmt.Println(tenantID)
	// fmt.Println(subID)
	// fmt.Println(clusterName)
	// fmt.Println(vmName)

	authorizer, ctx, err := AKSAuthorizer(clientID, clientSec, tenantID)
	if err != nil {
		fmt.Println("AKSAuth failed", err)
	}

	mc := AKSClusterInfo(authorizer, ctx, subID)
	var clusterData ManagedCluster

	for _, d := range mc {
		if d.Name == clusterName {
			clusterData = d
			break
		}
	}

	nodeRG := clusterData.NodeResourceGrouop
	vmssNames := clusterData.AgentPool
	resourceURL := azure.PublicCloud.ResourceManagerEndpoint
	vmsClient := compute.NewVirtualMachineScaleSetVMsClientWithBaseURI(resourceURL, subID)
	vmsClient.Authorizer = authorizer
	// filterStr := "name eq '" + vmName + "'"
	var targetVmss string = ""
	var tagetVMID string
	for _, d := range vmssNames {
		for list, err := vmsClient.ListComplete(ctx, nodeRG, d.VmssName, "", "", ""); list.NotDone(); err = list.Next() {
			if err != nil {
				fmt.Println("got error while traverising vms list: ", err)
			}
			i := list.Value()
			fmt.Println(*i.Name, *i.InstanceID)
			if *i.Name == vmName {
				targetVmss = d.VmssName
				tagetVMID = *i.InstanceID
				break
			}
		}
		if targetVmss != "" {
			break
		}
	}
	// fmt.Println(vmName)
	// fmt.Println(nodeRG)
	// fmt.Println(vmssNames)
	// fmt.Println(targetVmss)
	// fmt.Println(tagetVMID)

	progress, err := vmsClient.PowerOff(ctx, nodeRG, targetVmss, tagetVMID, nil)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	json.NewEncoder(w).Encode(progress)
}

func StartAKSNode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// // http://192.168.0.89:4885/apis/aksnodepower?cluster=azure-cluster-1&node=aks-agentpool-17101166-vmss_0
	// clientID := ""
	// clientSec := ""
	// tenantID := ""
	// subID := "dc80d3cf-4e1a-4b9a-8785-65c4b739e8d2"

	// clusterName := r.URL.Query().Get("cluster")
	// vmName := r.URL.Query().Get("node")

	//Post
	body := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지

	clientID := body["clientId"].(string)
	clientSec := body["clientSec"].(string)
	tenantID := body["tenantId"].(string)
	subID := body["subId"].(string)

	clusterName := body["cluster"].(string)
	vmName := body["node"].(string)

	// fmt.Println(clientID)
	// fmt.Println(clientSec)
	// fmt.Println(tenantID)
	// fmt.Println(subID)
	// fmt.Println(clusterName)
	// fmt.Println(vmName)

	authorizer, ctx, err := AKSAuthorizer(clientID, clientSec, tenantID)
	if err != nil {
		fmt.Println("AKSAuth failed", err)
	}

	mc := AKSClusterInfo(authorizer, ctx, subID)
	var clusterData ManagedCluster

	for _, d := range mc {
		if d.Name == clusterName {
			clusterData = d
			break
		}
	}

	nodeRG := clusterData.NodeResourceGrouop
	vmssNames := clusterData.AgentPool
	resourceURL := azure.PublicCloud.ResourceManagerEndpoint
	vmsClient := compute.NewVirtualMachineScaleSetVMsClientWithBaseURI(resourceURL, subID)
	vmsClient.Authorizer = authorizer
	// filterStr := "name eq '" + vmName + "'"
	var targetVmss string = ""
	var tagetVMID string
	for _, d := range vmssNames {
		for list, err := vmsClient.ListComplete(ctx, nodeRG, d.VmssName, "", "", ""); list.NotDone(); err = list.Next() {
			if err != nil {
				fmt.Println("got error while traverising vms list: ", err)
			}
			i := list.Value()
			fmt.Println(*i.Name, *i.InstanceID)
			if *i.Name == vmName {
				targetVmss = d.VmssName
				tagetVMID = *i.InstanceID
				break
			}
		}
		if targetVmss != "" {
			break
		}
	}
	// fmt.Println(vmName)
	// fmt.Println(nodeRG)
	// fmt.Println(vmssNames)
	// fmt.Println(targetVmss)
	// fmt.Println(tagetVMID)

	// vmssClient.PowerOff(ctx, config.GroupName(), vmssName, nil, nil)
	// vmssClient.Start(ctx, config.GroupName(), vmssName, nil)
	progress, err := vmsClient.Start(ctx, nodeRG, targetVmss, tagetVMID)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	json.NewEncoder(w).Encode(progress)
}

// aks resource change
func AKSChangeVMSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// http://192.168.0.89:4885/apis/akschangevmss?cluster=azure-cluster-2&pool=agentpool2
	// clientID := ""
	// clientSec := ""
	// tenantID := ""
	// subID := "dc80d3cf-4e1a-4b9a-8785-65c4b739e8d2"
	// skuTierStr := "Standard"
	// skuNameStr := "Standard_B1s"

	// clusterName := r.URL.Query().Get("cluster")
	// targetPool := r.URL.Query().Get("pool")

	data := GetJsonBody(r.Body)

	clientID := data["clientId"].(string)
	clientSec := data["clientSec"].(string)
	tenantID := data["tenantId"].(string)
	subID := data["subId"].(string)
	skuTierStr := data["skuTierStr"].(string)
	skuNameStr := data["skuNameStr"].(string)

	clusterName := data["cluster"].(string)
	targetPool := data["poolName"].(string)

	// fmt.Println(clientID)
	// fmt.Println(clientSec)
	// fmt.Println(tenantID)
	// fmt.Println(subID)
	// fmt.Println(skuTierStr)
	// fmt.Println(skuNameStr)
	// fmt.Println(clusterName)
	// fmt.Println(targetPool)
	authorizer, ctx, err := AKSAuthorizer(clientID, clientSec, tenantID)
	if err != nil {
		fmt.Println("AKSAuth failed", err)
	}

	mc := AKSClusterInfo(authorizer, ctx, subID)
	var clusterData ManagedCluster

	for _, d := range mc {
		if d.Name == clusterName {
			clusterData = d
			break
		}
	}
	resourceURL := azure.PublicCloud.ResourceManagerEndpoint

	aksClient := compute.NewVirtualMachineScaleSetsClientWithBaseURI(resourceURL, subID)
	aksClient.Authorizer = authorizer
	targetVMSS := ""
	nodeRG := clusterData.NodeResourceGrouop
	for _, a := range clusterData.AgentPool {
		if a.Name == targetPool {
			targetVMSS = a.VmssName
		}
	}

	vmss, err := aksClient.Get(ctx, nodeRG, targetVMSS)
	location := vmss.Location
	skuCapa := vmss.Sku.Capacity

	task, err := aksClient.CreateOrUpdate(
		ctx,
		nodeRG,
		targetVMSS,
		compute.VirtualMachineScaleSet{
			Location: location,
			Sku: &compute.Sku{
				Tier:     &skuTierStr,
				Name:     &skuNameStr,
				Capacity: skuCapa,
			},
		},
	)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}
	// fmt.Println(task.Status())
	json.NewEncoder(w).Encode(task)

	// // get available Skus
	// for _, vmss := range vmssNames {
	// 	skus, err := aksClient.ListSkus(ctx, nodeRG, vmss)
	// 	if err != nil {
	// 		json.NewEncoder(w).Encode(err)
	// 	}
	// 	json.NewEncoder(w).Encode(skus.Values())
	// }
}

// add/remove aks node
func AddAKSnode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	// http://192.168.0.89:4885/apis/addaksnode?cluster=azure-cluster-2&pool=agentpool2&nodecnt=2

	data := GetJsonBody(r.Body)
	defer r.Body.Close() // 리소스 누출 방지

	clientID := data["clientId"].(string)
	clientSec := data["clientSec"].(string)
	tenantID := data["tenantId"].(string)
	subID := data["subId"].(string)
	clusterName := data["cluster"].(string)
	targetAgentPoolName := data["nodePool"].(string)
	// nodeCountStr := strconv.FormatFloat(data["desiredCnt"].(float64), 'f', 6, 64)
	nodeCount, _ := strconv.ParseInt(data["desiredCnt"].(string), 10, 64)

	// fmt.Println(clientID)
	// fmt.Println(clientSec)
	// fmt.Println(tenantID)
	// fmt.Println(subID)
	// fmt.Println(clusterName)
	// fmt.Println(targetAgentPoolName)
	// fmt.Println(nodeCount)

	// clientID := ""
	// clientSec := ""
	// tenantID := ""
	// subID := "dc80d3cf-4e1a-4b9a-8785-65c4b739e8d2"

	// clusterName := r.URL.Query().Get("cluster")
	// targetAgentPoolName := r.URL.Query().Get("pool")
	// nodeCountStr := r.URL.Query().Get("nodecnt")
	// nodeCount, err := strconv.ParseInt(nodeCountStr, 10, 64)

	authorizer, ctx, err := AKSAuthorizer(clientID, clientSec, tenantID)
	if err != nil {
		fmt.Println("AKSAuth failed", err)
	}

	mc := AKSClusterInfo(authorizer, ctx, subID)
	var clusterData ManagedCluster

	for _, d := range mc {
		// fmt.Println(clusterName, d.Name)
		if d.Name == clusterName {
			clusterData = d
			break
		}
	}

	resourceURL := azure.PublicCloud.ResourceManagerEndpoint

	aksClient := containerservice.NewManagedClustersClientWithBaseURI(resourceURL, subID)
	aksClient.Authorizer = authorizer
	resourceGroupName := clusterData.ResourceGroup
	resourceName := clusterName
	location := clusterData.Location
	// fmt.Println("==========================")
	// fmt.Println(clusterData.ProvisionState)
	pvstate := clusterData.ProvisionState

	for i := 0; i < 100; i++ {
		// fmt.Println(pvstate)
		if pvstate != "Succeeded" {
			mc := AKSClusterInfo(authorizer, ctx, subID)
			var data ManagedCluster
			for _, d := range mc {

				if d.Name == clusterName {
					data = d
					break
				}
			}
			pvstate = data.ProvisionState
		} else {
			break
		}
		time.Sleep(time.Second * 3)
	}

	res, err := aksClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		resourceName,
		containerservice.ManagedCluster{
			Location: &location,
			ManagedClusterProperties: &containerservice.ManagedClusterProperties{
				AgentPoolProfiles: &[]containerservice.ManagedClusterAgentPoolProfile{
					{
						Count: to.Int32Ptr(int32(nodeCount)),
						Name:  to.StringPtr(targetAgentPoolName),
						Mode:  containerservice.AgentPoolMode("System"),
					},
				},
			},
		},
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Count Change Success")
	json.NewEncoder(w).Encode(res)

	// // get provision state after change config
	// c, err := aksClient.Get(ctx, resourceGroupName, resourceName)
	// fmt.Println(c.AgentPoolProfiles)
	// json.NewEncoder(w).Encode(c.AgentPoolProfiles)

	// c, err := aksClient.ListComplete(ctx)
	// fmt.Println(clusterData)
	// json.NewEncoder(w).Encode(c.Value())

}
