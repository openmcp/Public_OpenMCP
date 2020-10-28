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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"io/ioutil"
	"openmcp/openmcp/omcpctl/apiServerMethod"
	"openmcp/openmcp/omcpctl/resource"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"strings"
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
openmcpctl get odeploy <ODEPLOYNAME> -o yaml

openmcpctl get ohas <OHASNAME>
openmcpctl get ohas <OHASNAME> -o yaml`,

	Run: func(cmd *cobra.Command, args []string) {
		getResource(args)
	},
}


func getMetaInfo(body []byte) (cobrautil.MetaInfo, error){
	var metainfo cobrautil.MetaInfo
	err := yaml.Unmarshal(body, &metainfo)
	return metainfo, err
}

func getResource(args []string) {
	if len(args) == 0 && cobrautil.Option_file == ""{
		fmt.Println("error: Required resource not specified.")
		return
	}

	resourceKinds := ""
	resourceName := ""
	resourceNamespace := ""

	if len(args) >= 1 {
		resourceKinds = args[0]
	}
	if len(args) >= 2 {
		resourceName = args[1]
	}

	c := cobrautil.GetKubeConfig("/root/.kube/config")
	for _, kubecontext := range c.Contexts {
		if cobrautil.Option_allcluster{

		} else if cobrautil.Option_context != "" && cobrautil.Option_context != kubecontext.Name{
			continue
		} else if cobrautil.Option_context == "" && c.CurrentContext != kubecontext.Name{
			continue
		}
		clusterContext := kubecontext.Context.Cluster

		if cobrautil.Option_file != "" {
			filenameList := cobrautil.GetFileNameList()
			for _, filename := range filenameList {
				yamlFile, err := ioutil.ReadFile(filename)
				if err != nil {
					panic(err)
					continue
				}
				metainfo, err := getMetaInfo(yamlFile)
				if err != nil {
					panic(err)
					continue
				}
				if metainfo.Kind == "Status" {
					fmt.Println(metainfo.Message)
					continue
				}

				resourceKind := cobrautil.KindMap[metainfo.Kind]
				resourceName = metainfo.Metadata.Name
				resourceNamespace = metainfo.Metadata.Namespace

				err = getCore(resourceKind, resourceName, resourceNamespace, clusterContext)
				if err != nil {
					continue
				}

			}
		} else {
			resourceKindList := strings.Split(resourceKinds, ",")

			for _, resourceKind := range resourceKindList {
				resourceNamespace = cobrautil.Option_namespace
				err := getCore(resourceKind, resourceName, resourceNamespace, clusterContext)
				if err != nil {
					continue
				}
			}
		}



	}



}
func getCore(resourceKind, resourceName, resourceNamespace, clusterContext string) error{
	LINK := cobrautil.GetLinkParser(resourceKind, resourceName, resourceNamespace, clusterContext)
	fmt.Println(LINK)

	body, err := apiServerMethod.GetAPIServer(LINK)
	if err != nil {
		fmt.Println("error: the server doesn't have a resource type '" + resourceKind + "'")
		return err
	}

	if cobrautil.Option_filetype == "yaml"{
		var prettyYaml map[string]interface{}

		err = yaml.Unmarshal(body, &prettyYaml)
		if err != nil {
			panic(err.Error())
			return err
		}
		yamlResource, _ := yaml.Marshal(prettyYaml)

		fmt.Println(string(yamlResource))

	} else if cobrautil.Option_filetype == "json" {
		var prettyJSON bytes.Buffer

		err = json.Indent(&prettyJSON, body, "", "\t")
		if err != nil {
			panic(err.Error())
			return err
		}
		fmt.Println(string(prettyJSON.Bytes()))

	} else {
		metainfo, err := getMetaInfo(body)
		if err != nil {
			panic(err)
			return err
		}
		if metainfo.Kind == "Status" {
			fmt.Println(metainfo.Message)
			return cobrautil.NewError("")
		}


		fmt.Println("metainfo.Kind : ",  metainfo.Kind)
		fmt.Println("Cluster : ", clusterContext)

		if  metainfo.Kind == "CronJob"{
			resource.PrintCronJob(body)
		} else if  metainfo.Kind == "CronJobList" {
			resource.PrintCronJobList(body)
		} else if  metainfo.Kind == "DaemonSet"{
			resource.PrintDaemonSet(body)
		} else if  metainfo.Kind == "DaemonSetList"{
			resource.PrintDaemonSetList(body)
		} else if  metainfo.Kind == "Deployment"{
			resource.PrintDeployment(body)
		} else if  metainfo.Kind == "DeploymentList" {
			resource.PrintDeploymentList(body)
		} else if  metainfo.Kind == "Job"{
			resource.PrintJob(body)
		} else if  metainfo.Kind == "JobList"{
			resource.PrintJobList(body)
		} else if  metainfo.Kind == "Pod"{
			resource.PrintPod(body)
		} else if  metainfo.Kind == "PodList"{
			resource.PrintPodList(body)
		} else if  metainfo.Kind == "ReplicaSet"{
			resource.PrintReplicaSet(body)
		} else if  metainfo.Kind == "ReplicaSetList"{
			resource.PrintReplicaSetList(body)
		} else if  metainfo.Kind == "ReplicaController"{
			resource.PrintReplicationController(body)
		} else if  metainfo.Kind == "ReplicaControllerList"{
			resource.PrintReplicationControllerList(body)
		} else if  metainfo.Kind == "StatefulSet"{
			resource.PrintStatefulSet(body)
		} else if  metainfo.Kind == "StatefulSetList"{
			resource.PrintStatefulSetList(body)
		} else if  metainfo.Kind == "Endpoints"{
			resource.PrintEndpoints(body)
		} else if  metainfo.Kind == "EndpointsList"{
			resource.PrintEndpointsList(body)
		} else if  metainfo.Kind == "Ingress"{
			resource.PrintIngress(body)
		} else if  metainfo.Kind == "IngressList"{
			resource.PrintIngressList(body)
		} else if  metainfo.Kind == "Service"{
			resource.PrintService(body)
		} else if  metainfo.Kind == "ServiceList"{
			resource.PrintServiceList(body)
		} else if  metainfo.Kind == "ConfigMap"{
			resource.PrintConfigMap(body)
		} else if  metainfo.Kind == "ConfigMapList"{
			resource.PrintConfigMapList(body)
		} else if  metainfo.Kind == "Secret"{
			resource.PrintSecret(body)
		} else if  metainfo.Kind == "SecretList"{
			resource.PrintSecretList(body)
		} else if  metainfo.Kind == "PersistentVolumeClaim"{
			resource.PrintPersistentVolumeClaim(body)
		} else if  metainfo.Kind == "PersistentVolumeClaimList"{
			resource.PrintPersistentVolumeClaimList(body)
		} else if  metainfo.Kind == "StorageClass"{
			resource.PrintStorageClass(body)
		} else if  metainfo.Kind == "StorageClassList"{
			resource.PrintStorageClassList(body)
		} else if  metainfo.Kind == "CustomResourceDefinition"{
			resource.PrintCustomResourceDefinition(body)
		} else if  metainfo.Kind == "CustomResourceDefinitionList"{
			resource.PrintCustomResourceDefinitionList(body)
		} else if  metainfo.Kind == "Event"{
			resource.PrintEvent(body)
		} else if  metainfo.Kind == "EventList"{
			resource.PrintEventList(body)
		} else if  metainfo.Kind == "HorizontalPodAutoscaler"{
			resource.PrintHorizontalPodAutoscaler(body)
		} else if  metainfo.Kind == "HorizontalPodAutoscalerList"{
			resource.PrintHorizontalPodAutoscalerList(body)
		} else if  metainfo.Kind == "VerticalPodAutoscaler"{
			resource.PrintVerticalPodAutoscaler(body)
		} else if  metainfo.Kind == "VerticalPodAutoscalerList"{
			resource.PrintVerticalPodAutoscalerList(body)
		} else if  metainfo.Kind == "PodDisruptionBudget"{
			resource.PrintPodDisruptionBudget(body)
		} else if  metainfo.Kind == "PodDisruptionBudgetList"{
			resource.PrintPodDisruptionBudgetList(body)
		} else if  metainfo.Kind == "APIService"{
			resource.PrintAPIService(body)
		} else if  metainfo.Kind == "APIServiceList"{
			resource.PrintAPIServiceList(body)
		} else if  metainfo.Kind == "ClusterRole"{
			resource.PrintClusterRole(body)
		} else if  metainfo.Kind == "ClusterRoleList"{
			resource.PrintClusterRoleList(body)
		} else if  metainfo.Kind == "ClusterRoleBinding"{
			resource.PrintClusterRoleBinding(body)
		} else if  metainfo.Kind == "ClusterRoleBindingList"{
			resource.PrintClusterRoleBindingList(body)
		} else if  metainfo.Kind == "Namespace"{
			resource.PrintNamespace(body)
		} else if  metainfo.Kind == "NamespaceList"{
			resource.PrintNamespaceList(body)
		} else if  metainfo.Kind == "Node"{
			resource.PrintNode(body)
		} else if  metainfo.Kind == "NodeList"{
			resource.PrintNodeList(body)
		} else if  metainfo.Kind == "PersistentVolume"{
			resource.PrintPersistentVolume(body)
		} else if  metainfo.Kind == "PersistentVolumeList"{
			resource.PrintPersistentVolumeList(body)
		} else if  metainfo.Kind == "ResourceQuota"{
			resource.PrintResourceQuota(body)
		} else if  metainfo.Kind == "ResourceQuotaList"{
			resource.PrintResourceQuotaList(body)
		} else if  metainfo.Kind == "Role"{
			resource.PrintRole(body)
		} else if  metainfo.Kind == "RoleList"{
			resource.PrintRoleList(body)
		} else if  metainfo.Kind == "RoleBinding"{
			resource.PrintRoleBinding(body)
		} else if  metainfo.Kind == "RoleBindingList"{
			resource.PrintRoleBindingList(body)
		} else if  metainfo.Kind == "ServiceAccount"{
			resource.PrintServiceAccount(body)
		} else if  metainfo.Kind == "ServiceAccountList"{
			resource.PrintServiceAccountList(body)
		} else if  metainfo.Kind == "KubeFedCluster"{
			resource.PrintKubeFedCluster(body)
		} else if  metainfo.Kind == "KubeFedClusterList"{
			resource.PrintKubeFedClusterList(body)
		} else if  metainfo.Kind == "OpenMCPDeployment"{
			resource.PrintOpenMCPDeployment(body)
		} else if  metainfo.Kind == "OpenMCPDeploymentList"{
			resource.PrintOpenMCPDeploymentList(body)
		} else if  metainfo.Kind == "OpenMCPService"{
			resource.PrintOpenMCPService(body)
		} else if  metainfo.Kind == "OpenMCPServiceList"{
			resource.PrintOpenMCPServiceList(body)
		} else if  metainfo.Kind == "OpenMCPIngress"{
			resource.PrintOpenMCPIngress(body)
		} else if  metainfo.Kind == "OpenMCPIngressList"{
			resource.PrintOpenMCPIngressList(body)
		} else if  metainfo.Kind == "OpenMCPHybridAutoScaler"{
			resource.PrintOpenMCPHybridAutoScaler(body)
		} else if  metainfo.Kind == "OpenMCPHybridAutoScalerList"{
			resource.PrintOpenMCPHybridAutoScalerList(body)
		} else if  metainfo.Kind == "OpenMCPPolicy"{
			resource.PrintOpenMCPPolicy(body)
		} else if  metainfo.Kind == "OpenMCPPolicyList"{
			resource.PrintOpenMCPPolicyList(body)
		} else if  metainfo.Kind == "OpenMCPConfigMap"{
			resource.PrintOpenMCPConfigMap(body)
		} else if  metainfo.Kind == "OpenMCPConfigMapList"{
			resource.PrintOpenMCPConfigMapList(body)
		} else if  metainfo.Kind == "OpenMCPSecret"{
			resource.PrintOpenMCPSecret(body)
		} else if  metainfo.Kind == "OpenMCPSecretList"{
			resource.PrintOpenMCPSecretList(body)
		} else {
			fmt.Println("error: the server doesn't have a resource type \""+resourceKind+"\"")
			return cobrautil.NewError("")
		}

		fmt.Println()
	}
	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getCmd.Flags().StringVarP(&cobrautil.Option_filetype, "option","o", "", "input a option")
	getCmd.Flags().StringVarP(&cobrautil.Option_context, "context","c","", "input a option")
	getCmd.Flags().StringVarP(&cobrautil.Option_namespace, "namespace","n", "", "input a option")
	getCmd.Flags().StringVarP(&cobrautil.Option_file, "file","f", "", "input a option")
	getCmd.Flags().BoolVarP(&cobrautil.Option_allnamespace,"all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	getCmd.Flags().BoolVarP(&cobrautil.Option_allcluster,"all-clusters", "C", false, "If present, list the requested object(s) across all clusters.")

}
