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


func OpenMCPSecretInfo(osec *ketiv1alpha1.OpenMCPSecret) []string{

	namespace := osec.Namespace
	name := osec.Name

	secType := string(osec.Spec.Template.Type)
	dataCount := strconv.Itoa(len(osec.Spec.Template.Data))

	cluster := ""

	keys := make([]string, 0, len(osec.Status.ClusterMaps))
	for k := range osec.Status.ClusterMaps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys{
		if osec.Status.ClusterMaps[k] >= 1 {
			cluster += k+":"+strconv.Itoa(int(osec.Status.ClusterMaps[k]))+" "
		}
	}
	age := cobrautil.GetAge(osec.CreationTimestamp.Time)

	data := []string{namespace, name, secType, dataCount, cluster, age}

	return data
}
func PrintOpenMCPSecret(body []byte) {
	osec := ketiv1alpha1.OpenMCPSecret{}
	err := yaml.Unmarshal(body, &osec)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	data := OpenMCPSecretInfo(&osec)
	datas = append(datas, data)

	DrawOpenMCPSecretTable(datas)

}
func PrintOpenMCPSecretList(body []byte) {
	resourceStruct := ketiv1alpha1.OpenMCPSecretList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}

	for _, osec := range resourceStruct.Items {
		data := OpenMCPSecretInfo(&osec)
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

	DrawOpenMCPSecretTable(datas)

}

func DrawOpenMCPSecretTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "TYPE", "DATA", "CLUSTER", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}