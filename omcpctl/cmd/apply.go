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
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"openmcp/openmcp/omcpctl/apiServerMethod"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"strings"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

openmcpctl apply -f <FILENAME>
openmcpctl apply -f <FILENAME> --context <CLUSTERNAME>`,
	Run: func(cmd *cobra.Command, args []string) {
		ApplyResource(args)
	},
}

func ApplyResource(args []string){
	if cobrautil.Option_file == "" {
		fmt.Println("-f option needed")
		return
	}
	filenameList := cobrautil.GetFileNameList()


	for _, filename := range filenameList {
		yamlFile, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}


		err = RunApply(yamlFile)

		if err != nil {
			fmt.Println(err)
		}
	}

}
func PrepareApply(yamlFile []byte) (string, io.Reader) {
	var metainfo cobrautil.MetaInfo

	err := yaml.Unmarshal(yamlFile, &metainfo)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Value: %#v\n", metainfo.Kind)
	fmt.Printf("Value: %#v\n", metainfo.Metadata.Name)


	LINK := cobrautil.ApplyLinkParser(&metainfo)
	fmt.Println(LINK)

	body := strings.NewReader(string(yamlFile))
	return LINK, body
}
func RunApply(yamlFile []byte) error {
	LINK, body := PrepareApply(yamlFile)

	msg, err := apiServerMethod.PutAPIServer(LINK, body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	metainfo, err := getMetaInfo(msg)
	fmt.Println(string(msg))

	if err != nil {
		fmt.Println(err)
		return err
	}
	if metainfo.Message != "" {
		if metainfo.Reason == "NotFound"{
			err = RunCreate(yamlFile)
		} else {
			fmt.Println(metainfo.Message)
		}

	} else {
		fmt.Println(cobrautil.KindMap[metainfo.Kind] + " \""+metainfo.Metadata.Name+"\" configured")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	applyCmd.Flags().StringVarP(&cobrautil.Option_file, "file","f", "", "input a option")
	applyCmd.Flags().StringVarP(&cobrautil.Option_context, "context","c", "", "input a option")
	applyCmd.Flags().StringVarP(&cobrautil.Option_namespace, "namespace","n", "", "input a option")

}
