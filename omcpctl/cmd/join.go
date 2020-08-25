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
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"openmcp/openmcp/omcpctl/resource"
	"openmcp/openmcp/util/clusterManager"
	"path/filepath"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
	"strings"

	//"k8s.io/client-go/tools/clientcmd"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"openmcp/openmcp/util"
	//"openmcp/openmcp/util/clusterManager"
	//genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
)

// joinCmd represents the join command
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

openmcpctl join list
openmcpctl join cluster <CLUSTERIP>`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 && args[0] == "cluster" {
			if args[1] == "" {
				fmt.Println("You Must Provide Cluster IP")
			} else {
				joinCluster(args[1])
			}

		} else if len(args) != 0 && args[0] == "list" {
			fmt.Println("[ cluster list (join) ]")
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
			resource.GetClusterList()
		}
	},
}

func moveToUnjoin(memberIP string) {

	c := cobrautil.GetOmcpctlConf("/var/lib/omcpctl/config.yaml")

	util.CmdExec("umount -l /mnt")
	defer util.CmdExec("umount -l /mnt")

	util.CmdExec("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")

	openmcpIP := GetOutboundIP()

	util.CmdExec("mv /mnt/openmcp/" + openmcpIP + "/members/join/" + memberIP + " /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP)

}

func getDiffJoinIP() []string {
	joinErrorClusterIPs := []string{}
	kubeconfig, _ := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	genClient := genericclient.NewForConfigOrDie(kubeconfig)

	c := cobrautil.GetOmcpctlConf("/var/lib/omcpctl/config.yaml")

	util.CmdExec("umount -l /mnt")
	defer util.CmdExec("umount -l /mnt")

	util.CmdExec("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")
	openmcpIP := GetOutboundIP()
	nfsClusterJoinStr, err := util.CmdExec("ls /mnt/openmcp/" + openmcpIP + "/members/join")
	nfsClusterJoinList := strings.Split(nfsClusterJoinStr, "\n")
	nfsClusterJoinList = nfsClusterJoinList[:len(nfsClusterJoinList)-1]
	if err != nil {
		fmt.Println(err)
		return joinErrorClusterIPs
	}

	clusterList := clusterManager.ListKubeFedClusters(genClient, "kube-federation-system")

	for _, nfsJoinCluster := range nfsClusterJoinList {
		//	fmt.Println(nfsClusterJoin)
		find := false
		for _, cluster := range clusterList.Items {
			if strings.Contains(cluster.Spec.APIEndpoint, nfsJoinCluster) {
				find = true
				break
			}
		}
		if !find {
			//clusterIP := Splitter(cluster.Spec.APIEndpoint,"/:")[1]
			joinErrorClusterIPs = append(joinErrorClusterIPs, nfsJoinCluster)
		}

	}

	return joinErrorClusterIPs
}

func joinCluster(memberIP string) {
	c := cobrautil.GetOmcpctlConf("/var/lib/omcpctl/config.yaml")

	util.CmdExec("umount -l /mnt")
	defer util.CmdExec("umount -l /mnt")

	util.CmdExec("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")

	fmt.Println("Cluster Join Start")

	openmcpIP := GetOutboundIP()
	if !fileExists("/mnt/openmcp/" + openmcpIP) {
		fmt.Println("Failed Join List in OpenMCP Master: " + openmcpIP)
		fmt.Println("=> Not Yet Register OpenMCP.")
		fmt.Println("=> First You Must be Input the Next Command in 'OpenMCP Master Server(" + openmcpIP + ")' : omcpctl register openmcp")
		return
	}

	if !fileExists("/mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP) {
		fmt.Println("Failed UnJoin Cluster '" + memberIP + "' in OpenMCP Master: " + openmcpIP)
		fmt.Println("=> '" + memberIP + "' is Not Joinable Cluster in OpenMCP.")
		return
	}

	kc := cobrautil.GetKubeConfig("/mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP + "/config/config")
	context := kc.Contexts[0]
	cluster := kc.Clusters[0]
	user := kc.Users[0]

	kc = cobrautil.GetKubeConfig("/root/.kube/config")
	kc.Clusters = append(kc.Clusters, cluster)
	kc.Contexts = append(kc.Contexts, context)
	kc.Users = append(kc.Users, user)

	//cobrautil.WriteKubeConfig(kc, "/root/.kube/config_2")

	cobrautil.WriteKubeConfig(kc, "/root/.kube/config")
	util.CmdExec("mv /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP + " /mnt/openmcp/" + openmcpIP + "/members/join/" + memberIP)
	util.CmdExec("kubefedctl join " + cluster.Name + " --cluster-context " + cluster.Name + " --host-cluster-context openmcp --v=2")

	installInitCluster(cluster.Name, c.OpenmcpDir)

	fmt.Println("Cluster Join Completed - " + cluster.Name)
}

func installInitCluster(clusterName, openmcpDir string) {
	fmt.Println("Init Module Deployment Start - " + clusterName)
	install_dir := filepath.Join(openmcpDir, "install_openmcp/member")
	initYamls := []string{"custom-metrics-apiserver", "metallb", "metric-collector", "metrics-server", "nginx-ingress-controller"}

	util.CmdExec("kubectl create ns openmcp --context " + clusterName)
	for _, initYaml := range initYamls {
		//fmt.Println("kubectl create -f " + install_dir + "/" + initYaml + " --context " + clusterName)
		util.CmdExec("kubectl create -f " + install_dir + "/" + initYaml + " --context " + clusterName)
	}

	util.CmdExec("chmod 755 " + install_dir + "/vertical-pod-autoscaler/hack/*")
	util.CmdExec(install_dir + "/vertical-pod-autoscaler/hack/vpa-up.sh " + clusterName)
	fmt.Println("Init Module Deployment Finished - " + clusterName)
}

func init() {
	rootCmd.AddCommand(joinCmd)
}
