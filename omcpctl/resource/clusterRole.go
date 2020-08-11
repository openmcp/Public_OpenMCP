package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"

	rbacv1 "k8s.io/api/rbac/v1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
)

func ClusterRoleInfo(cr *rbacv1.ClusterRole) []string{


	age := cobrautil.GetAge(cr.CreationTimestamp.Time)

	data := []string{cr.Name, age}

	return data
}
func PrintClusterRole(body []byte) {
	cr := rbacv1.ClusterRole{}
	err := yaml.Unmarshal(body, &cr)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := ClusterRoleInfo(&cr)
	datas = append(datas, data)

	DrawClusterRoleTable(datas)

}
func PrintClusterRoleList(body []byte) {
	resourceStruct := rbacv1.ClusterRoleList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, cr := range resourceStruct.Items {
		data := ClusterRoleInfo(&cr)
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

	DrawClusterRoleTable(datas)

}

func DrawClusterRoleTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}