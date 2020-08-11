package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	corev1 "k8s.io/api/core/v1"
	"strconv"

	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
)

func ServiceAccountInfo(sa *corev1.ServiceAccount) []string{

	secrets := strconv.Itoa(len(sa.Secrets))

	age := cobrautil.GetAge(sa.CreationTimestamp.Time)

	data := []string{sa.Namespace, sa.Name, secrets ,age}

	return data
}
func PrintServiceAccount(body []byte) {
	sa := corev1.ServiceAccount{}
	err := yaml.Unmarshal(body, &sa)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := ServiceAccountInfo(&sa)
	datas = append(datas, data)

	DrawServiceAccountTable(datas)

}
func PrintServiceAccountList(body []byte) {
	resourceStruct := corev1.ServiceAccountList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, sa := range resourceStruct.Items {
		data := ServiceAccountInfo(&sa)
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

	DrawServiceAccountTable(datas)

}

func DrawServiceAccountTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "SECRETS", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}