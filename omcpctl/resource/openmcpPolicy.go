package resource

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/ghodss/yaml"

	cobrautil "openmcp/openmcp/omcpctl/util"
	ketiv1alpha1 "openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"

	"os"

	"strings"
)


func OpenMCPPolicyInfo(opol *ketiv1alpha1.OpenMCPPolicy) []string{

	namespace := opol.Namespace
	name := opol.Name
	status := opol.Spec.PolicyStatus
	policies := []string{}
	for _, pol := range opol.Spec.Template.Spec.Policies{
		policy := pol.Type + "("+ strings.Join(pol.Value,",") + ")"
		policies = append(policies, policy)
	}
	rangeApp := opol.Spec.RangeOfApplication
	targetCont := opol.Spec.Template.Spec.TargetController.Kind
	age := cobrautil.GetAge(opol.CreationTimestamp.Time)

	data := []string{namespace, name, status, strings.Join(policies,"\n"), rangeApp, targetCont, age}
	

	return data
}
func PrintOpenMCPPolicy(body []byte) {
	od := ketiv1alpha1.OpenMCPPolicy{}
	err := yaml.Unmarshal(body, &od)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := OpenMCPPolicyInfo(&od)
	datas = append(datas, data)

	DrawOpenMCPPolicyTable(datas)

}
func PrintOpenMCPPolicyList(body []byte) {
	resourceStruct := ketiv1alpha1.OpenMCPPolicyList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, od := range resourceStruct.Items {
		data := OpenMCPPolicyInfo(&od)
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

	DrawOpenMCPPolicyTable(datas)

}

func DrawOpenMCPPolicyTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "STATUS", "POLICIES", "RANGE_APP", "TARGET_CONT","AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}


//func GetOpenmcpPolicy(opolName string) {
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
//	opol, err := crdClient.OpenMCPPolicy(ns).Get(opolName, metav1.GetOptions{})
//	if err != nil && errors.IsNotFound(err) {
//		fmt.Println("No resources found in " + ns + " namespace.")
//		return
//	} else if err != nil {
//		fmt.Println(err)
//		return
//	}
//	if cobrautil.Option_filetype == "yaml"{
//		data, err := yaml.Marshal(&opol)
//		if err != nil {
//			log.Fatalf("Unmarshal: %v", err)
//		}
//		fmt.Println(string(data))
//		return
//	}
//	datas := [][]string{}
//	namespace := opol.Namespace
//	name := opol.Name
//	status := opol.Spec.PolicyStatus
//	policies := []string{}
//	for _, pol := range opol.Spec.Template.Spec.Policies{
//		policy := pol.Type + "("+ strings.Join(pol.Value,",") + ")"
//		policies = append(policies, policy)
//	}
//	rangeApp := opol.Spec.RangeOfApplication
//	targetCont := opol.Spec.Template.Spec.TargetController.Kind
//	age := ""
//
//	data := []string{namespace, name, status, strings.Join(policies,"\n"), rangeApp, targetCont, age}
//	datas = append(datas, data)
//
//	table := tablewriter.NewWriter(os.Stdout)
//	table.SetHeader([]string{"NS", "NAME", "STATUS", "POLICIES", "RANGE_APP", "TARGET_CONT","AGE"})
//	table.SetBorder(false)
//	table.AppendBulk(datas)
//	table.Render()
//}
//func GetOpenmcpPolicyList() {
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
//	opolList, err := crdClient.OpenMCPPolicy(ns).List(metav1.ListOptions{})
//	if err != nil && errors.IsNotFound(err) {
//		fmt.Println("No resources found in " + ns + " namespace.")
//		return
//	} else if err != nil {
//		fmt.Println(err)
//		return
//	}
//	if cobrautil.Option_filetype == "yaml"{
//		data, err := yaml.Marshal(&opolList)
//		if err != nil {
//			log.Fatalf("Unmarshal: %v", err)
//		}
//		fmt.Println(string(data))
//		return
//	}
//	datas := [][]string{}
//	for _, opol := range opolList.Items{
//		namespace := opol.Namespace
//		name := opol.Name
//		status := opol.Spec.PolicyStatus
//		policies := []string{}
//		for _, pol := range opol.Spec.Template.Spec.Policies{
//			policy := pol.Type + "("+ strings.Join(pol.Value,",") + ")"
//			policies = append(policies, policy)
//		}
//		rangeApp := opol.Spec.RangeOfApplication
//		targetCont := opol.Spec.Template.Spec.TargetController.Kind
//		age := ""
//
//		data := []string{namespace, name, status, strings.Join(policies,"\n"), rangeApp, targetCont, age}
//		datas = append(datas, data)
//	}
//
//	table := tablewriter.NewWriter(os.Stdout)
//	table.SetHeader([]string{"NS", "NAME", "STATUS", "POLICIES", "RANGE_APP", "TARGET_CONT","AGE"})
//	table.SetBorder(false)
//	table.AppendBulk(datas)
//	table.Render()
//}