package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	dnsv1alpha1 "openmcp/openmcp/apis/dns/v1alpha1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
)


func OpenMCPDNSEndpointInfo(ode *dnsv1alpha1.OpenMCPDNSEndpoint) [][]string{
	datas := [][]string{}

	namespace := ode.Namespace
	name := ode.Name
	if len(ode.Spec.Domains) == 0 {
		return datas
	}
	//domain := ode.Spec.Domains[0]

	age := cobrautil.GetAge(ode.CreationTimestamp.Time)

	for _, endpoint := range ode.Spec.Endpoints {
		dnsName := endpoint.DNSName
		targets := endpoint.Targets

		for _, target := range targets {

			data := []string{namespace, name, dnsName, target, age}
			datas = append(datas, data)
		}

	}

	return datas
}
func PrintOpenMCPDNSEndpoint(body []byte) {
	ode := dnsv1alpha1.OpenMCPDNSEndpoint{}
	err := yaml.Unmarshal(body, &ode)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	datas2 := OpenMCPDNSEndpointInfo(&ode)
	for _, data := range datas2 {
		datas = append(datas, data)
	}

	DrawOpenMCPDNSEndpointTable(datas)

}
func PrintOpenMCPDNSEndpointList(body []byte) {
	resourceStruct := dnsv1alpha1.OpenMCPDNSEndpointList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}

	for _, ode := range resourceStruct.Items {
		datas2 := OpenMCPDNSEndpointInfo(&ode)
		for _, data := range datas2 {
			datas = append(datas, data)
		}
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

	DrawOpenMCPDNSEndpointTable(datas)

}

func DrawOpenMCPDNSEndpointTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "DOMAIN", "IP", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}