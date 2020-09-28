package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	corev1 "k8s.io/api/core/v1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
)

func PersistentVolumeInfo(pv *corev1.PersistentVolume) []string{
	capacity := pv.Spec.Capacity.Storage().String()
	accessModes :=""
	if  len(pv.Spec.AccessModes) != 0 {
		accessModes = string(pv.Spec.AccessModes[0])
	}
	reclaimPolicy := string(pv.Spec.PersistentVolumeReclaimPolicy)

	status := string(pv.Status.Phase)
	claim := ""
	if pv.Spec.ClaimRef != nil {
		claim = pv.Spec.ClaimRef.Namespace +"/"+ pv.Spec.ClaimRef.Name
	}

	storageClass := pv.Spec.StorageClassName
	reason := ""

	age := cobrautil.GetAge(pv.CreationTimestamp.Time)

	data := []string{pv.Name, capacity, accessModes, reclaimPolicy, status, claim, storageClass, reason, age}

	return data
}
func PrintPersistentVolume(body []byte) {
	pv := corev1.PersistentVolume{}
	err := yaml.Unmarshal(body, &pv)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	data := PersistentVolumeInfo(&pv)
	datas = append(datas, data)

	DrawPersistentVolumeTable(datas)

}
func PrintPersistentVolumeList(body []byte) {
	resourceStruct := corev1.PersistentVolumeList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}

	for _, pv := range resourceStruct.Items {
		data := PersistentVolumeInfo(&pv)
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

	DrawPersistentVolumeTable(datas)

}

func DrawPersistentVolumeTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "CAPACITY", "ACCESS MODES", "RECLAIM POLICY", "STATUS", "CLAIM", "STORAGECLASS", "REASON", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}