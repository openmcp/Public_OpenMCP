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
	"github.com/spf13/cobra"
	"openmcp/openmcp/omcpctl/apiServerMethod"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"strings"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

omcpctl exec (POD | TYPE/NAME) [-c CONTAINER] [flags] -- COMMAND [args...] [options]`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		fmt.Println(cobrautil.Option_containerName)
		fmt.Println(cobrautil.Option_stdin)
		fmt.Println(cobrautil.Option_tty)

		ExecSendToAPIServer(args[0], args[1:], cobrautil.Option_context)
	},
}

func ExecSendToAPIServer(podname string, command []string, context string){
	LINK := cobrautil.ExecLinkParser(podname, context)
	fmt.Println(LINK)

	parString := strings.Join(command, ",")
	ioR := strings.NewReader(parString)

	fmt.Println(parString)
	fmt.Println(ioR)

	msg, err := apiServerMethod.PostAPIServer(LINK, ioR)
	if err != nil {
		fmt.Println("fail - ",err)
		return
	}

	metainfo2, err := getMetaInfo(msg)

	if err != nil {
		return
	}
	if metainfo2.Message != "" {
		fmt.Println(metainfo2.Message)
	} else {
		//fmt.Println(podname + " exec")
	}
}

func init() {
	rootCmd.AddCommand(execCmd)

	execCmd.Flags().StringVarP(&cobrautil.Option_context, "context","x", "", "input a option")
	execCmd.Flags().StringVarP(&cobrautil.Option_namespace, "namespace","n", "", "input a option")

	execCmd.Flags().StringVarP(&cobrautil.Option_containerName, "container","c", "", "input a option")
	execCmd.Flags().BoolVarP(&cobrautil.Option_stdin, "stdin","i", false, "Default False")
	execCmd.Flags().BoolVarP(&cobrautil.Option_tty, "tty","t", false, "Default False")

}