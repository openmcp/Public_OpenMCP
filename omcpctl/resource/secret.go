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

func SecretInfo(sec *corev1.Secret) []string{

	secType := string(sec.Type)
	dataNum := strconv.Itoa(len(sec.Data))
	age := cobrautil.GetAge(sec.CreationTimestamp.Time)

	data := []string{sec.Namespace, sec.Name, secType, dataNum, age}

	return data
}
func PrintSecret(body []byte) {
	sec := corev1.Secret{}
	err := yaml.Unmarshal(body, &sec)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := SecretInfo(&sec)
	datas = append(datas, data)

	DrawSecretTable(datas)

}
func PrintSecretList(body []byte) {
	resourceStruct := corev1.SecretList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, sec := range resourceStruct.Items {
		data := SecretInfo(&sec)
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

	DrawSecretTable(datas)

}

func DrawSecretTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "TYPE", "DATA", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}