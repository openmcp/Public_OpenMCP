package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	corev1 "k8s.io/api/core/v1"
	"strconv"
	"strings"

	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
)

func EndpointsInfo(ep *corev1.Endpoints) []string{

	endpoints := []string{}
	endpointsIps := []string{}
	endpointsPorts := []string{}
	for _, subset := range ep.Subsets {
		for _, address := range subset.Addresses {
			endpointsIps = append(endpointsIps, address.IP)
		}

		for _, port := range subset.Ports {
			endpointsPorts = append(endpointsPorts, strconv.Itoa(int(port.Port)))
		}
	}

	for _, endpointIp := range endpointsIps {
		for _, endpointsPort := range endpointsPorts {
			endpoints = append(endpoints, endpointIp+":"+endpointsPort)

			if len(endpoints) == 3 {
				break
			}
		}
		if len(endpoints) == 3 {
			break
		}
	}
	if len(endpointsIps) * len(endpointsPorts) -3 >= 1 {
		endpoints = append(endpoints, "+ "+strconv.Itoa(len(endpointsIps) * len(endpointsPorts) -3 ) + " more ...")
	}

	age := cobrautil.GetAge(ep.CreationTimestamp.Time)

	data := []string{ep.Namespace, ep.Name, strings.Join(endpoints,","), age}

	return data
}
func PrintEndpoints(body []byte) {
	ep := corev1.Endpoints{}
	err := yaml.Unmarshal(body, &ep)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := EndpointsInfo(&ep)
	datas = append(datas, data)

	DrawEndpointsTable(datas)

}
func PrintEndpointsList(body []byte) {
	resourceStruct := corev1.EndpointsList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, ep := range resourceStruct.Items {
		data := EndpointsInfo(&ep)
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

	DrawEndpointsTable(datas)

}

func DrawEndpointsTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "ENDPOINTS", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}