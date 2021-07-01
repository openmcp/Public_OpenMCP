package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	batchv1 "k8s.io/api/batch/v1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
	"strconv"
)

func JobInfo(job *batchv1.Job) []string{

	numOfComplete := 0
	for _, c := range job.Status.Conditions {
		if c.Type == "Complete" {
			numOfComplete += 1
		}
	}
	completions := strconv.Itoa(numOfComplete) + "/" + strconv.Itoa(int(*job.Spec.Completions))

	age := cobrautil.GetAge(job.CreationTimestamp.Time)
	duration := age
	if job.Status.CompletionTime != nil {
		duration = cobrautil.GetDuration(job.CreationTimestamp.Time, job.Status.CompletionTime.Time)
	}


	data := []string{job.Namespace, job.Name, completions, duration, age}

	return data
}
func PrintJob(body []byte) {
	job := batchv1.Job{}
	err := yaml.Unmarshal(body, &job)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	data := JobInfo(&job)
	datas = append(datas, data)

	DrawJobTable(datas)

}
func PrintJobList(body []byte) {
	resourceStruct := batchv1.JobList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}

	for _, job := range resourceStruct.Items {
		data := JobInfo(&job)
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

	DrawJobTable(datas)

}

func DrawJobTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "COMPLETIONS", "DURATION", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}