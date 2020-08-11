package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	cobrautil "openmcp/openmcp/omcpctl/util"
	ketiv1alpha1 "openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"
	"os"
	"sort"
	"strconv"
)


func OpenMCPDeploymentInfo(od *ketiv1alpha1.OpenMCPDeployment) []string{

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
	age := cobrautil.GetAge(od.CreationTimestamp.Time)

	data := []string{namespace, name, cluster, age}

	return data
}
func PrintOpenMCPDeployment(body []byte) {
	od := ketiv1alpha1.OpenMCPDeployment{}
	err := yaml.Unmarshal(body, &od)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := OpenMCPDeploymentInfo(&od)
	datas = append(datas, data)

	DrawOpenMCPDeploymentTable(datas)

}
func PrintOpenMCPDeploymentList(body []byte) {
	resourceStruct := ketiv1alpha1.OpenMCPDeploymentList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, od := range resourceStruct.Items {
		data := OpenMCPDeploymentInfo(&od)
		datas = append(datas, data)
	}

	if len(resourceStruct.Items) == 0 {
		ns := "default"
		if cobrautil.Option_namespace != "" {
			ns = cobrautil.Option_namespace
		}
		errMsg := "No resources found"
		if !cobrautil.Option_allnamespace {
			errMsg += " in "+ ns +" namespace."
		}
		fmt.Println(errMsg)
		return
	}

	DrawOpenMCPDeploymentTable(datas)

}

func DrawOpenMCPDeploymentTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "CLUSTER", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}










//func GetOpenMCPDeployment(odeployName string) {
//	kubeconfig, _ := cobrautil.BuildConfigFromFlags("openmcp", "/root/.kube/config")
//	crdClient, _ := clientV1alpha1.NewForConfig(kubeconfig)
//
//	ns := ""
//	if cobrautil.Option_allnamespace {
//		ns = metav1.NamespaceAll
//	} else {
//		if cobrautil.Option_namespace != "" {
//			ns = cobrautil.Option_namespace
//		} else {
//			ns = "default"
//		}
//	}
//	od, err := crdClient.OpenMCPDeployment(ns).Get(odeployName, metav1.GetOptions{})
//	if err != nil && errors.IsNotFound(err) {
//		fmt.Println("No resources found in " + ns + " namespace.")
//		return
//	} else if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	if cobrautil.Option_filetype == "yaml"{
//		data, err := yaml.Marshal(&od)
//		if err != nil {
//			log.Fatalf("Unmarshal: %v", err)
//		}
//		fmt.Println(string(data))
//		return
//	}
//
//	datas := [][]string{}
//	namespace := od.Namespace
//	name := od.Name
//	cluster := ""
//
//	keys := make([]string, 0, len(od.Status.ClusterMaps))
//	for k := range od.Status.ClusterMaps {
//		keys = append(keys, k)
//	}
//	sort.Strings(keys)
//	for _, k := range keys{
//		if od.Status.ClusterMaps[k] >= 1 {
//			cluster += k+":"+strconv.Itoa(int(od.Status.ClusterMaps[k]))+" "
//		}
//	}
//	age := ""
//	data := []string{namespace, name, cluster, age}
//	datas = append(datas, data)
//
//	table := tablewriter.NewWriter(os.Stdout)
//	table.SetHeader([]string{"NS", "NAME", "CLUSTER", "AGE"})
//	table.SetBorder(false)
//	table.AppendBulk(datas)
//	table.Render()
//
//}
//func GetOpenMCPDeploymentList() {
//	kubeconfig, _ := cobrautil.BuildConfigFromFlags("openmcp", "/root/.kube/config")
//	crdClient, _ := clientV1alpha1.NewForConfig(kubeconfig)
//
//	ns := ""
//	if cobrautil.Option_allnamespace {
//		ns = metav1.NamespaceAll
//	} else {
//		if cobrautil.Option_namespace != ""{
//			ns = cobrautil.Option_namespace
//		}
//	}
//	odList, err := crdClient.OpenMCPDeployment(ns).List(metav1.ListOptions{})
//
//	if err != nil && errors.IsNotFound(err) {
//		fmt.Println("No Resource exist")
//	}
//	if len(odList.Items) == 0 {
//		if ns == ""{
//			fmt.Println("No resources found.")
//		} else{
//			fmt.Println("No resources found in "+ns+" namespace.")
//		}
//
//		return
//	}
//	if cobrautil.Option_filetype == "yaml"{
//		data, err := yaml.Marshal(&odList)
//		if err != nil {
//			log.Fatalf("Unmarshal: %v", err)
//		}
//		fmt.Println(string(data))
//		return
//	}
//
//	datas := [][]string{}
//	for _, od := range odList.Items {
//		namespace := od.Namespace
//		name := od.Name
//		cluster := ""
//
//		keys := make([]string, 0, len(od.Status.ClusterMaps))
//		for k := range od.Status.ClusterMaps {
//			keys = append(keys, k)
//		}
//		sort.Strings(keys)
//		for _, k := range keys{
//			if od.Status.ClusterMaps[k] >= 1 {
//				cluster += k+":"+strconv.Itoa(int(od.Status.ClusterMaps[k]))+" "
//			}
//		}
//		age := ""
//		data := []string{namespace, name, cluster, age}
//		datas = append(datas, data)
//	}
//
//	table := tablewriter.NewWriter(os.Stdout)
//	table.SetHeader([]string{"NS", "NAME", "CLUSTER", "AGE"})
//	table.SetBorder(false)
//	table.AppendBulk(datas)
//	table.Render()
//
//
//}
