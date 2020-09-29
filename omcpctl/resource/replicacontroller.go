package resource

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/ghodss/yaml"

	corev1 "k8s.io/api/core/v1"

	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
	"strconv"
)

func ReplicationControllerInfo(rc *corev1.ReplicationController) []string{

	desired := strconv.Itoa(int(*rc.Spec.Replicas))
	current := strconv.Itoa(int(rc.Status.Replicas))
	ready := strconv.Itoa(int(rc.Status.ReadyReplicas))

	age := cobrautil.GetAge(rc.CreationTimestamp.Time)

	data := []string{rc.Namespace, rc.Name, desired, current, ready, age}

	return data
}
func PrintReplicationController(body []byte) {
	rc := corev1.ReplicationController{}
	err := yaml.Unmarshal(body, &rc)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	data := ReplicationControllerInfo(&rc)
	datas = append(datas, data)

	DrawReplicationControllerTable(datas)

}
func PrintReplicationControllerList(body []byte) {
	resourceStruct := corev1.ReplicationControllerList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}

	for _, rc := range resourceStruct.Items {
		data := ReplicationControllerInfo(&rc)
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

	DrawReplicationControllerTable(datas)

}

func DrawReplicationControllerTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "DESIRED", "CURRENT", "READY", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}