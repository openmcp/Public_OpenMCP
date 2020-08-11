package resource

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/ghodss/yaml"
	policyv1beta1 "k8s.io/api/policy/v1beta1"

	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
	"strconv"
)

func PodDisruptionBudgetInfo(pdb *policyv1beta1.PodDisruptionBudget) []string{

	minAvailable := "N/A"
	if pdb.Spec.MinAvailable != nil {
		minAvailable = pdb.Spec.MinAvailable.String()
	}
	maxUnavailable := "N/A"
	if pdb.Spec.MaxUnavailable != nil {
		minAvailable = pdb.Spec.MaxUnavailable.String()
	}
	allowedDisruptions := strconv.Itoa(int(pdb.Status.DisruptionsAllowed))

	age := cobrautil.GetAge(pdb.CreationTimestamp.Time)

	data := []string{pdb.Namespace, pdb.Name, minAvailable, maxUnavailable, allowedDisruptions, age}

	return data
}
func PrintPodDisruptionBudget(body []byte) {
	pdb := policyv1beta1.PodDisruptionBudget{}
	err := yaml.Unmarshal(body, &pdb)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := PodDisruptionBudgetInfo(&pdb)
	datas = append(datas, data)

	DrawPodDisruptionBudgetTable(datas)

}
func PrintPodDisruptionBudgetList(body []byte) {
	resourceStruct := policyv1beta1.PodDisruptionBudgetList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, pdb := range resourceStruct.Items {
		data := PodDisruptionBudgetInfo(&pdb)
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

	DrawPodDisruptionBudgetTable(datas)

}

func DrawPodDisruptionBudgetTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "MIN AVAILABLE", "MAX UNAVAILABLE", "ALLOWED DISRUPTIONS", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}