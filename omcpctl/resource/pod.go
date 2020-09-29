package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	corev1 "k8s.io/api/core/v1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
	"strconv"
)

func PodInfo(pod *corev1.Pod) []string{
	readyCount := 0
	for _, cs := range pod.Status.ContainerStatuses{
		if cs.Ready {
			readyCount += 1
			break
		}
	}
	ready := strconv.Itoa(readyCount)+"/"+strconv.Itoa(len(pod.Spec.Containers))

	status := string(pod.Status.Phase)
	if pod.ObjectMeta.DeletionTimestamp != nil {
		status = "Terminating"
	} else {
		for _, cs := range pod.Status.ContainerStatuses{
			if cs.State.Waiting != nil {
				status = cs.State.Waiting.Reason
			}
		}
	}

	restart := 0
	for _, cs := range pod.Status.ContainerStatuses{
		restart += int(cs.RestartCount)
	}

	age := cobrautil.GetAge(pod.CreationTimestamp.Time)

	data := []string{pod.Namespace, pod.Name, ready,
		status, strconv.Itoa(restart), age}

	return data
}
func PrintPod(body []byte) {
	pod := corev1.Pod{}
	err := yaml.Unmarshal(body, &pod)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	data := PodInfo(&pod)
	datas = append(datas, data)

	DrawPodTable(datas)

}
func PrintPodList(body []byte) {
	resourceStruct := corev1.PodList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}

	for _, pod := range resourceStruct.Items {
		data := PodInfo(&pod)
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

	DrawPodTable(datas)

}

func DrawPodTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "READY", "STATUS", "RESTARTS", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}