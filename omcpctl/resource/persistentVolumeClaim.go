package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	corev1 "k8s.io/api/core/v1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
)

func PersistentVolumeClaimInfo(pvc *corev1.PersistentVolumeClaim) []string{

	status := string(pvc.Status.Phase)
	volume := pvc.Spec.VolumeName
	capacity := pvc.Status.Capacity.Storage().String()
	if capacity == "0"{
		capacity = ""
	}

	accessmodes :=""
	if  len(pvc.Status.AccessModes) != 0 {
		accessmodes = string(pvc.Status.AccessModes[0])
	}
	storageclass := *pvc.Spec.StorageClassName

	age := cobrautil.GetAge(pvc.CreationTimestamp.Time)

	data := []string{pvc.Namespace, pvc.Name, status, volume, capacity, accessmodes, storageclass, age}

	return data
}
func PrintPersistentVolumeClaim(body []byte) {
	pvc := corev1.PersistentVolumeClaim{}
	err := yaml.Unmarshal(body, &pvc)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	data := PersistentVolumeClaimInfo(&pvc)
	datas = append(datas, data)

	DrawPersistentVolumeClaimTable(datas)

}
func PrintPersistentVolumeClaimList(body []byte) {
	resourceStruct := corev1.PersistentVolumeClaimList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}

	for _, pvc := range resourceStruct.Items {
		data := PersistentVolumeClaimInfo(&pvc)
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

	DrawPersistentVolumeClaimTable(datas)

}

func DrawPersistentVolumeClaimTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "STATUS", "VOLUME", "CAPACITY", "ACCESS", "MODES", "STORAGECLASS", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}