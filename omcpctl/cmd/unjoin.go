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
	"log"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"openmcp/openmcp/util"
	"openmcp/openmcp/util/clusterManager"
	"path/filepath"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
	"strings"
	"time"
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

openmcpctl unjoin cluster <CLUSTERIP>
openmcpctl unjoin gke-cluster <CLUSTERNAME>
openmcpctl unjoin eks-cluster <CLUSTERNAME>`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 && args[0] == "cluster" {
			if args[1] == "" {
				fmt.Println("You Must Provide Cluster IP")
			} else {
				unjoinCluster(args[1])
			}
		} else if len(args) != 0 && (args[0] == "gke-cluster"|| args[0] == "eks-cluster") {
			if args[1] == "" {
				fmt.Println("You Must Provide Cluster Name")
			} else {
				unjoinCloudCluster(args[1])
			}

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

	util.CmdExec2("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")
	openmcpIP := cobrautil.GetOutboundIP()
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

	util.CmdExec2("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")

	openmcpIP := cobrautil.GetOutboundIP()

	util.CmdExec2("mv /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP + " /mnt/openmcp/" + openmcpIP + "/members/join/" + memberIP)

}

func removeInitCluster(clusterName, openmcpDir string) {
	install_dir := filepath.Join(openmcpDir, "install_openmcp/member")
	initYamls := []string{"custom-metrics-apiserver", "metallb", "metric-collector", "metrics-server", "nginx-ingress-controller"}

	for _, initYaml := range initYamls {
		fmt.Println(initYaml)
		util.CmdExec2("kubectl delete -f " + install_dir + "/" + initYaml + " --context " + clusterName)
	}

	util.CmdExec2("chmod 755 " + install_dir + "/vertical-pod-autoscaler/hack/*")

	util.CmdExec2(install_dir + "/vertical-pod-autoscaler/hack/vpa-down.sh " + clusterName)


}

func unjoinCluster(memberIP string) {
	for {
		lockErr := Lock.TryLock()
		if lockErr != nil {
			fmt.Println("Mount Dir Using Another Works. Wait...")
			time.Sleep(time.Second)
		}
		break
	}

	totalStart := time.Now()
	fmt.Println("***** [Start] Cluster UnJoin Start : '", memberIP, "' *****")

	c := cobrautil.GetOmcpctlConf("/var/lib/omcpctl/config.yaml")

	util.CmdExec("umount -l /mnt")
	defer util.CmdExec("umount -l /mnt")

	util.CmdExec2("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")


	openmcpIP := cobrautil.GetOutboundIP()
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

	kubeconfig, _ := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	genClient := genericclient.NewForConfigOrDie(kubeconfig)
	clusterList := clusterManager.ListKubeFedClusters(genClient, "kube-federation-system")

	checkJoin := 0

	for _, cluster := range clusterList.Items {
		if strings.Contains(cluster.Spec.APIEndpoint, memberIP) {
			checkJoin = 1
			break
		}
	}

	if checkJoin == 0 {
		fmt.Println("ERROR - Fail to find cluster")
		return
	}

	start1 := time.Now()
	fmt.Println("***** [Start] 1. Cluster Config Get *****")

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
	elapsed1 := time.Since(start1)
	log.Printf("Cluster Config Get Time : %s", elapsed1)
	fmt.Println("***** [End] 1. Cluster Config Get ***** ")

	start2 := time.Now()
	fmt.Println("***** [Start] 2. Init Service Remove *****")

	removeInitCluster(target_name, c.OpenmcpDir)

	elapsed2 := time.Since(start2)
	log.Printf("Init Service Remove Time : %s", elapsed2)
	fmt.Println("***** [End] 2. Init Service Remove ***** ")

	start3 := time.Now()
	fmt.Println("***** [Start] 3. Cluster UnJoin *****")

	util.CmdExec2("kubefedctl unjoin " + target_name + " --cluster-context " + target_name + " --host-cluster-context openmcp --v=2")

	elapsed3 := time.Since(start3)
	log.Printf("Cluster Unjoin Time : %s", elapsed3)
	fmt.Println("***** [End] 3. Cluster UnJoin ***** ")


	start4 := time.Now()
	fmt.Println("***** [Start] 4. Cluster Config Remove *****")

	kc.Clusters = append(kc.Clusters[:target_name_index], kc.Clusters[target_name_index+1:]...)
	kc.Contexts = append(kc.Contexts[:target_context_index], kc.Contexts[target_context_index+1:]...)
	kc.Users = append(kc.Users[:target_user_index], kc.Users[target_user_index+1:]...)

	cobrautil.WriteKubeConfig(kc, "/root/.kube/config")

	elapsed4 := time.Since(start4)
	log.Printf("Cluster Config Remove Time : %s", elapsed4)
	fmt.Println("***** [End] 4. Cluster Config Remove ***** ")

	util.CmdExec2("mv /mnt/openmcp/" + openmcpIP + "/members/join/" + memberIP + " /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP)
	Lock.Unlock()

	totalElapsed := time.Since(totalStart)
	log.Printf("Cluster UnJoin Total Elapsed Time : %s", totalElapsed)
	fmt.Println("***** [End] Cluster UnJoin Completed - " + target_name, "*****")

	elapsed := time.Since(totalStart)
	log.Printf("Cluster Join Elapsed Time : %s", elapsed)
}

func unjoinCloudCluster(memberName string) {
	for {
		lockErr := Lock.TryLock()
		if lockErr != nil {
			fmt.Println("Mount Dir Using Another Works. Wait...")
			time.Sleep(time.Second)
		}
		break
	}

	totalStart := time.Now()
	fmt.Println("***** [Start] Cluster UnJoin Start : '", memberName, "' *****")

	c := cobrautil.GetOmcpctlConf("/var/lib/omcpctl/config.yaml")

	util.CmdExec("umount -l /mnt")
	defer util.CmdExec("umount -l /mnt")

	util.CmdExec2("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")
	Lock.Unlock()

	openmcpIP := cobrautil.GetOutboundIP()
	if !fileExists("/mnt/openmcp/" + openmcpIP) {
		fmt.Println("Failed UnJoin Cluster '" + memberName + "' in OpenMCP Master: " + openmcpIP)
		fmt.Println("=> Not Yet Register OpenMCP.")
		fmt.Println("=> First You Must be Input the Next Command in 'OpenMCP Master Server(" + openmcpIP + ")' : ompcpctl register openmcp")

		return
	}

	alreadyJoined, err := cobrautil.CheckAlreadyJoinClusterWithPublicClusterName(memberName, "gke")
	if err != nil {
		fmt.Println("CheckAlreadyJoinCluster Error : ", err)
		return
	}

	if !alreadyJoined {
		return
	}

	start1 := time.Now()
	fmt.Println("***** [Start] 1. Init Service Remove *****")

	removeInitCluster(memberName, c.OpenmcpDir)

	elapsed1 := time.Since(start1)
	log.Printf("Init Service Remove Time : %s", elapsed1)
	fmt.Println("***** [End] 1. Init Service Remove ***** ")

	start2 := time.Now()
	fmt.Println("***** [Start] 2. Cluster UnJoin *****")

	util.CmdExec2("kubefedctl unjoin " + memberName + " --cluster-context " + memberName + " --host-cluster-context openmcp --v=2")

	elapsed2 := time.Since(start2)
	log.Printf("Cluster Unjoin Time : %s", elapsed2)
	fmt.Println("***** [End] 2. Cluster UnJoin ***** ")

	start3 := time.Now()
	fmt.Println("***** [Start] 3. Cluster Config Remove *****")

	kc := cobrautil.GetKubeConfig("/root/.kube/config")

	target_name := ""
	target_user := ""
	var target_name_index int
	var target_context_index int
	var target_user_index int

	for i, cluster := range kc.Clusters {
		if memberName == cluster.Name {
			target_name = memberName
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

	kc.Clusters = append(kc.Clusters[:target_name_index], kc.Clusters[target_name_index+1:]...)
	kc.Contexts = append(kc.Contexts[:target_context_index], kc.Contexts[target_context_index+1:]...)
	kc.Users = append(kc.Users[:target_user_index], kc.Users[target_user_index+1:]...)

	cobrautil.WriteKubeConfig(kc, "/root/.kube/config")

	elapsed3 := time.Since(start3)
	log.Printf("Cluster Config Remove Time : %s", elapsed3)
	fmt.Println("***** [End] 3. Cluster Config Remove ***** ")

	totalElapsed := time.Since(totalStart)
	log.Printf("Cluster UnJoin Total Elapsed Time : %s", totalElapsed)
	fmt.Println("***** [End] Cluster UnJoin Completed - " + memberName, "*****")

	elapsed := time.Since(totalStart)
	log.Printf("Cluster Join Elapsed Time : %s", elapsed)
}

func init() {
	rootCmd.AddCommand(unjoinCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
