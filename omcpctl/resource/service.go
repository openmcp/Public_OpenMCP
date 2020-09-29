package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	corev1 "k8s.io/api/core/v1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
	"strconv"
	"strings"
)

func ServiceInfo(svc corev1.Service) []string{
	svctype := string(svc.Spec.Type)

	clusterIp  := svc.Spec.ClusterIP

	externalIp := ""
	if svctype != "LoadBalancer" {
		externalIp = "<none>"
	} else {
		externalIps := []string{}
		if len(svc.Status.LoadBalancer.Ingress) != 0 {
			for _, ing := range svc.Status.LoadBalancer.Ingress {
				externalIps = append(externalIps, ing.IP)
			}
		}

		if len(externalIps) == 0 {
			externalIp = "<pending>"
		} else {
			externalIp = strings.Join(externalIps,",")
		}
	}

	ports := []string{}
	for _, p := range svc.Spec.Ports {
		port := ""
		if svctype == "ClusterIP"{
			port = strconv.Itoa(int(p.Port)) + "/"+ string(p.Protocol)
		} else {
			port = strconv.Itoa(int(p.Port)) + ":" + strconv.Itoa(int(p.NodePort)) + "/"+ string(p.Protocol)
		}
		ports = append(ports, port)
	}

	age := cobrautil.GetAge(svc.ObjectMeta.CreationTimestamp.Time)


	data := []string{svc.Namespace, svc.Name, svctype, clusterIp, externalIp, strings.Join(ports,","), age}

	return data
}
func PrintService(body []byte) {
	svc := corev1.Service{}
	err := yaml.Unmarshal(body, &svc)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}

	data := ServiceInfo(svc)
	datas = append(datas, data)

	drawServiceTable(datas)

}
func PrintServiceList(body []byte) {
	resourceStruct := corev1.ServiceList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}

	for _, svc := range resourceStruct.Items {
		data := ServiceInfo(svc)
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

	drawServiceTable(datas)

}
func drawServiceTable(datas [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "TYPE", "CLUSTER-IP", "EXTERNAL-IP", "PORT(S)", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}