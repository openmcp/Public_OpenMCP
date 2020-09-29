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