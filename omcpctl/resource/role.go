package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"

	rbacv1 "k8s.io/api/rbac/v1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
)

func RoleInfo(role *rbacv1.Role) []string{


	age := cobrautil.GetAge(role.CreationTimestamp.Time)

	data := []string{role.Namespace, role.Name, age}

	return data
}
func PrintRole(body []byte) {
	role := rbacv1.Role{}
	err := yaml.Unmarshal(body, &role)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := RoleInfo(&role)
	datas = append(datas, data)

	DrawRoleTable(datas)

}
func PrintRoleList(body []byte) {
	resourceStruct := rbacv1.RoleList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, role := range resourceStruct.Items {
		data := RoleInfo(&role)
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

	DrawRoleTable(datas)

}

func DrawRoleTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}