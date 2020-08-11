package resource

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/ghodss/yaml"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"

)

func CustomResourceDefinitionInfo(crd *apiextensionsv1.CustomResourceDefinition) []string{
	createdAt := crd.CreationTimestamp.Time.String()

	data := []string{crd.Name, createdAt}

	return data
}
func PrintCustomResourceDefinition(body []byte) {
	crd := apiextensionsv1.CustomResourceDefinition{}
	err := yaml.Unmarshal(body, &crd)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := CustomResourceDefinitionInfo(&crd)
	datas = append(datas, data)

	DrawCustomResourceDefinitionTable(datas)

}
func PrintCustomResourceDefinitionList(body []byte) {
	resourceStruct := apiextensionsv1.CustomResourceDefinitionList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, crd := range resourceStruct.Items {
		data := CustomResourceDefinitionInfo(&crd)
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

	DrawCustomResourceDefinitionTable(datas)

}

func DrawCustomResourceDefinitionTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "CREATED AT"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}