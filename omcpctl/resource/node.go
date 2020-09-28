package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	corev1 "k8s.io/api/core/v1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
)


func NodeInfo(no *corev1.Node) []string{

	nodeName := no.Name
	nodeStatus := ""
	for _, cond := range no.Status.Conditions {
		if cond.Type == "Ready"{
			nodeStatus = string(cond.Status)
			break
		}
	}
	nodeRegion := no.Labels["failure-domain.beta.kubernetes.io/region"]
	nodeZone := no.Labels["failure-domain.beta.kubernetes.io/zone"]
	nodeAddress := no.Status.Addresses[0].Address
	nodeMaster := "false"
	if _, ok := no.Labels["node-role.kubernetes.io/master"]; ok {
		nodeMaster = "true"
	}
	nodeSchedule := "Yes"
	for _, taint := range no.Spec.Taints{
		if taint.Effect == "NoSchedule" {
			nodeSchedule ="No"
		}
	}

	age := cobrautil.GetAge(no.CreationTimestamp.Time)

	data := []string{nodeName, nodeStatus, nodeRegion, nodeZone, nodeAddress, nodeMaster, nodeSchedule, age}

	return data
}
func PrintNode(body []byte) {
	no := corev1.Node{}
	err := yaml.Unmarshal(body, &no)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	data := NodeInfo(&no)
	datas = append(datas, data)

	DrawNodeTable(datas)

}
func PrintNodeList(body []byte) {
	resourceStruct := corev1.NodeList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, no := range resourceStruct.Items {
		data := NodeInfo(&no)
		datas = append(datas, data)
	}

	if len(resourceStruct.Items) == 0 {
		ns := "default"
		if cobrautil.Option_namespace != "" {
			ns = cobrautil.Option_namespace
		}
		fmt.Println("No resources found in "+ ns +" Node.")
		return
	}

	DrawNodeTable(datas)

}

func DrawNodeTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NodeName", "Status", "Region", "Zones", "address", "Master", "Schedule", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}














//func GetNode(nodeName string) {
//	c := cobrautil.GetKubeConfig("/root/.kube/config")
//
//	for _, kubecontext := range c.Contexts {
//		if cobrautil.Option_context != "" && cobrautil.Option_context != kubecontext.Name{
//			continue
//		}
//		kubeconfig, _ := cobrautil.BuildConfigFromFlags(kubecontext.Name, "/root/.kube/config")
//		genClient := genericclient.NewForConfigOrDie(kubeconfig)
//
//		node := &corev1.Node{}
//		err := genClient.Get(context.Background(), node,"default", nodeName)
//		if err != nil && errors.IsNotFound(err){
//			continue
//		}
//
//		if cobrautil.Option_filetype == ""{
//			datas := [][]string{}
//			nodeName := node.Name
//			nodeStatus := ""
//			for _, cond := range node.Status.Conditions {
//				if cond.Type == "Ready"{
//					nodeStatus = string(cond.Status)
//					break
//				}
//			}
//			nodeRegion := node.Labels["failure-domain.beta.kubernetes.io/region"]
//			nodeZone := node.Labels["failure-domain.beta.kubernetes.io/zone"]
//			nodeAddress := node.Status.Addresses[0].Address
//			nodeMaster := "false"
//			if _, ok := node.Labels["node-role.kubernetes.io/master"]; ok {
//				nodeMaster = "true"
//			}
//			nodeSchedule := "Yes"
//			for _, taint := range node.Spec.Taints{
//				if taint.Effect == "NoSchedule" {
//					nodeSchedule ="No"
//				}
//			}
//
//			nodeClusterName := kubecontext.Context.Cluster
//
//			data := []string{nodeName, nodeStatus, nodeRegion, nodeZone, nodeAddress, nodeMaster, nodeSchedule, nodeClusterName}
//
//			datas = append(datas, data)
//			table := tablewriter.NewWriter(os.Stdout)
//			table.SetHeader([]string{"NodeName", "Status", "Region", "Zones", "address", "Master", "Schedule", "ClusterName", "AGE"})
//			table.SetBorder(false)
//			table.AppendBulk(datas)
//			table.Render()
//
//		} else if cobrautil.Option_filetype == "yaml"{
//			if cobrautil.Option_context == "" {
//				fmt.Println("The Command Must Need Option '--context cobrautil.Option_context'")
//			} else {
//				res, _ := util.CmdExec("kubectl get node " + nodeName + " --context "+ cobrautil.Option_context +" -o yaml")
//				fmt.Println(res)
//			}
//
//		}
//
//	}
//
//}
//
//func GetNodeList() {
//	c := cobrautil.GetKubeConfig("/root/.kube/config")
//	datas := [][]string{}
//
//	for _, kubecontext := range c.Contexts {
//		if cobrautil.Option_context != "" && cobrautil.Option_context != kubecontext.Name{
//			continue
//		}
//		kubeconfig, _ := cobrautil.BuildConfigFromFlags(kubecontext.Name, "/root/.kube/config")
//		genClient := genericclient.NewForConfigOrDie(kubeconfig)
//
//		nodeList := &corev1.NodeList{}
//		err := genClient.List(context.Background(), nodeList,"default")
//		if err != nil && errors.IsNotFound(err){
//			continue
//		}
//		if cobrautil.Option_filetype == ""{
//			for _, node := range nodeList.Items {
//				nodeName := node.Name
//				nodeStatus := ""
//				for _, cond := range node.Status.Conditions {
//					if cond.Type == "Ready"{
//						nodeStatus = string(cond.Status)
//						break
//					}
//				}
//				nodeRegion := node.Labels["failure-domain.beta.kubernetes.io/region"]
//				nodeZone := node.Labels["failure-domain.beta.kubernetes.io/zone"]
//				nodeAddress := node.Status.Addresses[0].Address
//				nodeMaster := "false"
//				if _, ok := node.Labels["node-role.kubernetes.io/master"]; ok {
//					nodeMaster = "true"
//				}
//				nodeSchedule := "Yes"
//				for _, taint := range node.Spec.Taints{
//					if taint.Effect == "NoSchedule" {
//						nodeSchedule ="No"
//					}
//				}
//
//				nodeClusterName := kubecontext.Context.Cluster
//
//				data := []string{nodeName, nodeStatus, nodeRegion, nodeZone, nodeAddress, nodeMaster, nodeSchedule, nodeClusterName}
//				datas = append(datas, data)
//
//			}
//		} else if cobrautil.Option_filetype == "yaml"{
//			if cobrautil.Option_context == "" {
//				fmt.Println("The Command Must Need Option '--context cobrautil.Option_context'")
//			} else {
//				res, _ := util.CmdExec("kubectl get node --context "+ cobrautil.Option_context +" -o yaml")
//				fmt.Println(res)
//			}
//		}
//
//
//
//	}
//	if cobrautil.Option_filetype == ""{
//		table := tablewriter.NewWriter(os.Stdout)
//		table.SetHeader([]string{"NodeName", "Status", "Region", "Zones", "address", "Master", "Schedule", "ClusterName", "AGE"})
//		table.SetBorder(false)
//		table.AppendBulk(datas)
//		table.Render()
//	}
//
//}
//
