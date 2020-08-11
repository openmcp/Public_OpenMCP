package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"

	rbacv1 "k8s.io/api/rbac/v1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
)

func RoleBindingInfo(rb *rbacv1.RoleBinding) []string{


	age := cobrautil.GetAge(rb.CreationTimestamp.Time)

	data := []string{rb.Namespace, rb.Name, age}

	return data
}
func PrintRoleBinding(body []byte) {
	rb := rbacv1.RoleBinding{}
	err := yaml.Unmarshal(body, &rb)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := RoleBindingInfo(&rb)
	datas = append(datas, data)

	DrawRoleBindingTable(datas)

}
func PrintRoleBindingList(body []byte) {
	resourceStruct := rbacv1.RoleBindingList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, rb := range resourceStruct.Items {
		data := RoleBindingInfo(&rb)
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

	DrawRoleBindingTable(datas)

}

func DrawRoleBindingTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}