package resource

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/ghodss/yaml"

	cobrautil "openmcp/openmcp/omcpctl/util"
	ketiv1alpha1 "openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"

	"os"

	"strconv"
	"strings"
)


func OpenMCPHybridAutoScalerInfo(ohas *ketiv1alpha1.OpenMCPHybridAutoScaler) []string{

	namespace := ohas.Namespace
	name := ohas.Name
	min_replica := strconv.Itoa(int(*ohas.Spec.HpaTemplate.Spec.MinReplicas))
	max_replica := strconv.Itoa(int(ohas.Spec.HpaTemplate.Spec.MaxReplicas))
	resources := []string{}

	for _, metric := range ohas.Spec.HpaTemplate.Spec.Metrics {
		resource := string(metric.Resource.Name) + "/"
		if metric.Resource.Target.Type == "Utilization" {
			resource += "averageUtilization(" + strconv.Itoa(int(*metric.Resource.Target.AverageUtilization)) + ")"
		} else if metric.Resource.Target.Type == "averagevalue" {
			resource += "averageValue(" + metric.Resource.Target.AverageValue.String()+ ")"
		} else if metric.Resource.Target.Type == "value"{
			resource += "value(" + metric.Resource.Target.Value.String()+ ")"
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
	data := []string{namespace, name,  min_replica, max_replica,  strings.Join(resources, "\n"), reference_odeploy, vpa_mode, strings.Join(policies, "\n"), age}

	return data
}
func PrintOpenMCPHybridAutoScaler(body []byte) {
	od := ketiv1alpha1.OpenMCPHybridAutoScaler{}
	err := yaml.Unmarshal(body, &od)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := OpenMCPHybridAutoScalerInfo(&od)
	datas = append(datas, data)

	DrawOpenMCPHybridAutoScaleTable(datas)

}
func PrintOpenMCPHybridAutoScalerList(body []byte) {
	resourceStruct := ketiv1alpha1.OpenMCPHybridAutoScalerList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
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


//
//func GetOpenmcpHybridAutoScale(ohasName string) {
//	kubeconfig, _ := cobrautil.BuildConfigFromFlags("openmcp", "/root/.kube/config")
//	crdClient, _ := clientV1alpha1.NewForConfig(kubeconfig)
//	ns := ""
//	if cobrautil.Option_allnamespace {
//		ns = metav1.NamespaceAll
//	} else {
//		if cobrautil.Option_namespace != "" {
//			ns = cobrautil.Option_namespace
//		} else {
//			ns = "default"
//		}
//	}
//	ohas, err := crdClient.OpenMCPHybridAutoScaler(ns).Get(ohasName, metav1.GetOptions{})
//	if err != nil && errors.IsNotFound(err) {
//		fmt.Println("No resources found in " + ns + " namespace.")
//		return
//	} else if err != nil {
//		fmt.Println(err)
//		return
//	}
//	if cobrautil.Option_filetype == "yaml"{
//		data, err := yaml.Marshal(&ohas)
//		if err != nil {
//			log.Fatalf("Unmarshal: %v", err)
//		}
//		fmt.Println(string(data))
//		return
//	}
//	datas := [][]string{}
//	namespace := ohas.Namespace
//	name := ohas.Name
//	min_replica := strconv.Itoa(int(*ohas.Spec.HpaTemplate.Spec.MinReplicas))
//	max_replica := strconv.Itoa(int(ohas.Spec.HpaTemplate.Spec.MaxReplicas))
//	resources := []string{}
//
//	for _, metric := range ohas.Spec.HpaTemplate.Spec.Metrics {
//		resource := string(metric.Resource.Name) + "/"
//		if metric.Resource.Target.Type == "Utilization" {
//			resource += "averageUtilization(" + strconv.Itoa(int(*metric.Resource.Target.AverageUtilization)) + ")"
//		} else if metric.Resource.Target.Type == "averagevalue" {
//			resource += "averageValue(" + metric.Resource.Target.AverageValue.String()+ ")"
//		} else if metric.Resource.Target.Type == "value"{
//			resource += "value(" + metric.Resource.Target.Value.String()+ ")"
//		}
//
//		resources = append(resources, resource)
//	}
//
//	reference_odeploy := ohas.Spec.HpaTemplate.Spec.ScaleTargetRef.Name
//	vpa_mode := ohas.Spec.VpaMode
//	policies := []string{}
//	for _, policy := range ohas.Status.Policies{
//		policyStr := policy.Type + "(" + strings.Join(policy.Value, "/") + ")"
//
//		policies = append(policies, policyStr)
//	}
//	age := ""
//	data := []string{namespace, name,  min_replica, max_replica,  strings.Join(resources, "\n"), reference_odeploy, vpa_mode, strings.Join(policies, "\n"), age}
//	datas = append(datas, data)
//
//	table := tablewriter.NewWriter(os.Stdout)
//	table.SetHeader([]string{"NS", "NAME", "MIN_REPLICAS", "MAX_REPLICAS", "RESOURCES", "REFERENCE_ODEPLOY", "VPA_MODE", "POLICIES", "AGE"})
//	table.SetBorder(false)
//	table.AppendBulk(datas)
//	table.Render()
//
//}
//func GetOpenmcpHybridAutoScaleList() {
//	kubeconfig, _ := cobrautil.BuildConfigFromFlags("openmcp", "/root/.kube/config")
//	crdClient, _ := clientV1alpha1.NewForConfig(kubeconfig)
//	ns := ""
//	if cobrautil.Option_allnamespace {
//		ns = metav1.NamespaceAll
//	} else {
//		if cobrautil.Option_namespace != "" {
//			ns = cobrautil.Option_namespace
//		} else {
//			ns = "default"
//		}
//	}
//	ohasList, err := crdClient.OpenMCPHybridAutoScaler(ns).List(metav1.ListOptions{})
//	if err != nil && errors.IsNotFound(err) {
//		fmt.Println("No Resource exist")
//		return
//	}
//	if len(ohasList.Items) == 0 {
//		fmt.Println("No resources found in "+ns+" namespace.")
//		return
//	}
//	if cobrautil.Option_filetype == "yaml"{
//		data, err := yaml.Marshal(&ohasList)
//		if err != nil {
//			log.Fatalf("Unmarshal: %v", err)
//		}
//		fmt.Println(string(data))
//		return
//	}
//	datas := [][]string{}
//	for _, ohas := range ohasList.Items{
//		namespace := ohas.Namespace
//		name := ohas.Name
//		min_replica := strconv.Itoa(int(*ohas.Spec.HpaTemplate.Spec.MinReplicas))
//		max_replica := strconv.Itoa(int(ohas.Spec.HpaTemplate.Spec.MaxReplicas))
//		resources := []string{}
//
//		for _, metric := range ohas.Spec.HpaTemplate.Spec.Metrics {
//			resource := string(metric.Resource.Name) + "/"
//			if metric.Resource.Target.Type == "Utilization" {
//				resource += "averageUtilization(" + strconv.Itoa(int(*metric.Resource.Target.AverageUtilization)) + ")"
//			} else if metric.Resource.Target.Type == "averagevalue" {
//				resource += "averageValue(" + metric.Resource.Target.AverageValue.String()+ ")"
//			} else if metric.Resource.Target.Type == "value"{
//				resource += "value(" + metric.Resource.Target.Value.String()+ ")"
//			}
//
//			resources = append(resources, resource)
//		}
//
//		reference_odeploy := ohas.Spec.HpaTemplate.Spec.ScaleTargetRef.Name
//		vpa_mode := ohas.Spec.VpaMode
//		policies := []string{}
//		for _, policy := range ohas.Status.Policies{
//			policyStr := policy.Type + "(" + strings.Join(policy.Value, "/") + ")"
//
//			policies = append(policies, policyStr)
//		}
//		age := ""
//		data := []string{namespace, name,  min_replica, max_replica,  strings.Join(resources, "\n"), reference_odeploy, vpa_mode, strings.Join(policies, "\n"), age}
//		datas = append(datas, data)
//
//	}
//
//	table := tablewriter.NewWriter(os.Stdout)
//	table.SetHeader([]string{"NS", "NAME", "MIN_REPLICAS", "MAX_REPLICAS", "RESOURCES", "REFERENCE_ODEPLOY", "VPA_MODE", "POLICIES", "AGE"})
//	table.SetBorder(false)
//	table.AppendBulk(datas)
//	table.Render()
//}