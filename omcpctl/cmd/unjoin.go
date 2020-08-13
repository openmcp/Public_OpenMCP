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
	"fmt"
	"github.com/olekukonko/tablewriter"
	"k8s.io/client-go/tools/clientcmd"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"openmcp/openmcp/util"
	"openmcp/openmcp/util/clusterManager"
	"os"
	"path/filepath"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
	"strings"

	"github.com/spf13/cobra"
)

// unjoinCmd represents the unjoin command
var unjoinCmd = &cobra.Command{
	Use:   "unjoin",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

openmcpctl unjoin list
openmcpctl unjoin cluster <CLUSTERIP>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 && args[0] == "cluster" {
			if args[1] == "" {
				fmt.Println("You Must Provide Cluster IP")
			} else {
				unjoinCluster(args[1])
			}
		} else if len(args) != 0 && args[0] == "list" {
			fmt.Println("[ cluster list (unjoin) ]")
			//실제로 조인되어 있지 않은 클러스터 정보가 join 디렉토리에 들어가 있는 경우
			// => unjoin
			clusterIps := getDiffJoinIP()
			for _, clusterIp := range clusterIps {
				fmt.Println("Not Correct - ", clusterIp)
				moveToUnjoin(clusterIp)
				unjoinCluster(clusterIp)
			}
			//실제로 조인된 클러스터 정보가 unjoin 디렉토리에 들어가 있는 경우
			// => join
			clusterIps = getDiffUnjoinIP()
			for _, clusterIp := range clusterIps {
				fmt.Println("Not Correct - ", clusterIp)
				moveToJoin(clusterIp)
				joinCluster(clusterIp)
			}
			GetUnjoinClusterList()
		}
	},
}

func getDiffUnjoinIP() []string {
	unjoinErrorClusterIPs := []string{}
	kubeconfig, _ := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	genClient := genericclient.NewForConfigOrDie(kubeconfig)

	c := cobrautil.GetOmcpctlConf("/var/lib/omcpctl/config.yaml")

	util.CmdExec("umount -l /mnt")
	defer util.CmdExec("umount -l /mnt")

	util.CmdExec("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")
	openmcpIP := GetOutboundIP()
	nfsClusterUnjoinStr, err := util.CmdExec("ls /mnt/openmcp/" + openmcpIP + "/members/unjoin")
	nfsClusterUnjoinList := strings.Split(nfsClusterUnjoinStr, "\n")
	nfsClusterUnjoinList = nfsClusterUnjoinList[:len(nfsClusterUnjoinList)-1]
	if err != nil {
		fmt.Println(err)
		return unjoinErrorClusterIPs
	}

	clusterList := clusterManager.ListKubeFedClusters(genClient, "kube-federation-system")

	for _, nfsUnjoinCluster := range nfsClusterUnjoinList {
		for _, cluster := range clusterList.Items {
			if strings.Contains(cluster.Spec.APIEndpoint, nfsUnjoinCluster) {
				unjoinErrorClusterIPs = append(unjoinErrorClusterIPs, nfsUnjoinCluster)
				break
			}
		}
	}

	return unjoinErrorClusterIPs
}

func moveToJoin(memberIP string) {

	c := cobrautil.GetOmcpctlConf("/var/lib/omcpctl/config.yaml")

	util.CmdExec("umount -l /mnt")
	defer util.CmdExec("umount -l /mnt")

	util.CmdExec("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")

	openmcpIP := GetOutboundIP()

	util.CmdExec("mv /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP + " /mnt/openmcp/" + openmcpIP + "/members/join/" + memberIP)

}

func GetUnjoinClusterList() {
	c := cobrautil.GetOmcpctlConf("/var/lib/omcpctl/config.yaml")

	util.CmdExec("umount -l /mnt")
	defer util.CmdExec("umount -l /mnt")

	util.CmdExec("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")
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

		data := []string{kc.Clusters[0].Name, kc.Clusters[0].Cluster.Server}
		datas = append(datas, data)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ClusterName", "apiEndpoint"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}

func removeInitCluster(clusterName, openmcpDir string) {
	install_dir := filepath.Join(openmcpDir, "install_openmcp/member")
	initYamls := []string{"custom-metrics-apiserver", "metallb", "metric-collector", "metrics-server", "nginx-ingress-controller"}

	for _, initYaml := range initYamls {
		util.CmdExec("kubectl delete -f " + install_dir + "/" + initYaml + " --context " + clusterName)
	}

	util.CmdExec("chmod 755 " + install_dir + "/vertical-pod-autoscaler/hack/*")
	util.CmdExec(install_dir + "/vertical-pod-autoscaler/hack/vpa-down.sh " + clusterName)
	util.CmdExec("kubectl delete ns openmcp --context " + clusterName)
}

func unjoinCluster(memberIP string) {
	c := cobrautil.GetOmcpctlConf("/var/lib/omcpctl/config.yaml")

	util.CmdExec("umount -l /mnt")
	defer util.CmdExec("umount -l /mnt")

	util.CmdExec("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")

	fmt.Println("Cluster UnJoin Start")

	openmcpIP := GetOutboundIP()
	if !fileExists("/mnt/openmcp/" + openmcpIP) {
		fmt.Println("Failed UnJoin Cluster '" + memberIP + "' in OpenMCP Master: " + openmcpIP)
		fmt.Println("=> Not Yet Register OpenMCP.")
		fmt.Println("=> First You Must be Input the Next Command in 'OpenMCP Master Server(" + openmcpIP + ")' : ompcpctl register openmcp")

		return
	}

	if !fileExists("/mnt/openmcp/" + openmcpIP + "/members/join/" + memberIP) {
		fmt.Println("Failed UnJoin Cluster '" + memberIP + "' in OpenMCP Master: " + openmcpIP)
		fmt.Println("=> '" + memberIP + "' is Not Joined Cluster in OpenMCP.")

		return
	}

	kc := cobrautil.GetKubeConfig("/root/.kube/config")

	target_name := ""
	target_user := ""
	var target_name_index int
	var target_context_index int
	var target_user_index int

	for i, cluster := range kc.Clusters {
		if strings.Contains(cluster.Cluster.Server, memberIP) {
			target_name = cluster.Name
			target_name_index = i
			break
		}
	}
	for j, context := range kc.Contexts {
		if target_name == context.Context.Cluster {
			target_user = context.Context.User
			target_context_index = j
			break
		}
	}
	for k, user := range kc.Users {
		if target_user == user.Name {
			target_user_index = k
			break
		}
	}

	removeInitCluster(target_name, c.OpenmcpDir)

	util.CmdExec("kubefedctl unjoin " + target_name + " --cluster-context " + target_name + " --host-cluster-context openmcp --v=2")
	kc.Clusters = append(kc.Clusters[:target_name_index], kc.Clusters[target_name_index+1:]...)
	kc.Contexts = append(kc.Contexts[:target_context_index], kc.Contexts[target_context_index+1:]...)
	kc.Users = append(kc.Users[:target_user_index], kc.Users[target_user_index+1:]...)

	cobrautil.WriteKubeConfig(kc, "/root/.kube/config")
	util.CmdExec("mv /mnt/openmcp/" + openmcpIP + "/members/join/" + memberIP + " /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP)

	fmt.Println("Cluster Unjoin Completed - " + target_name)
}
func init() {
	rootCmd.AddCommand(unjoinCmd)
}
