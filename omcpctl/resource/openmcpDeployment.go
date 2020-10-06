package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	cobrautil "openmcp/openmcp/omcpctl/util"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"
	"os"
	"sort"
	"strconv"
)


func OpenMCPDeploymentInfo(od *resourcev1alpha1.OpenMCPDeployment) []string{

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
	od := resourcev1alpha1.OpenMCPDeployment{}
	err := yaml.Unmarshal(body, &od)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	data := OpenMCPDeploymentInfo(&od)
	datas = append(datas, data)

	DrawOpenMCPDeploymentTable(datas)

}
func PrintOpenMCPDeploymentList(body []byte) {
	resourceStruct := resourcev1alpha1.OpenMCPDeploymentList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
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
