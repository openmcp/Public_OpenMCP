package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	appsv1 "k8s.io/api/apps/v1"

	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
	"strconv"
)

func StatefulSetInfo(sts *appsv1.StatefulSet) []string{


	ready := strconv.Itoa(int(sts.Status.ReadyReplicas)) + "/" + strconv.Itoa(int(*sts.Spec.Replicas))

	age := cobrautil.GetAge(sts.CreationTimestamp.Time)

	data := []string{sts.Namespace, sts.Name, ready, age}

	return data
}
func PrintStatefulSet(body []byte) {
	sts := appsv1.StatefulSet{}
	err := yaml.Unmarshal(body, &sts)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := StatefulSetInfo(&sts)
	datas = append(datas, data)

	DrawStatefulSetTable(datas)

}
func PrintStatefulSetList(body []byte) {
	resourceStruct := appsv1.StatefulSetList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, sts := range resourceStruct.Items {
		data := StatefulSetInfo(&sts)
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

	DrawStatefulSetTable(datas)

}

func DrawStatefulSetTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "READY", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}