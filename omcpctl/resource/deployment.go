package resource

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"

	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
	"strconv"
)

func DeploymentInfo(dep *appsv1.Deployment) []string{


	ready := strconv.Itoa(int(dep.Status.ReadyReplicas)) + "/" + strconv.Itoa(int(dep.Status.Replicas))

	uptodate := strconv.Itoa(int(dep.Status.ReadyReplicas))
	available := strconv.Itoa(int(dep.Status.AvailableReplicas))


	age := cobrautil.GetAge(dep.CreationTimestamp.Time)

	data := []string{dep.Namespace, dep.Name, ready, uptodate, available, age}

	return data
}
func PrintDeployment(body []byte) {
	dep := appsv1.Deployment{}
	err := yaml.Unmarshal(body, &dep)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := DeploymentInfo(&dep)
	datas = append(datas, data)

	DrawDeploymentTable(datas)

}
func PrintDeploymentList(body []byte) {
	resourceStruct := appsv1.DeploymentList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, dep := range resourceStruct.Items {
		data := DeploymentInfo(&dep)
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

	DrawDeploymentTable(datas)

}

func DrawDeploymentTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "READY", "UP-TO-DATE", "AVAILABLE", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}