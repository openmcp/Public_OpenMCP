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
	//"admiralty.io/multicluster-controller/pkg/cluster"
	"context"
	"github.com/ghodss/yaml"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	clientV1alpha1 "openmcp/openmcp/openmcp-resource-controller/clientset/v1alpha1"
	"sort"
	"strconv"

	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//ketiv1alpha1 "openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"

	cobrautil "openmcp/openmcp/cobra-test/util"
	"openmcp/openmcp/util"

	"sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"

	"fmt"

	"github.com/olekukonko/tablewriter"

	"openmcp/openmcp/util/clusterManager"
	"os"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"

	//"sigs.k8s.io/controller-runtime/pkg/client"
)



// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

openmcpctl get cluster
openmcpctl get cluster -o yaml

openmcpctl get cluster <CLUSTERNAME>
openmcpctl get cluster <CLUSTERNAME> -o yaml

openmcpctl get node
openmcpctl get node -o yaml
openmcpctl get node --context <CLUUSTERNAME>
openmcpctl get node --context <CLUUSTERNAME> -o yaml

openmcpctl get node <NODENAME>
openmcpctl get node <NODENAME> -o yaml
openmcpctl get node <NODENAME> --context <CLUSTERNAME>
openmcpctl get node <NODENAME> --context <CLUSTERNAME> -o yaml

openmcpctl get odeploy <ODEPLOYNAME>
openmcpctl get odeploy <ODEPLOYNAME> -o yaml`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 && args[0] == "cluster" {
			if len(args) == 1 {
				getClusterList()
			} else if len(args) == 2 {
				clusterName := args[1]
				getCluster(clusterName)
			}
		} else if len(args) != 0 && args[0] == "node" {
			if len(args) == 1 {
				getNodeList()
			} else if len(args) == 2 {
				nodeName := args[1]
				getNode(nodeName)
			}
		} else if  len(args) != 0 && (args[0] == "openmcpdeployment" ||  args[0] == "odeploy"){
			if len(args) == 1 {
				getOpenmcpDeploymentList()
			} else if len(args) == 2 {
				odName := args[1]
				getOpenmcpDeployment(odName)
			}
		} else if  len(args) != 0 && (args[0] == "openmcpdeployment" ||  args[0] =="osvc"){

		} else if  len(args) != 0 && (args[0] == "openmcpdeployment" ||  args[0] == "oing"){

		} else if  len(args) != 0 && (args[0] == "openmcpdeployment" ||  args[0] == "ohas"){

		} else if  len(args) != 0 && (args[0] == "openmcpdeployment" ||  args[0] == "opol"){

		} else {
			cmdStr := "kubectl get"
			for i := 0; i < len(args); i++ {
				cmdStr = cmdStr + " " + args[i]
			}
			if option_file != "" {
				cmdStr = cmdStr + " -f " + option_file
			}
			if option_namespace != "" {
				cmdStr = cmdStr + " -f " + option_namespace
			}
			if option_allnamespace != false {
				cmdStr = cmdStr + " -A"
			}
			if option_context != ""{
				cmdStr = cmdStr + " --context " + option_context
			}
			util.CmdExec2(cmdStr)
		}
	},
}
func getOpenmcpDeployment(odName string) {
	kubeconfig, _ := cobrautil.BuildConfigFromFlags("openmcp", "/root/.kube/config")
	crdClient, _ := clientV1alpha1.NewForConfig(kubeconfig)

	ns := ""
	if option_allnamespace {
		ns = metav1.NamespaceAll
	} else {
		if option_namespace != "" {
			ns = option_namespace
		} else {
			ns = "default"
		}
	}
	od, err := crdClient.OpenMCPDeployment(ns).Get(odName, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		fmt.Println("No resources found in " + ns + " namespace.")
		return
	} else if err != nil {
		fmt.Println(err)
		return
	}

	if option_filetype == "yaml"{
		data, err := yaml.Marshal(&od)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		fmt.Println(string(data))
		return
	}

	datas := [][]string{}
	namespace := od.Namespace
	name := od.Name
	cluster := ""

	keys := make([]string, 0, len(od.Status.ClusterMaps))
	for k := range od.Status.ClusterMaps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys{
		if od.Status.ClusterMaps[k] >= 1 {
			cluster += k+":"+strconv.Itoa(int(od.Status.ClusterMaps[k]))+" "
		}
	}
	age := ""
	data := []string{namespace, name, cluster, age}
	datas = append(datas, data)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "CLUSTER", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()

}
func getOpenmcpDeploymentList() {
	kubeconfig, _ := cobrautil.BuildConfigFromFlags("openmcp", "/root/.kube/config")
	crdClient, _ := clientV1alpha1.NewForConfig(kubeconfig)

	ns := ""
	if option_allnamespace {
		ns = metav1.NamespaceAll
	} else {
		if option_namespace != ""{
			ns = option_namespace
		} else {
			ns = "default"
		}
	}
	odList, err := crdClient.OpenMCPDeployment(ns).List(metav1.ListOptions{})

	if err != nil && errors.IsNotFound(err) {
		fmt.Println("No Resource exist")
	}
	if len(odList.Items) == 0 {
		fmt.Println("No resources found in "+ns+" namespace.")
		return
	}
	if option_filetype == "yaml"{
		data, err := yaml.Marshal(&odList)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
		fmt.Println(string(data))
		return
	}

	datas := [][]string{}
	for _, od := range odList.Items {
		namespace := od.Namespace
		name := od.Name
		cluster := ""

		keys := make([]string, 0, len(od.Status.ClusterMaps))
		for k := range od.Status.ClusterMaps {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys{
			if od.Status.ClusterMaps[k] >= 1 {
				cluster += k+":"+strconv.Itoa(int(od.Status.ClusterMaps[k]))+" "
			}
		}
		age := ""
		data := []string{namespace, name, cluster, age}
		datas = append(datas, data)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "CLUSTER", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()


}
func getNode(nodeName string) {
	c := cobrautil.GetKubeConfig("/root/.kube/config")

	for _, kubecontext := range c.Contexts {
		if option_context != "" && option_context != kubecontext.Name{
			continue
		}
		kubeconfig, _ := cobrautil.BuildConfigFromFlags(kubecontext.Name, "/root/.kube/config")
		genClient := genericclient.NewForConfigOrDie(kubeconfig)

		node := &corev1.Node{}
		err := genClient.Get(context.Background(), node,"default", nodeName)
		if err != nil && errors.IsNotFound(err){
			continue
		}

		if option_filetype == ""{
			datas := [][]string{}
			nodeName := node.Name
			nodeStatus := ""
			for _, cond := range node.Status.Conditions {
				if cond.Type == "Ready"{
					nodeStatus = string(cond.Status)
					break
				}
			}
			nodeRegion := node.Labels["failure-domain.beta.kubernetes.io/region"]
			nodeZone := node.Labels["failure-domain.beta.kubernetes.io/zone"]
			nodeAddress := node.Status.Addresses[0].Address
			nodeMaster := "false"
			if _, ok := node.Labels["node-role.kubernetes.io/master"]; ok {
				nodeMaster = "true"
			}
			nodeSchedule := "Yes"
			for _, taint := range node.Spec.Taints{
				if taint.Effect == "NoSchedule" {
					nodeSchedule ="No"
				}
			}

			nodeClusterName := kubecontext.Context.Cluster

			data := []string{nodeName, nodeStatus, nodeRegion, nodeZone, nodeAddress, nodeMaster, nodeSchedule, nodeClusterName}

			datas = append(datas, data)
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"NodeName", "Status", "Region", "Zones", "address", "Master", "Schedule", "ClusterName", "AGE"})
			table.SetBorder(false)
			table.AppendBulk(datas)
			table.Render()

		} else if option_filetype == "yaml"{
			if option_context == "" {
				fmt.Println("The Command Must Need Option '--context option_context'")
			} else {
				res, _ := util.CmdExec("kubectl get node " + nodeName + " --context "+ option_context +" -o yaml")
				fmt.Println(res)
			}

		}

	}

}

func getNodeList() {
	c := cobrautil.GetKubeConfig("/root/.kube/config")
	datas := [][]string{}

	for _, kubecontext := range c.Contexts {
		if option_context != "" && option_context != kubecontext.Name{
			continue
		}
		kubeconfig, _ := cobrautil.BuildConfigFromFlags(kubecontext.Name, "/root/.kube/config")
		genClient := genericclient.NewForConfigOrDie(kubeconfig)

		nodeList := &corev1.NodeList{}
		err := genClient.List(context.Background(), nodeList,"default")
		if err != nil && errors.IsNotFound(err){
			continue
		}
		if option_filetype == ""{
			for _, node := range nodeList.Items {
				nodeName := node.Name
				nodeStatus := ""
				for _, cond := range node.Status.Conditions {
					if cond.Type == "Ready"{
						nodeStatus = string(cond.Status)
						break
					}
				}
				nodeRegion := node.Labels["failure-domain.beta.kubernetes.io/region"]
				nodeZone := node.Labels["failure-domain.beta.kubernetes.io/zone"]
				nodeAddress := node.Status.Addresses[0].Address
				nodeMaster := "false"
				if _, ok := node.Labels["node-role.kubernetes.io/master"]; ok {
					nodeMaster = "true"
				}
				nodeSchedule := "Yes"
				for _, taint := range node.Spec.Taints{
					if taint.Effect == "NoSchedule" {
						nodeSchedule ="No"
					}
				}

				nodeClusterName := kubecontext.Context.Cluster

				data := []string{nodeName, nodeStatus, nodeRegion, nodeZone, nodeAddress, nodeMaster, nodeSchedule, nodeClusterName}
				datas = append(datas, data)

			}
		} else if option_filetype == "yaml"{
			if option_context == "" {
				fmt.Println("The Command Must Need Option '--context option_context'")
			} else {
				res, _ := util.CmdExec("kubectl get node --context "+ option_context +" -o yaml")
				fmt.Println(res)
			}
		}



	}
	if option_filetype == ""{
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"NodeName", "Status", "Region", "Zones", "address", "Master", "Schedule", "ClusterName", "AGE"})
		table.SetBorder(false)
		table.AppendBulk(datas)
		table.Render()
	}

}

func getCluster(clusterName string){
	kubeconfig, _ := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	genClient := genericclient.NewForConfigOrDie(kubeconfig)

	kubeFedCluster := &v1beta1.KubeFedCluster{}
	genClient.Get(context.Background(), kubeFedCluster,"kube-federation-system", clusterName)

	if option_filetype == ""{
		datas := [][]string{}
		data := []string{kubeFedCluster.Name, string(kubeFedCluster.Status.Conditions[0].Status), *kubeFedCluster.Status.Region, strings.Join(kubeFedCluster.Status.Zones, ","), kubeFedCluster.Spec.APIEndpoint, kubeFedCluster.GenerateName }

		datas = append(datas, data)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ClusterName", "Status", "Region", "Zones", "apiEndpoint", "AGE"})
		table.SetBorder(false)
		table.AppendBulk(datas)
		table.Render()

	} else if option_filetype == "yaml"{
		res, _ := util.CmdExec("kubectl get kubefedclusters " + clusterName + " -n kube-federation-system -o yaml")
		fmt.Println(res)
	}
}

func getClusterList(){
	kubeconfig, _ := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	genClient := genericclient.NewForConfigOrDie(kubeconfig)

	clusterList := clusterManager.ListKubeFedClusters(genClient, "kube-federation-system")

	if option_filetype == ""{
		datas := [][]string{}
		for _, Cluster := range clusterList.Items {
			data := []string{Cluster.Name, string(Cluster.Status.Conditions[0].Status), *Cluster.Status.Region, strings.Join(Cluster.Status.Zones, ","), Cluster.Spec.APIEndpoint, Cluster.GenerateName }
			datas = append(datas, data)

		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ClusterName", "Status", "Region", "Zones", "apiEndpoint", "AGE"})
		table.SetBorder(false)
		table.AppendBulk(datas)
		table.Render()

	} else if option_filetype == "yaml"{
		res, _ := util.CmdExec("kubectl get kubefedclusters -n kube-federation-system -o yaml")
		fmt.Println(res)
	}
}

func init() {
	rootCmd.AddCommand(getCmd)


	getCmd.Flags().StringVarP(&option_filetype, "option","o", "", "input a option")
	getCmd.Flags().StringVarP(&option_context, "context","c", "", "input a option")
	getCmd.Flags().StringVarP(&option_namespace, "namespace","n", "", "input a option")
	getCmd.Flags().StringVarP(&option_file, "file","f", "", "input a option")
	getCmd.Flags().BoolVarP(&option_allnamespace,"all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")


	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
