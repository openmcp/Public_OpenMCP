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

func OpenMCPServiceInfo(osvc *resourcev1alpha1.OpenMCPService) []string{

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
	osvc := resourcev1alpha1.OpenMCPService{}
	err := yaml.Unmarshal(body, &osvc)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	data := OpenMCPServiceInfo(&osvc)
	datas = append(datas, data)

	DrawOpenMCPServiceTable(datas)

}
func PrintOpenMCPServiceList(body []byte) {
	resourceStruct := resourcev1alpha1.OpenMCPServiceList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
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
