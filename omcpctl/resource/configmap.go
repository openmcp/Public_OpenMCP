package resource

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/ghodss/yaml"
	corev1 "k8s.io/api/core/v1"

	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
	"strconv"
)

func ConfigMapInfo(cm *corev1.ConfigMap) []string{

	dataNum := strconv.Itoa(len(cm.Data))
	age := cobrautil.GetAge(cm.CreationTimestamp.Time)

	data := []string{cm.Namespace, cm.Name, dataNum, age}

	return data
}
func PrintConfigMap(body []byte) {
	cm := corev1.ConfigMap{}
	err := yaml.Unmarshal(body, &cm)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := ConfigMapInfo(&cm)
	datas = append(datas, data)

	DrawConfigMapTable(datas)

}
func PrintConfigMapList(body []byte) {
	resourceStruct := corev1.ConfigMapList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, cm := range resourceStruct.Items {
		data := ConfigMapInfo(&cm)
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

	DrawConfigMapTable(datas)

}

func DrawConfigMapTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "DATA", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}