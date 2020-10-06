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


func OpenMCPConfigMapInfo(ocm *resourcev1alpha1.OpenMCPConfigMap) []string{

	namespace := ocm.Namespace
	name := ocm.Name
	dataCount := strconv.Itoa(len(ocm.Spec.Template.Data))

	cluster := ""

	keys := make([]string, 0, len(ocm.Status.ClusterMaps))
	for k := range ocm.Status.ClusterMaps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys{
		if ocm.Status.ClusterMaps[k] >= 1 {
			cluster += k+":"+strconv.Itoa(int(ocm.Status.ClusterMaps[k]))+" "
		}
	}
	age := cobrautil.GetAge(ocm.CreationTimestamp.Time)

	data := []string{namespace, name, dataCount, cluster, age}

	return data
}
func PrintOpenMCPConfigMap(body []byte) {
	ocm := resourcev1alpha1.OpenMCPConfigMap{}
	err := yaml.Unmarshal(body, &ocm)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	data := OpenMCPConfigMapInfo(&ocm)
	datas = append(datas, data)

	DrawOpenMCPConfigMapTable(datas)

}
func PrintOpenMCPConfigMapList(body []byte) {
	resourceStruct := resourcev1alpha1.OpenMCPConfigMapList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}

	for _, ocm := range resourceStruct.Items {
		data := OpenMCPConfigMapInfo(&ocm)
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

	DrawOpenMCPConfigMapTable(datas)

}

func DrawOpenMCPConfigMapTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "DATA", "CLUSTER", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}