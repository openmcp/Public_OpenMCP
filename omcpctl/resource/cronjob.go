package resource

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/ghodss/yaml"

	batchv1beta1 "k8s.io/api/batch/v1beta1"

	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
	"strconv"
)

func CronJobInfo(cj *batchv1beta1.CronJob) []string{
	schedule := cj.Spec.Schedule
	suspend := strconv.FormatBool(*cj.Spec.Suspend)
	active := "0"
	lastschedule := cobrautil.GetAge(cj.Status.LastScheduleTime.Time)

	age := cobrautil.GetAge(cj.CreationTimestamp.Time)

	data := []string{cj.Namespace, cj.Name, schedule, suspend, active, lastschedule, age}

	return data
}
func PrintCronJob(body []byte) {
	ds := batchv1beta1.CronJob{}
	err := yaml.Unmarshal(body, &ds)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := CronJobInfo(&ds)
	datas = append(datas, data)

	DrawCronJobTable(datas)

}
func PrintCronJobList(body []byte) {
	resourceStruct := batchv1beta1.CronJobList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, ds := range resourceStruct.Items {
		data := CronJobInfo(&ds)
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

	DrawCronJobTable(datas)

}

func DrawCronJobTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "SCHEDULE", "SUSPEND", "ACTIVE", "LAST SCHEDULE", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}