package resource

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/ghodss/yaml"
	appsv1 "k8s.io/api/apps/v1"
	"strings"

	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
	"strconv"
)

func DaemonSetInfo(ds *appsv1.DaemonSet) []string{
	desired := strconv.Itoa(int(ds.Status.DesiredNumberScheduled))
	current := strconv.Itoa(int(ds.Status.CurrentNumberScheduled))

	ready := strconv.Itoa(int(ds.Status.NumberReady))

	uptodate := strconv.Itoa(int(ds.Status.UpdatedNumberScheduled))
	available := strconv.Itoa(int(ds.Status.NumberAvailable))

	nodeselector := []string{}
	for k, v := range ds.Spec.Template.Spec.NodeSelector {
		nodeselector = append(nodeselector, k+"="+v)
	}


	age := cobrautil.GetAge(ds.CreationTimestamp.Time)

	data := []string{ds.Namespace, ds.Name, desired, current, ready, uptodate, available, strings.Join(nodeselector, ","), age}

	return data
}
func PrintDaemonSet(body []byte) {
	ds := appsv1.DaemonSet{}
	err := yaml.Unmarshal(body, &ds)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := DaemonSetInfo(&ds)
	datas = append(datas, data)

	DrawDaemonSetTable(datas)

}
func PrintDaemonSetList(body []byte) {
	resourceStruct := appsv1.DaemonSetList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, ds := range resourceStruct.Items {
		data := DaemonSetInfo(&ds)
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

	DrawDaemonSetTable(datas)

}

func DrawDaemonSetTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "DESIRED", "CURRENT", "READY", "UP-TO-DATE", "AVAILABLE", "NODE SELECTOR", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}