package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	corev1 "k8s.io/api/core/v1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
)

func NamespaceInfo(ns *corev1.Namespace) []string{

	status := string(ns.Status.Phase)
	age := cobrautil.GetAge(ns.CreationTimestamp.Time)

	data := []string{ns.Name, status, age}

	return data
}
func PrintNamespace(body []byte) {
	ns := corev1.Namespace{}
	err := yaml.Unmarshal(body, &ns)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	data := NamespaceInfo(&ns)
	datas = append(datas, data)

	DrawNamespaceTable(datas)

}
func PrintNamespaceList(body []byte) {
	resourceStruct := corev1.NamespaceList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}

	for _, ns := range resourceStruct.Items {
		data := NamespaceInfo(&ns)
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

	DrawNamespaceTable(datas)

}

func DrawNamespaceTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "STATUS", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}