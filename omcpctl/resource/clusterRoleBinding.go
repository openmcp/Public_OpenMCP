package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	rbacv1 "k8s.io/api/rbac/v1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
)

func ClusterRoleBindingInfo(crb *rbacv1.ClusterRoleBinding) []string{


	age := cobrautil.GetAge(crb.CreationTimestamp.Time)

	data := []string{crb.Name, age}

	return data
}
func PrintClusterRoleBinding(body []byte) {
	crb := rbacv1.ClusterRoleBinding{}
	err := yaml.Unmarshal(body, &crb)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	data := ClusterRoleBindingInfo(&crb)
	datas = append(datas, data)

	DrawClusterRoleBindingTable(datas)

}
func PrintClusterRoleBindingList(body []byte) {
	resourceStruct := rbacv1.ClusterRoleBindingList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}

	for _, crb := range resourceStruct.Items {
		data := ClusterRoleBindingInfo(&crb)
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

	DrawClusterRoleBindingTable(datas)

}

func DrawClusterRoleBindingTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}