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
	"fmt"
	"github.com/ghodss/yaml"
	"openmcp/openmcp/omcpctl/apiServerMethod"
	cobrautil "openmcp/openmcp/omcpctl/util"

	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		getLogs(args)
	},
}
func getLogs(args []string) {
	if len(args) == 0 {
		fmt.Println("POD or TYPE/NAME is a required argument for the logs command")
		return
	}
	c := cobrautil.GetKubeConfig("/root/.kube/config")

	clusterContext := c.CurrentContext
	if  cobrautil.Option_context != "" {
		clusterContext = cobrautil.Option_context
	}

	resourceKind := "pod"
	resourceName := args[0]

	LINK := cobrautil.LogLinkParser(resourceKind, resourceName, clusterContext)
	fmt.Println(LINK)

	body, err := apiServerMethod.GetAPIServer(LINK)
	if err != nil {
		fmt.Println("error: the server doesn't have a resource type '" + resourceKind + "'")
		return
	}

	var prettyYaml map[string]interface{}
	err = yaml.Unmarshal(body, &prettyYaml)
	if err != nil {
		fmt.Println(string(body))
	} else {
		fmt.Println(prettyYaml["message"])

	}

}

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().StringVarP(&cobrautil.Option_context, "context","c","", "input a option")
	logsCmd.Flags().StringVarP(&cobrautil.Option_namespace, "namespace","n", "", "input a option")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
