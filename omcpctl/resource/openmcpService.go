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

func OpenMCPServiceInfo(osvc *ketiv1alpha1.OpenMCPService) []string{

	namespace := osvc.Namespace
	name := osvc.Name
	clusters := ""
	keys := make([]string, 0, len(osvc.Status.ClusterMaps))
	for k := range osvc.Status.ClusterMaps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys{
		if osvc.Status.ClusterMaps[k] >= 1 {
			clusters += k+":"+strconv.Itoa(int(osvc.Status.ClusterMaps[k]))+" "
		}
	}
	age := cobrautil.GetAge(osvc.CreationTimestamp.Time)
	data := []string{namespace, name, clusters, age}

	return data
}
func PrintOpenMCPService(body []byte) {
	osvc := ketiv1alpha1.OpenMCPService{}
	err := yaml.Unmarshal(body, &osvc)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := OpenMCPServiceInfo(&osvc)
	datas = append(datas, data)

	DrawOpenMCPServiceTable(datas)

}
func PrintOpenMCPServiceList(body []byte) {
	resourceStruct := ketiv1alpha1.OpenMCPServiceList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, osvc := range resourceStruct.Items {
		data := OpenMCPServiceInfo(&osvc)
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

	DrawOpenMCPServiceTable(datas)

}

func DrawOpenMCPServiceTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "CLUSTER", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}











//func GetOpenMCPService(osvcName string) {
//	kubeconfig, _ := cobrautil.BuildConfigFromFlags("OpenMCP", "/root/.kube/config")
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
//	osvc, err := crdClient.OpenMCPService(ns).Get(osvcName, metav1.GetOptions{})
//	if err != nil && errors.IsNotFound(err) {
//		fmt.Println("No resources found in " + ns + " namespace.")
//		return
//	} else if err != nil {
//		fmt.Println(err)
//		return
//	}
//	if cobrautil.Option_filetype == "yaml"{
//		data, err := yaml.Marshal(&osvc)
//		if err != nil {
//			log.Fatalf("Unmarshal: %v", err)
//		}
//		fmt.Println(string(data))
//		return
//	}
//	datas := [][]string{}
//	namespace := osvc.Namespace
//	name := osvc.Name
//	clusters := ""
//	keys := make([]string, 0, len(osvc.Status.ClusterMaps))
//	for k := range osvc.Status.ClusterMaps {
//		keys = append(keys, k)
//	}
//	sort.Strings(keys)
//	for _, k := range keys{
//		if osvc.Status.ClusterMaps[k] >= 1 {
//			clusters += k+":"+strconv.Itoa(int(osvc.Status.ClusterMaps[k]))+" "
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
//func GetOpenMCPServiceList() {
//	kubeconfig, _ := cobrautil.BuildConfigFromFlags("OpenMCP", "/root/.kube/config")
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
//	osvcList, err := crdClient.OpenMCPService(ns).List(metav1.ListOptions{})
//	if err != nil && errors.IsNotFound(err) {
//		fmt.Println("No resources found in " + ns + " namespace.")
//		return
//	} else if err != nil {
//		fmt.Println(err)
//		return
//	}
//	if cobrautil.Option_filetype == "yaml"{
//		data, err := yaml.Marshal(&osvcList)
//		if err != nil {
//			log.Fatalf("Unmarshal: %v", err)
//		}
//		fmt.Println(string(data))
//		return
//	}
//	datas := [][]string{}
//	for _, osvc := range osvcList.Items {
//		namespace := osvc.Namespace
//		name := osvc.Name
//		clusters := ""
//		keys := make([]string, 0, len(osvc.Status.ClusterMaps))
//		for k := range osvc.Status.ClusterMaps {
//			keys = append(keys, k)
//		}
//		sort.Strings(keys)
//		for _, k := range keys{
//			if osvc.Status.ClusterMaps[k] >= 1 {
//				clusters += k+":"+strconv.Itoa(int(osvc.Status.ClusterMaps[k]))+" "
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