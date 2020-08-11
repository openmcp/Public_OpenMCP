package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	corev1 "k8s.io/api/core/v1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
)


func ResourceQuotaInfo(quota *corev1.ResourceQuota) []string{


	usedReqCpu := quota.Status.Used["requests.cpu"]
	usedReqMem := quota.Status.Used["requests.memory"]
	usedLimitCpu := quota.Status.Used["limits.cpu"]
	usedLimitMem := quota.Status.Used["limits.memory"]

	hardReqCpu := quota.Status.Hard["requests.cpu"]
	hardReqMem := quota.Status.Hard["requests.memory"]
	hardLimitCpu := quota.Status.Hard["limits.cpu"]
	hardLimitMem := quota.Status.Hard["limits.memory"]

	requestCpu := usedReqCpu.String() + "/" + hardReqCpu.String()
	requestMemory := usedReqMem.String() + "/" + hardReqMem.String()
	limitCpu := usedLimitCpu.String() + "/" + hardLimitCpu.String()
	limitMemory := usedLimitMem.String() + "/" + hardLimitMem.String()

	createadAt := quota.CreationTimestamp.Time.String()

	data := []string{quota.Namespace, quota.Name, requestCpu, requestMemory, limitCpu, limitMemory, createadAt}

	return data
}
func PrintResourceQuota(body []byte) {
	quota := corev1.ResourceQuota{}
	err := yaml.Unmarshal(body, &quota)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := ResourceQuotaInfo(&quota)
	datas = append(datas, data)

	DrawResourceQuotaTable(datas)

}
func PrintResourceQuotaList(body []byte) {
	resourceStruct := corev1.ResourceQuotaList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, quota := range resourceStruct.Items {
		data := ResourceQuotaInfo(&quota)
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

	DrawResourceQuotaTable(datas)

}

func DrawResourceQuotaTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "REQUEST CPU", "REQUEST MEMORY","LIMIT CPU","LIMIT MEMORY", "CREATED AT"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}