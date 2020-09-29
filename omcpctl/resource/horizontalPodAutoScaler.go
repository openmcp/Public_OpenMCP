package resource

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	autov1 "k8s.io/api/autoscaling/v1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"os"
	"strconv"
)

type hpaCurrentInfo struct {
	Type     string `json:"type"`
	Resource struct {
		Name                      string    `json:"name"`
		CurrentAverageValue       string    `json:"currentAverageValue"`
		CurrentAverageUtilization int       `json:"currentAverageUtilization"`
	} `json:"resource,omitempty"`
	Object struct {
		Target struct {
			Kind       string `json:"kind"`
			Name       string `json:"name"`
			APIVersion string `json:"apiVersion"`
		} `json:"target"`
		MetricName   string `json:"metricName"`
		CurrentValue string `json:"currentValue"`
	} `json:"object,omitempty"`
}

type hpaTargetInfo struct {
	Type     string `json:"type"`
	Resource struct {
		Name                        string    `json:"name"`
		TargetAverageValue          string    `json:"targetAverageValue"`
		TargetAverageUtilization    int       `json:"targetAverageUtilization"`
	} `json:"resource,omitempty"`
	Object struct {
		Target struct {
			Kind       string `json:"kind"`
			Name       string `json:"name"`
			APIVersion string `json:"apiVersion"`
		} `json:"target"`
		MetricName   string    `json:"metricName"`
		TargetValue  string    `json:"targetValue"`
	} `json:"object,omitempty"`
}

func HorizontalPodAutoscalerInfo(hpa *autov1.HorizontalPodAutoscaler) []string{

	var current_metrics []hpaCurrentInfo
	var target_metrics []hpaTargetInfo
	targets := ""
	num := 0

	_ = json.Unmarshal([]byte(hpa.Annotations["autoscaling.alpha.kubernetes.io/current-metrics"]), &current_metrics)
	_ = json.Unmarshal([]byte(hpa.Annotations["autoscaling.alpha.kubernetes.io/metrics"]), &target_metrics)

	if current_metrics != nil {

	}else {

	}

	if hpa.Spec.TargetCPUUtilizationPercentage != nil {
		num = 1
		if current_metrics != nil {
			for i := range current_metrics {
				if current_metrics[i].Resource.Name == "cpu" {
					targets += strconv.Itoa(current_metrics[i].Resource.CurrentAverageUtilization) + "%/" + strconv.Itoa(int(*hpa.Spec.TargetCPUUtilizationPercentage)) + "%"
					break
				}
			}
		}else {
			targets += "<unknown>/" + strconv.Itoa(int(*hpa.Spec.TargetCPUUtilizationPercentage)) + "%"
		}
	}

	for i := range target_metrics {
		currentAverageValue := ""
		currentAverageUtilization := 0
		targetAverageValue := ""
		targetAverageUtilization := 0

		if num >= 2 {
			l := 0
			if hpa.Spec.TargetCPUUtilizationPercentage != nil {
				l = len(target_metrics) - 1
			}else {
				l = len(target_metrics) - 2
			}
			targets += " + " + strconv.Itoa(l) + " more..."
			break
		}

		if current_metrics != nil {
			if current_metrics[i].Type == "Resource" {
				currentAverageValue = current_metrics[i].Resource.CurrentAverageValue
				currentAverageUtilization = current_metrics[i].Resource.CurrentAverageUtilization
				targetAverageValue = target_metrics[i].Resource.TargetAverageValue
				targetAverageUtilization = target_metrics[i].Resource.TargetAverageUtilization
			} else if current_metrics[i].Type == "Object" {
				currentAverageValue = current_metrics[i].Object.CurrentValue
				targetAverageValue = target_metrics[i].Object.TargetValue
			}

			if targetAverageValue != "" {
				if num > 0 {
					targets += ", "
				}
				targets += currentAverageValue + "/" + targetAverageValue
				num += 1
			} else if targetAverageUtilization > 0 {
				if num > 0 {
					targets += ", "
				}
				targets += strconv.Itoa(currentAverageUtilization) + "%/" + strconv.Itoa(targetAverageUtilization) + "%"
				num += 1
			}

		}else {
			if target_metrics[i].Type == "Resource" {
				targetAverageValue = target_metrics[i].Resource.TargetAverageValue
				targetAverageUtilization = target_metrics[i].Resource.TargetAverageUtilization
			} else if target_metrics[i].Type == "Object" {
				targetAverageValue = target_metrics[i].Object.TargetValue
			}

			if targetAverageValue != "" {
				if num > 0 {
					targets += ", "
				}
				targets += "<unknown>/" + targetAverageValue
				num += 1
			} else if targetAverageUtilization > 0 {
				if i > 0 {
					targets += ", "
				}
				targets += "<unknown>/" + strconv.Itoa(targetAverageUtilization) + "%"
				num += 1
			}
		}
	}

	reference := hpa.Spec.ScaleTargetRef.Name + "("+hpa.Spec.ScaleTargetRef.Kind+")"
	minpods := strconv.Itoa(int(*hpa.Spec.MinReplicas))
	maxpods := strconv.Itoa(int(hpa.Spec.MaxReplicas))
	replicas := strconv.Itoa(int(hpa.Status.CurrentReplicas))

	age := cobrautil.GetAge(hpa.CreationTimestamp.Time)

	data := []string{hpa.Namespace, hpa.Name, reference, targets, minpods, maxpods, replicas, age}

	return data
}

func PrintHorizontalPodAutoscaler(body []byte) {
	hpa := autov1.HorizontalPodAutoscaler{}
	err := yaml.Unmarshal(body, &hpa)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}


	data := HorizontalPodAutoscalerInfo(&hpa)
	datas = append(datas, data)

	DrawHorizontalPodAutoscalerTable(datas)

}

func PrintHorizontalPodAutoscalerList(body []byte) {
	resourceStruct := autov1.HorizontalPodAutoscalerList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		panic(err.Error())
	}
	datas := [][]string{}

	for _, hpa := range resourceStruct.Items {
		data := HorizontalPodAutoscalerInfo(&hpa)
		datas = append(datas, data)
	}

	if len(resourceStruct.Items) == 0 {
		ns := "default"
		if cobrautil.Option_namespace != "" {
			ns = cobrautil.Option_namespace
		}
		fmt.Println("No resources found in "+ ns +" namespace.")
		return
	}

	DrawHorizontalPodAutoscalerTable(datas)

}

func DrawHorizontalPodAutoscalerTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "NAME", "REFERENCE", "TARGETS", "MINPODS", "MAXPODS", "REPLICAS", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}