package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	"strings"

	"k8s.io/api/extensions/v1beta1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
)

func IngressInfo(ing *v1beta1.Ingress) []string{

	hosts := []string{}

	for _, rule := range ing.Spec.Rules {
		hosts = append(hosts, rule.Host)
	}

	address := ing.Status.LoadBalancer.Ingress[0].IP
	ports := "80"


	age := cobrautil.GetAge(ing.CreationTimestamp.Time)

	data := []string{ing.Namespace, ing.Name, strings.Join(hosts, ","), address, ports, age}

	return data
}
func PrintIngress(body []byte) {
	ing := v1beta1.Ingress{}
	err := yaml.Unmarshal(body, &ing)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := IngressInfo(&ing)
	datas = append(datas, data)

	DrawIngressTable(datas)

}
func PrintIngressList(body []byte) {
	resourceStruct := v1beta1.IngressList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, ing := range resourceStruct.Items {
		data := IngressInfo(&ing)
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

	DrawIngressTable(datas)

}

func DrawIngressTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "HOSTS", "ADDRESS", "PORTS", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}