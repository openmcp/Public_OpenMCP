package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
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
