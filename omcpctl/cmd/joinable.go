/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"k8s.io/client-go/tools/clientcmd"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"openmcp/openmcp/util"
	"openmcp/openmcp/util/clusterManager"
	"os"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
	"strings"

	"github.com/spf13/cobra"
)

// joinableCmd represents the joinable command
var joinableCmd = &cobra.Command{
	Use:   "joinable",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 && args[0] == "list" {
			GetJoinableClusterList()
		}
	},
}


func GetJoinableClusterList() {
	c := cobrautil.GetOmcpctlConf("/var/lib/omcpctl/config.yaml")

	util.CmdExec("umount -l /mnt")
	defer util.CmdExec("umount -l /mnt")

	util.CmdExec2("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")
	openmcpIP := GetOutboundIP()
	nfsClusterJoinStr, err := util.CmdExec("ls /mnt/openmcp/" + openmcpIP + "/members/unjoin")
	nfsClusterJoinList := strings.Split(nfsClusterJoinStr, "\n")
	nfsClusterJoinList = nfsClusterJoinList[:len(nfsClusterJoinList)-1]
	if err != nil {
		fmt.Println(err)
	}

	datas := [][]string{}

	for i := range nfsClusterJoinList {
		kc := cobrautil.GetKubeConfig("/mnt/openmcp/" + openmcpIP + "/members/unjoin/" + nfsClusterJoinList[i] + "/config/config")

		data := []string{kc.Clusters[0].Name, "https://"+kc.Clusters[0].Cluster.Server, ""}
		datas = append(datas, data)
	}

	gke_datas := getGKEClusterData()
	for _, gke_data := range gke_datas {
		datas = append(datas, gke_data)
	}

	eks_datas := getEKSClusterData()
	for _, eks_data := range eks_datas {
		datas = append(datas, eks_data)
	}


	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ClusterName", "apiEndpoint", "Platform"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}

func getGKEClusterData() [][]string{

	kubeconfig, _ := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	genClient := genericclient.NewForConfigOrDie(kubeconfig)

	clusterList := clusterManager.ListKubeFedClusters(genClient, "kube-federation-system")

	datas := [][]string{}
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
		clusterName := ss[0]
		masterIP := "https://"+ss[3]
		platform := "gke"
		status := ss[7]

		if status != "RUNNING"{
			continue
		}
		isAlreadyJoined := false
		for _, joinedCluster := range clusterList.Items{
			if clusterName == joinedCluster.Name && strings.Contains(joinedCluster.Spec.APIEndpoint, masterIP){
				isAlreadyJoined = true
				break
			}
		}
		if !isAlreadyJoined {
			data := []string{clusterName, masterIP, platform}
			datas = append(datas, data)
		}



	}


	return datas
}

func getEKSClusterData() [][]string{

	kubeconfig, _ := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	genClient := genericclient.NewForConfigOrDie(kubeconfig)

	clusterList := clusterManager.ListKubeFedClusters(genClient, "kube-federation-system")

	datas := [][]string{}
	s, err := util.CmdExec("aws eks list-clusters")
	if err != nil {
		//fmt.Println(err)
		return datas
	}
	//fmt.Println(s)
	jsonData := make(map[string]interface{})
	err = json.Unmarshal([]byte(s), &jsonData)
	if err != nil {
		//fmt.Println(err)
		return datas
	}
	if _, ok := jsonData["clusters"]; !ok {
		return datas
	}
	eksClusterNamesInteface := jsonData["clusters"].([]interface{})
	for _, clusterNameInteface := range eksClusterNamesInteface{
		clusterName := clusterNameInteface.(string)
		ss, err := util.CmdExec("aws eks describe-cluster --name "+clusterName+" | cat")

		if err != nil {
			fmt.Println(err)
		}
		err = json.Unmarshal([]byte(ss), &jsonData)
		if err != nil {
			fmt.Println(err)
		}
		clusterInfo := jsonData["cluster"].(map[string]interface{})
		//fmt.Println(clusterName)
		apiEndpoint := ""
		if _, ok := clusterInfo["endpoint"]; ok {
			apiEndpoint = clusterInfo["endpoint"].(string)
		} else {
			continue
		}

		apiEndpoint = strings.ToLower(apiEndpoint)
		//fmt.Println(apiEndpoint)
		platform := "eks"

		isAlreadyJoined := false
		for _, joinedCluster := range clusterList.Items{
			if clusterName == joinedCluster.Name && strings.Contains(joinedCluster.Spec.APIEndpoint, apiEndpoint){
				isAlreadyJoined = true
				break
			}
		}
		if !isAlreadyJoined {
			data := []string{clusterName, apiEndpoint, platform}
			datas = append(datas, data)
		}

	}

	return datas
}
func init() {
	rootCmd.AddCommand(joinableCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// joinableCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// joinableCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
