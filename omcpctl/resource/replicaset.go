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

func ReplicaSetInfo(rs *appsv1.ReplicaSet) []string{

	desired := strconv.Itoa(int(*rs.Spec.Replicas))
	current := strconv.Itoa(int(rs.Status.Replicas))
	ready := strconv.Itoa(int(rs.Status.ReadyReplicas))

	age := cobrautil.GetAge(rs.CreationTimestamp.Time)

	data := []string{rs.Namespace, rs.Name, desired, current, ready, age}

	return data
}
func PrintReplicaSet(body []byte) {
	rs := appsv1.ReplicaSet{}
	err := yaml.Unmarshal(body, &rs)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	data := ReplicaSetInfo(&rs)
	datas = append(datas, data)

	DrawReplicaSetTable(datas)

}
func PrintReplicaSetList(body []byte) {
	resourceStruct := appsv1.ReplicaSetList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}

	for _, rs := range resourceStruct.Items {
		data := ReplicaSetInfo(&rs)
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

	DrawReplicaSetTable(datas)

}

func DrawReplicaSetTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "DESIRED", "CURRENT", "READY", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}