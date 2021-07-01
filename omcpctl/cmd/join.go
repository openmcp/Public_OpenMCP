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
	"context"
	"fmt"
	"log"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"openmcp/openmcp/util"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
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

omcpctl join list
omcpctl join cluster <CLUSTERIP>
omcpctl join gke-cluster <CLUSTERNAME>
omcpctl join eks-cluster <CLUSTERNAME>
omcpctl join aks-cluster <CLUSTERNAME> <RESOURCEGROUPNAME>`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Run 'omcpctl join --help' to view all commands")
		}

		if len(args) != 0 && args[0] == "list" {
			listCluster()
		} else if len(args) != 0 && args[0] == "cluster" {
			if args[1] == "" {
				fmt.Println("Run 'omcpctl join --help' to view all commands")
			} else {
				joinCluster(args[1])
			}
		} else if len(args) != 0 && args[0] == "gke-cluster" {
			if args[1] == "" {
				fmt.Println("Run 'omcpctl join --help' to view all commands")
			} else {
				joinGKECluster(args[1])
			}
		} else if len(args) != 0 && args[0] == "eks-cluster" {
			if args[1] == "" {
				fmt.Println("Run 'omcpctl join --help' to view all commands")
			} else {
				joinEKSCluster(args[1])
			}

		} else if len(args) != 0 && args[0] == "aks-cluster" {
			if args[1] == "" {
				fmt.Println("Run 'omcpctl join --help' to view all commands")
			} else {
				joinAKSCluster(args[1], args[2])
			}

		}
	},
}

/*func moveToUnjoin(memberIP string) {

	c := cobrautil.GetOmcpctlConf("/var/lib/omcpctl/config.yaml")

	util.CmdExec("umount -l /mnt")
	defer util.CmdExec("umount -l /mnt")

	util.CmdExec2("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")

	openmcpIP := cobrautil.GetOutboundIP()

	util.CmdExec2("mv /mnt/openmcp/" + openmcpIP + "/members/join/" + memberIP + " /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP)

}

func getDiffJoinIP() []string {
	joinErrorClusterIPs := []string{}
	kubeconfig, _ := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	genClient := genericclient.NewForConfigOrDie(kubeconfig)

	c := cobrautil.GetOmcpctlConf("/var/lib/omcpctl/config.yaml")

	util.CmdExec("umount -l /mnt")
	defer util.CmdExec("umount -l /mnt")

	util.CmdExec2("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")
	openmcpIP := cobrautil.GetOutboundIP()
	nfsClusterJoinStr, err := util.CmdExec("ls /mnt/openmcp/" + openmcpIP + "/members/join")
	nfsClusterJoinList := strings.Split(nfsClusterJoinStr, "\n")
	nfsClusterJoinList = nfsClusterJoinList[:len(nfsClusterJoinList)-1]
	if err != nil {
		fmt.Println(err)
		return joinErrorClusterIPs
	}

	clusterList := clusterManager.ListKubeFedClusters(genClient, "kube-federation-system")

	for _, nfsJoinCluster := range nfsClusterJoinList {

		find := false
		for _, cluster := range clusterList.Items {
			if strings.Contains(cluster.Spec.APIEndpoint, nfsJoinCluster) {
				find = true
				break
			}
		}
		if !find {
			joinErrorClusterIPs = append(joinErrorClusterIPs, nfsJoinCluster)
		}

	}

	return joinErrorClusterIPs
}*/
func listCluster() {
	util.CmdExec2("omcpctl get cluster -A")
}
func joinCluster(memberIP string) {
	for {
		lockErr := Lock.TryLock()
		if lockErr != nil {
			fmt.Println("Mount Dir Using Another Works. Wait...")
			time.Sleep(time.Second)
		} else {
			break
		}

	}

	totalStart := time.Now()
	fmt.Println("***** [Start] Cluster Join Start : '", memberIP, "' *****")

	c := cobrautil.GetOmcpctlConf("/var/lib/omcpctl/config.yaml")

	util.CmdExec("umount -l /mnt")
	defer util.CmdExec("umount -l /mnt")

	util.CmdExec2("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")

	//openmcpIP := cobrautil.GetOutboundIP()
	openmcpIP := cobrautil.GetEndpointIP()
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

	alreadyJoined, err := cobrautil.CheckAlreadyJoinClusterWithIP(memberIP)
	if err != nil {
		fmt.Println("CheckAlreadyJoinClusterWithIP Error : ", err)
		return
	}

	if alreadyJoined {
		return
	}

	start1 := time.Now()
	fmt.Println("***** [Start] 1. Cluster Config Merge *****")

	kc := cobrautil.GetKubeConfig("/mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP + "/config/config")
	context := kc.Contexts[0]
	cluster := kc.Clusters[0]
	user := kc.Users[0]

	kc = cobrautil.GetKubeConfig("/root/.kube/config")
	kc.Clusters = append(kc.Clusters, cluster)
	kc.Contexts = append(kc.Contexts, context)
	kc.Users = append(kc.Users, user)

	cobrautil.WriteKubeConfig(kc, "/root/.kube/config")

	// namespace terminating stuck force delete
	util.CmdExec2("kubectl get namespace kube-federation-system --context " + cluster.Name + " -o json |jq '.spec = {\"finalizers\":[]}' >temp.json")
	util.CmdExec2("kubectl replace --raw \"/api/v1/namespaces/kube-federation-system/finalize\" -f ./temp.json --context " + cluster.Name)
	util.CmdExec2("rm temp.json")

	elapsed1 := time.Since(start1)
	log.Printf("Cluster Config Merge Time : %s", elapsed1)
	fmt.Println("***** [End] 1. Cluster Config Merge ***** ")

	start2 := time.Now()
	fmt.Println("***** [Start] 2. Cluster Join *****")
	util.CmdExec2("mv /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP + " /mnt/openmcp/" + openmcpIP + "/members/join/" + memberIP)
	util.CmdExec2("kubefedctl join " + cluster.Name + " --cluster-context " + cluster.Name + " --host-cluster-context openmcp --v=2")

	Lock.Unlock()

	elapsed2 := time.Since(start2)
	log.Printf("Cluster Join Time : %s", elapsed2)
	fmt.Println("***** [End] 2. Cluster Join ***** ")

	start3 := time.Now()
	fmt.Println("***** [Start] 3. Init Service Deployments *****")

	installInitCluster(cluster.Name, c.OpenmcpDir, "coredns")

	elapsed3 := time.Since(start3)
	log.Printf("Init Service Deployments Time : %s", elapsed3)
	fmt.Println("***** [End] 3. Init Service Deployments ***** ")

	totalElapsed := time.Since(totalStart)
	log.Printf("Cluster Join Total Elapsed Time : %s", totalElapsed)
	fmt.Println("***** [End] Cluster Join Completed - "+cluster.Name, "*****")
}

func joinGKECluster(memberName string) {

	for {
		lockErr := Lock.TryLock()
		if lockErr != nil {
			fmt.Println("Mount Dir Using Another Works. Wait...")
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	fmt.Println("gke cluster name : ", memberName)
	totalStart := time.Now()
	fmt.Println("***** [Start] Cluster Join Start : '", memberName, "' *****")

	c := cobrautil.GetOmcpctlConf("/var/lib/omcpctl/config.yaml")

	util.CmdExec("umount -l /mnt")
	defer util.CmdExec("umount -l /mnt")

	util.CmdExec2("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")

	//openmcpIP := cobrautil.GetOutboundIP()
	openmcpIP := cobrautil.GetEndpointIP()
	if !fileExists("/mnt/openmcp/" + openmcpIP) {
		fmt.Println("Failed Join List in OpenMCP Master: " + openmcpIP)
		fmt.Println("=> Not Yet Register OpenMCP.")
		fmt.Println("=> First You Must be Input the Next Command in 'OpenMCP Master Server(" + openmcpIP + ")' : omcpctl register openmcp")
		return
	}
	Lock.Unlock()

	alreadyJoined, err := cobrautil.CheckAlreadyJoinClusterWithPublicClusterName(memberName, "gke", "")
	if err != nil {
		fmt.Println("CheckAlreadyJoinClusterWithPublicClusterName Error : ", err)
		return
	}

	if alreadyJoined {
		return
	}

	_, err = util.CmdExec("gcloud container clusters get-credentials " + memberName)
	if err != nil {
		fmt.Println("[", err, "] No cluster found for name: "+memberName)
		return
	}

	kc := cobrautil.GetKubeConfig("/root/.kube/config")

	for i, c := range kc.Clusters {
		if strings.Contains(c.Name, memberName) {
			kc.Clusters[i].Name = memberName
			break
		}

	}
	for i, c := range kc.Contexts {
		if strings.Contains(c.Name, memberName) {
			kc.Contexts[i].Name = memberName
			kc.Contexts[i].Context.User = memberName + "-admin"
			kc.Contexts[i].Context.Cluster = memberName
			break

		}
	}
	for i, c := range kc.Users {
		if strings.Contains(c.Name, memberName) {
			kc.Users[i].Name = memberName + "-admin"
			break

		}
	}
	kc.CurrentContext = "openmcp"
	cobrautil.WriteKubeConfig(kc, "/root/.kube/config")

	// namespace terminating stuck force delete
	util.CmdExec2("kubectl get namespace kube-federation-system --context " + memberName + " -o json |jq '.spec = {\"finalizers\":[]}' >temp.json")
	util.CmdExec2("kubectl replace --raw \"/api/v1/namespaces/kube-federation-system/finalize\" -f ./temp.json --context " + memberName)
	util.CmdExec2("rm temp.json")

	start2 := time.Now()
	fmt.Println("***** [Start] 1. Cluster Join *****")
	util.CmdExec2("kubefedctl join " + memberName + " --cluster-context " + memberName + " --host-cluster-context openmcp --v=2")

	elapsed2 := time.Since(start2)
	log.Printf("Cluster Join Time : %s", elapsed2)
	fmt.Println("***** [End] 1. Cluster Join ***** ")

	start3 := time.Now()
	fmt.Println("***** [Start] 2. Init Service Deployments *****")

	installInitCluster(memberName, c.OpenmcpDir, "kube-dns")

	elapsed3 := time.Since(start3)
	log.Printf("Init Service Deployments Time : %s", elapsed3)
	fmt.Println("***** [End] 2. Init Service Deployments ***** ")

	kubeconfig, _ := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	genClient := genericclient.NewForConfigOrDie(kubeconfig)

	kubefedcluster := &fedv1b1.KubeFedCluster{}
	err = genClient.Get(context.TODO(), kubefedcluster, "kube-federation-system", memberName)
	if err != nil {
		fmt.Println(err)
	}
	labels := make(map[string]string)
	labels["platform"] = "gke"
	kubefedcluster.Labels = labels

	err = genClient.Update(context.TODO(), kubefedcluster)
	if err != nil {
		fmt.Println(err)
	}

	totalElapsed := time.Since(totalStart)
	log.Printf("Cluster Join Total Elapsed Time : %s", totalElapsed)
	fmt.Println("***** [End] Cluster Join Completed - "+memberName, "*****")
}

func joinEKSCluster(memberName string) {
	for {
		lockErr := Lock.TryLock()
		if lockErr != nil {
			fmt.Println("Mount Dir Using Another Works. Wait...")
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	totalStart := time.Now()
	fmt.Println("***** [Start] Cluster Join Start : '", memberName, "' *****")

	c := cobrautil.GetOmcpctlConf("/var/lib/omcpctl/config.yaml")

	util.CmdExec("umount -l /mnt")
	defer util.CmdExec("umount -l /mnt")

	util.CmdExec2("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")

	//openmcpIP := cobrautil.GetOutboundIP()
	openmcpIP := cobrautil.GetEndpointIP()
	if !fileExists("/mnt/openmcp/" + openmcpIP) {
		fmt.Println("Failed Join List in OpenMCP Master: " + openmcpIP)
		fmt.Println("=> Not Yet Register OpenMCP.")
		fmt.Println("=> First You Must be Input the Next Command in 'OpenMCP Master Server(" + openmcpIP + ")' : omcpctl register openmcp")
		return
	}
	Lock.Unlock()

	_, err := util.CmdExec("aws eks update-kubeconfig --name " + memberName)
	if err != nil {
		fmt.Println("[", err, "] No cluster found for name: "+memberName)
		return
	}

	alreadyJoined, err := cobrautil.CheckAlreadyJoinClusterWithPublicClusterName(memberName, "eks", "")
	if err != nil {
		fmt.Println("CheckAlreadyJoinClusterWithPublicClusterName Error : ", err)
		return
	}

	if alreadyJoined {
		return
	}

	kc := cobrautil.GetKubeConfig("/root/.kube/config")

	for i, c := range kc.Clusters {
		if strings.Contains(c.Name, memberName) {
			kc.Clusters[i].Name = memberName

			a := c.Cluster.Server
			lower_a := strings.ToLower(a)
			fmt.Println(a, " => ", lower_a)

			kc.Clusters[i].Cluster.Server = lower_a

			break
		}

	}
	for i, c := range kc.Contexts {
		if strings.Contains(c.Name, memberName) {
			kc.Contexts[i].Name = memberName
			kc.Contexts[i].Context.User = memberName + "-admin"
			kc.Contexts[i].Context.Cluster = memberName
			break

		}
	}
	for i, c := range kc.Users {
		if strings.Contains(c.Name, memberName) {
			kc.Users[i].Name = memberName + "-admin"
			break

		}
	}
	kc.CurrentContext = "openmcp"
	cobrautil.WriteKubeConfig(kc, "/root/.kube/config")

	// namespace terminating stuck force delete
	util.CmdExec2("kubectl get namespace kube-federation-system --context " + memberName + " -o json |jq '.spec = {\"finalizers\":[]}' >temp.json")
	util.CmdExec2("kubectl replace --raw \"/api/v1/namespaces/kube-federation-system/finalize\" -f ./temp.json --context " + memberName)
	util.CmdExec2("rm temp.json")

	start2 := time.Now()
	fmt.Println("***** [Start] 1. Cluster Join *****")
	util.CmdExec2("kubefedctl join " + memberName + " --cluster-context " + memberName + " --host-cluster-context openmcp --v=2")

	elapsed2 := time.Since(start2)
	log.Printf("Cluster Join Time : %s", elapsed2)
	fmt.Println("***** [End] 1. Cluster Join ***** ")

	start3 := time.Now()
	fmt.Println("***** [Start] 2. Init Service Deployments *****")

	installInitCluster(memberName, c.OpenmcpDir, "coredns")

	elapsed3 := time.Since(start3)
	log.Printf("Init Service Deployments Time : %s", elapsed3)
	fmt.Println("***** [End] 2. Init Service Deployments ***** ")

	kubeconfig, _ := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	genClient := genericclient.NewForConfigOrDie(kubeconfig)
	kubefedcluster := &fedv1b1.KubeFedCluster{}
	err = genClient.Get(context.TODO(), kubefedcluster, "kube-federation-system", memberName)
	if err != nil {
		fmt.Println(err)
	}
	labels := make(map[string]string)
	labels["platform"] = "eks"
	kubefedcluster.Labels = labels

	err = genClient.Update(context.TODO(), kubefedcluster)
	if err != nil {
		fmt.Println(err)
	}

	totalElapsed := time.Since(totalStart)
	log.Printf("Cluster Join Total Elapsed Time : %s", totalElapsed)
	fmt.Println("***** [End] Cluster Join Completed - "+memberName, "*****")
}

func joinAKSCluster(memberName string, resourceGroupName string) {

	for {
		lockErr := Lock.TryLock()
		if lockErr != nil {
			fmt.Println("Mount Dir Using Another Works. Wait...")
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	totalStart := time.Now()
	fmt.Println("***** [Start] Cluster Join Start : '", memberName, "' *****")

	c := cobrautil.GetOmcpctlConf("/var/lib/omcpctl/config.yaml")

	util.CmdExec("umount -l /mnt")
	defer util.CmdExec("umount -l /mnt")

	util.CmdExec2("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")

	//openmcpIP := cobrautil.GetOutboundIP()
	openmcpIP := cobrautil.GetEndpointIP()
	if !fileExists("/mnt/openmcp/" + openmcpIP) {
		fmt.Println("Failed Join List in OpenMCP Master: " + openmcpIP)
		fmt.Println("=> Not Yet Register OpenMCP.")
		fmt.Println("=> First You Must be Input the Next Command in 'OpenMCP Master Server(" + openmcpIP + ")' : omcpctl register openmcp")
		return
	}

	Lock.Unlock()

	_, err := util.CmdExec("az aks get-credentials --name " + memberName + " --resource-group " + resourceGroupName)
	if err != nil {
		fmt.Println("[", err, "] No cluster found for name: "+memberName)
		return
	}

	alreadyJoined, err := cobrautil.CheckAlreadyJoinClusterWithPublicClusterName(memberName, "aks", resourceGroupName)
	if err != nil {
		fmt.Println("CheckAlreadyJoinClusterWithPublicClusterName Error : ", err)
		return
	}

	if alreadyJoined {
		return
	}

	kc := cobrautil.GetKubeConfig("/root/.kube/config")

	for i, c := range kc.Contexts {
		if strings.Contains(c.Name, memberName) {
			kc.Contexts[i].Name = memberName
			kc.Contexts[i].Context.User = memberName + "-admin"
			kc.Contexts[i].Context.Cluster = memberName
			break

		}
	}
	for i, c := range kc.Users {
		if strings.Contains(c.Name, memberName) {
			kc.Users[i].Name = memberName + "-admin"
			break

		}
	}
	kc.CurrentContext = "openmcp"
	cobrautil.WriteKubeConfig(kc, "/root/.kube/config")

	// namespace terminating stuck force delete
	util.CmdExec2("kubectl get namespace kube-federation-system --context " + memberName + " -o json |jq '.spec = {\"finalizers\":[]}' >temp.json")
	util.CmdExec2("kubectl replace --raw \"/api/v1/namespaces/kube-federation-system/finalize\" -f ./temp.json --context " + memberName)
	util.CmdExec2("rm temp.json")

	start2 := time.Now()
	fmt.Println("***** [Start] 1. Cluster Join *****")
	util.CmdExec2("kubefedctl join " + memberName + " --cluster-context " + memberName + " --host-cluster-context openmcp --v=2")

	elapsed2 := time.Since(start2)
	log.Printf("Cluster Join Time : %s", elapsed2)
	fmt.Println("***** [End] 1. Cluster Join ***** ")

	start3 := time.Now()
	fmt.Println("***** [Start] 2. Init Service Deployments *****")

	installInitCluster(memberName, c.OpenmcpDir, "coredns")

	elapsed3 := time.Since(start3)
	log.Printf("Init Service Deployments Time : %s", elapsed3)
	fmt.Println("***** [End] 2. Init Service Deployments ***** ")

	kubeconfig, _ := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	genClient := genericclient.NewForConfigOrDie(kubeconfig)

	kubefedcluster := &fedv1b1.KubeFedCluster{}
	err = genClient.Get(context.TODO(), kubefedcluster, "kube-federation-system", memberName)
	if err != nil {
		fmt.Println(err)
	}
	labels := make(map[string]string)
	labels["platform"] = "aks"
	kubefedcluster.Labels = labels

	err = genClient.Update(context.TODO(), kubefedcluster)
	if err != nil {
		fmt.Println(err)
	}

	totalElapsed := time.Since(totalStart)
	log.Printf("Cluster Join Total Elapsed Time : %s", totalElapsed)
	fmt.Println("***** [End] Cluster Join Completed - "+memberName, "*****")
}

func installInitCluster(clusterName, openmcpDir, dnsKind string) {
	fmt.Println("Init Module Deployment Start - " + clusterName)
	install_dir := filepath.Join(openmcpDir, "install_openmcp/member")

	util.CmdExec2("cp " + install_dir + "/metric-collector/operator/operator.yaml " + install_dir + "/metric-collector/operator.yaml")
	util.CmdExec2("sed -i 's|REPLACE_CLUSTER_NAME|\"" + clusterName + "\"|g' " + install_dir + "/metric-collector/operator.yaml")
	initYamls := []string{"namespace", "custom-metrics-apiserver", "metallb", "metric-collector", "metrics-server", "nginx-ingress-controller", "configmap"}

	util.CmdExec2("kubectl create ns openmcp --context " + clusterName)
	for _, initYaml := range initYamls {

		if initYaml == "configmap" {
			if dnsKind == "coredns" {
				util.CmdExec2("kubectl apply -f " + install_dir + "/" + initYaml + "/coredns --context " + clusterName)
			} else {
				util.CmdExec2("kubectl apply -f " + install_dir + "/" + initYaml + "/kubedns --context " + clusterName)
			}
		} else {
			util.CmdExec2("kubectl apply -f " + install_dir + "/" + initYaml + " --context " + clusterName)
		}

	}

	util.CmdExec2("chmod 755 " + install_dir + "/vertical-pod-autoscaler/hack/*")
	util.CmdExec2(install_dir + "/vertical-pod-autoscaler/hack/vpa-up.sh " + clusterName)
	fmt.Println("Init Module Deployment Finished - " + clusterName)

}

func init() {
	rootCmd.AddCommand(joinCmd)
	go RunCache()
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
