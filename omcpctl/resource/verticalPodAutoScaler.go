package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta2"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
)


func VerticalPodAutoscalerInfo(vpa *v1beta2.VerticalPodAutoscaler) []string{


	reference := vpa.Spec.TargetRef.Name + "(" + vpa.Spec.TargetRef.Kind + ")"
	lower_bound := ""
	target := ""
	uncapped_bound := ""
	upper_bound := ""
	if vpa.Status.Recommendation != nil {
		for _, cr := range vpa.Status.Recommendation.ContainerRecommendations{
			_ = cr
			lower_bound += cr.ContainerName +"(cpu:" + cr.LowerBound.Cpu().String() + "/mem: "+ cr.LowerBound.Memory().String()+")\n"
			target += cr.ContainerName +"(cpu:" + cr.Target.Cpu().String() + "/mem: "+ cr.Target.Memory().String()+")\n"
			uncapped_bound += cr.ContainerName +"(cpu:" + cr.UncappedTarget.Cpu().String() + "/mem: "+ cr.UncappedTarget.Memory().String()+")\n"
			upper_bound += cr.ContainerName +"(cpu:" + cr.UpperBound.Cpu().String() + "/mem: "+ cr.UpperBound.Memory().String()+")\n"
		}
	}
	age := cobrautil.GetAge(vpa.CreationTimestamp.Time)

	data := []string{vpa.Namespace, vpa.Name, reference, lower_bound, target, uncapped_bound, upper_bound, age}

	return data
}

func PrintVerticalPodAutoscaler(body []byte) {
	vpa := v1beta2.VerticalPodAutoscaler{}
	err := yaml.Unmarshal(body, &vpa)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	data := VerticalPodAutoscalerInfo(&vpa)
	datas = append(datas, data)

	DrawVerticalPodAutoscalerTable(datas)

}

func PrintVerticalPodAutoscalerList(body []byte) {
	resourceStruct := v1beta2.VerticalPodAutoscalerList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}

	for _, vpa := range resourceStruct.Items {
		data := VerticalPodAutoscalerInfo(&vpa)
		datas = append(datas, data)
	}

	if len(resourceStruct.Items) == 0 {
		ns := "default"
		if cobrautil.Option_namespace != "" {
			ns = cobrautil.Option_namespace
		}
		fmt.Println("No resources found in "+ ns +" namespace.")
		return
	}

	DrawVerticalPodAutoscalerTable(datas)

}

func DrawVerticalPodAutoscalerTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "REFERENCE", "LOWER BOUND", "TARGET", "UNCAPPED TARGET", "UPPER BOUND", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}