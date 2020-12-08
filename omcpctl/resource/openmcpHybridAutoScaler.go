package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	cobrautil "openmcp/openmcp/omcpctl/util"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"
	"os"
	"strconv"
	"strings"
)


func OpenMCPHybridAutoScalerInfo(ohas *resourcev1alpha1.OpenMCPHybridAutoScaler) []string{

	namespace := ohas.Namespace
	name := ohas.Name
	min_replica := strconv.Itoa(int(*ohas.Spec.HpaTemplate.Spec.MinReplicas))
	max_replica := strconv.Itoa(int(ohas.Spec.HpaTemplate.Spec.MaxReplicas))
	resources := []string{}

	for _, metric := range ohas.Spec.HpaTemplate.Spec.Metrics {
		resource := ""
		if metric.Type == "Resource" {
			resource = string(metric.Resource.Name) + "/"
			if metric.Resource.Target.Type == "Utilization" {
				resource += "averageUtilization(" + strconv.Itoa(int(*metric.Resource.Target.AverageUtilization)) + ")"
			} else if metric.Resource.Target.Type == "averagevalue" {
				resource += "averageValue(" + metric.Resource.Target.AverageValue.String()+ ")"
			} else if metric.Resource.Target.Type == "value"{
				resource += "value(" + metric.Resource.Target.Value.String()+ ")"
			}
		} else if metric.Type == "Object" {
			resource = metric.Object.Metric.Name + "(" + metric.Object.Target.Value.String() + " )"
		}

		resources = append(resources, resource)
	}

	reference_odeploy := ohas.Spec.HpaTemplate.Spec.ScaleTargetRef.Name
	vpa_mode := ohas.Spec.VpaMode
	policies := []string{}
	for _, policy := range ohas.Status.Policies{
		policyStr := policy.Type + "(" + strings.Join(policy.Value, "/") + ")"
		policies = append(policies, policyStr)
	}
	age := cobrautil.GetAge(ohas.CreationTimestamp.Time)
	data := []string{namespace, name,  min_replica, max_replica, strings.Join(resources, "\n"), reference_odeploy, vpa_mode, strings.Join(policies, "\n"), age}

	return data
}
func PrintOpenMCPHybridAutoScaler(body []byte) {
	od := resourcev1alpha1.OpenMCPHybridAutoScaler{}
	err := yaml.Unmarshal(body, &od)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	data := OpenMCPHybridAutoScalerInfo(&od)
	datas = append(datas, data)

	DrawOpenMCPHybridAutoScaleTable(datas)

}
func PrintOpenMCPHybridAutoScalerList(body []byte) {
	resourceStruct := resourcev1alpha1.OpenMCPHybridAutoScalerList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}

	for _, od := range resourceStruct.Items {
		data := OpenMCPHybridAutoScalerInfo(&od)
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

	DrawOpenMCPHybridAutoScaleTable(datas)

}

func DrawOpenMCPHybridAutoScaleTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "MIN_REPLICAS", "MAX_REPLICAS", "RESOURCES", "REFERENCE_ODEPLOY", "VPA_MODE", "POLICIES", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}