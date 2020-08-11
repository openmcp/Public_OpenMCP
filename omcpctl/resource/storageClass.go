package resource

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/ghodss/yaml"
	storagev1 "k8s.io/api/storage/v1"

	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
	"strconv"
)

func StorageClassInfo(sc *storagev1.StorageClass) []string{
	provisioner := sc.Provisioner
	reclaimpolicy := string(*sc.ReclaimPolicy)
	volumebindingmode := string(*sc.VolumeBindingMode)
	allowvolumeexpansion := "false"

	if sc.AllowVolumeExpansion != nil {
		allowvolumeexpansion = strconv.FormatBool(*sc.AllowVolumeExpansion)
	}

	age := cobrautil.GetAge(sc.CreationTimestamp.Time)

	data := []string{sc.Name, provisioner, reclaimpolicy, volumebindingmode, allowvolumeexpansion, age}

	return data
}
func PrintStorageClass(body []byte) {
	sc := storagev1.StorageClass{}
	err := yaml.Unmarshal(body, &sc)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := StorageClassInfo(&sc)
	datas = append(datas, data)

	DrawStorageClassTable(datas)

}
func PrintStorageClassList(body []byte) {
	resourceStruct := storagev1.StorageClassList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, sc := range resourceStruct.Items {
		data := StorageClassInfo(&sc)
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

	DrawStorageClassTable(datas)

}

func DrawStorageClassTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "PROVISIONER", "RECLAIMPOLICY", "VOLUMEBINDINGMODE", "ALLOWVOLUMEEXPANSION", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}