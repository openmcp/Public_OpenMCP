/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"openmcp/openmcp/omcpctl/apiServerMethod"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

example)
omcpctl set openmcppolicys log-level Level 5
omcpctl set opol log-level Level 2
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("error: Required resource not specified.")
			return
		}
		resourceKind := args[0]

		if cobrautil.ResourceMap[resourceKind] == "openmcppolicys"{
			setPolicy(args)
		}
	},
}
func setPolicy(args []string){
	if len(args) < 4 {
		fmt.Println("error: Required resource not specified.")
		return
	}
	resourceName := args[1]
	policyName := args[2]
	value := args[3]

	cobrautil.Option_namespace = "openmcp"

	LINK := cobrautil.GetLinkParser("openmcppolicys", resourceName, "openmcp")

	body, err := apiServerMethod.GetAPIServer(LINK)
	if err != nil {
		fmt.Println("error: the server doesn't have a resource type 'openmcppolicys'")
		return
	}

	metainfo, err := getMetaInfo(body)
	LINK = cobrautil.ApplyLinkParser(&metainfo)

	policyInfo, err := bytearrayToPolicy(body)
	if err != nil {
		fmt.Println(err)
		return
	}
	find := false
	prevValue := ""
	for i, policy := range policyInfo.Spec.Template.Spec.Policies{
		if policy.Type == policyName {
			prevValue = policyInfo.Spec.Template.Spec.Policies[i].Value[0]
			policyInfo.Spec.Template.Spec.Policies[i].Value[0] = value
			find = true
			break
		}
	}

	if find {
		fmt.Println("[Set Policy Complete] ResourceName : "+ resourceName+ ", policyName : "+ policyName+ ", value : "+ prevValue+ " -> "+ value)
		requestByte, err := policyToBytearray(policyInfo)
		if err != nil {
			fmt.Println(err)
			return
		}
		body := bytes.NewReader(requestByte)
		_, err = apiServerMethod.PutAPIServer(LINK, body)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("Can not Find Policy")
	}





}

func bytearrayToPolicy(body []byte) (v1alpha1.OpenMCPPolicy, error){
	var policyinfo v1alpha1.OpenMCPPolicy
	err := yaml.Unmarshal(body, &policyinfo)
	return policyinfo, err
}
func policyToBytearray(policyinfo v1alpha1.OpenMCPPolicy) ([]byte, error){

	body, err := yaml.Marshal(&policyinfo)
	return body, err
}
func init() {
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
