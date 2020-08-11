package resource

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/ghodss/yaml"
	ketiv1alpha1 "openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"

	cobrautil "openmcp/openmcp/omcpctl/util"

	"os"
	"sort"
	"strconv"
)

func OpenMCPIngressInfo(oing *ketiv1alpha1.OpenMCPIngress) []string{

	namespace := oing.Namespace
	name := oing.Name
	clusters := ""
	keys := make([]string, 0, len(oing.Status.ClusterMaps))
	for k := range oing.Status.ClusterMaps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys{
		if oing.Status.ClusterMaps[k] >= 1 {
			clusters += k+":"+strconv.Itoa(int(oing.Status.ClusterMaps[k]))+" "
		}
	}
	
	age := cobrautil.GetAge(oing.CreationTimestamp.Time)
	data := []string{namespace, name, clusters, age}

	return data
}
func PrintOpenMCPIngress(body []byte) {
	oing := ketiv1alpha1.OpenMCPIngress{}
	err := yaml.Unmarshal(body, &oing)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := OpenMCPIngressInfo(&oing)
	datas = append(datas, data)

	DrawOpenMCPIngressTable(datas)

}
func PrintOpenMCPIngressList(body []byte) {
	resourceStruct := ketiv1alpha1.OpenMCPIngressList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, oing := range resourceStruct.Items {
		data := OpenMCPIngressInfo(&oing)
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

	DrawOpenMCPIngressTable(datas)

}

func DrawOpenMCPIngressTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "CLUSTER", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}


//
//func GetOpenmcpIngress(oingName string) {
//	kubeconfig, _ := cobrautil.BuildConfigFromFlags("openmcp", "/root/.kube/config")
//	crdClient, _ := clientV1alpha1.NewForConfig(kubeconfig)
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
//	oing, err := crdClient.OpenMCPIngress(ns).Get(oingName, metav1.GetOptions{})
//	if err != nil && errors.IsNotFound(err) {
//		fmt.Println("No resources found in " + ns + " namespace.")
//		return
//	} else if err != nil {
//		fmt.Println(err)
//		return
//	}
//	if cobrautil.Option_filetype == "yaml"{
//		data, err := yaml.Marshal(&oing)
//		if err != nil {
//			log.Fatalf("Unmarshal: %v", err)
//		}
//		fmt.Println(string(data))
//		return
//	}
//	datas := [][]string{}
//	namespace := oing.Namespace
//	name := oing.Name
//	clusters := ""
//	keys := make([]string, 0, len(oing.Status.ClusterMaps))
//	for k := range oing.Status.ClusterMaps {
//		keys = append(keys, k)
//	}
//	sort.Strings(keys)
//	for _, k := range keys{
//		if oing.Status.ClusterMaps[k] >= 1 {
//			clusters += k+":"+strconv.Itoa(int(oing.Status.ClusterMaps[k]))+" "
//		}
//	}
//	age := ""
//	data := []string{namespace, name, clusters, age}
//	datas = append(datas, data)
//
//	table := tablewriter.NewWriter(os.Stdout)
//	table.SetHeader([]string{"NS", "NAME", "CLUSTERS", "AGE"})
//	table.SetBorder(false)
//	table.AppendBulk(datas)
//	table.Render()
//}
//func GetOpenmcpIngressList() {
//	kubeconfig, _ := cobrautil.BuildConfigFromFlags("openmcp", "/root/.kube/config")
//	crdClient, _ := clientV1alpha1.NewForConfig(kubeconfig)
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
//	oingList, err := crdClient.OpenMCPIngress(ns).List(metav1.ListOptions{})
//	if err != nil && errors.IsNotFound(err) {
//		fmt.Println("No resources found in " + ns + " namespace.")
//		return
//	} else if err != nil {
//		fmt.Println(err)
//		return
//	}
//	if cobrautil.Option_filetype == "yaml"{
//		data, err := yaml.Marshal(&oingList)
//		if err != nil {
//			log.Fatalf("Unmarshal: %v", err)
//		}
//		fmt.Println(string(data))
//		return
//	}
//	datas := [][]string{}
//	for _, oing := range oingList.Items {
//		namespace := oing.Namespace
//		name := oing.Name
//		clusters := ""
//		keys := make([]string, 0, len(oing.Status.ClusterMaps))
//		for k := range oing.Status.ClusterMaps {
//			keys = append(keys, k)
//		}
//		sort.Strings(keys)
//		for _, k := range keys{
//			if oing.Status.ClusterMaps[k] >= 1 {
//				clusters += k+":"+strconv.Itoa(int(oing.Status.ClusterMaps[k]))+" "
//			}
//		}
//		age := ""
//		data := []string{namespace, name, clusters, age}
//		datas = append(datas, data)
//	}
//
//	table := tablewriter.NewWriter(os.Stdout)
//	table.SetHeader([]string{"NS", "NAME", "CLUSTERS", "AGE"})
//	table.SetBorder(false)
//	table.AppendBulk(datas)
//	table.Render()
//}
