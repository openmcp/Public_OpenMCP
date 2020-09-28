package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	corev1 "k8s.io/api/core/v1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
)

func EventInfo(ev *corev1.Event) []string{
	lastSeen := cobrautil.GetAge(ev.LastTimestamp.Time)
	eventType := ev.Type
	reason := ev.Reason
	object := cobrautil.KindMap[ev.InvolvedObject.Kind]
	message := ev.Message


	data := []string{ev.Namespace, lastSeen, eventType, reason, object, message}

	return data
}
func PrintEvent(body []byte) {
	ev := corev1.Event{}
	err := yaml.Unmarshal(body, &ev)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	data := EventInfo(&ev)
	datas = append(datas, data)

	DrawEventTable(datas)

}
func PrintEventList(body []byte) {
	resourceStruct := corev1.EventList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}

	for _, ev := range resourceStruct.Items {
		data := EventInfo(&ev)
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

	DrawEventTable(datas)

}

func DrawEventTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "LAST SEEN", "TYPE", "REASON", "OBJECT", "MESSAGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}