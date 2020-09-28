package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	"k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
)

func APIServiceInfo(as *v1beta1.APIService) []string{

	service := ""
	available := ""
	for _, cd := range  as.Status.Conditions {
		if cd.Type == "Available" {
			service = cd.Reason
			available = string(cd.Status)
			break
		}
	}

	age := cobrautil.GetAge(as.CreationTimestamp.Time)

	data := []string{as.Name, service, available, age}

	return data
}
func PrintAPIService(body []byte) {
	as := v1beta1.APIService{}
	err := yaml.Unmarshal(body, &as)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	data := APIServiceInfo(&as)
	datas = append(datas, data)

	DrawAPIServiceTable(datas)

}
func PrintAPIServiceList(body []byte) {
	resourceStruct := v1beta1.APIServiceList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}

	for _, as := range resourceStruct.Items {
		data := APIServiceInfo(&as)
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

	DrawAPIServiceTable(datas)

}

func DrawAPIServiceTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "SERVICE", "AVAILABLE", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}