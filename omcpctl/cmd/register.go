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
	cobrautil "openmcp/openmcp/omcpctl/util"
	"openmcp/openmcp/util"
	"os"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

omcpctl register openmcp <OPENMCPIP>
omcpctl register member  <OPENMCPIP>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 && args[0] == "openmcp" {
			registerASOpenMCP()
		} else if len(args) != 0 && args[0] == "member" {
			if args[1] == "" {
				fmt.Println("You Must Provide Cluster IP")
			} else {
				registerMemberToOpenMCP(args[1])
			}
		}
	},
}



func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func registerASOpenMCP() {
	c := cobrautil.GetOmcpctlConf("/var/lib/omcpctl/config.yaml")

	util.CmdExec("umount -l /mnt")
	defer util.CmdExec("umount -l /mnt")

	util.CmdExec("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")

	openmcpIP := cobrautil.GetOutboundIP()

	if fileExists("/mnt/openmcp/" + openmcpIP) {
		fmt.Println("Failed Register OpenMCP Master")
		fmt.Println("=> Already Registered OpenMCP :" + openmcpIP)
		return
	}

	util.CmdExec("mkdir /mnt/openmcp")
	util.CmdExec("mkdir /mnt/openmcp/" + openmcpIP)
	util.CmdExec("mkdir /mnt/openmcp/" + openmcpIP + "/master")
	util.CmdExec("mkdir /mnt/openmcp/" + openmcpIP + "/master/config")
	util.CmdExec("mkdir /mnt/openmcp/" + openmcpIP + "/master/pki")
	util.CmdExec("mkdir /mnt/openmcp/" + openmcpIP + "/members")
	util.CmdExec("mkdir /mnt/openmcp/" + openmcpIP + "/members/join")
	util.CmdExec("mkdir /mnt/openmcp/" + openmcpIP + "/members/unjoin")

	util.CmdExec("cp ~/.kube/config /mnt/openmcp/" + openmcpIP + "/master/config/config")
	util.CmdExec("cp /etc/kubernetes/pki/etcd/ca.crt /mnt/openmcp/" + openmcpIP + "/master/pki/ca.crt")
	util.CmdExec("cp /etc/kubernetes/pki/etcd/server.crt /mnt/openmcp/" + openmcpIP + "/master/pki/server.crt")
	util.CmdExec("cp /etc/kubernetes/pki/etcd/server.key /mnt/openmcp/" + openmcpIP + "/master/pki/server.key")

	//SSH Public Key Copy
	util.CmdExec("cat /mnt/ssh/id_rsa.pub >> /root/.ssh/authorized_keys")

	fmt.Println("Success OpenMCP Master Register '" + openmcpIP + "'")
	return

}

func registerMemberToOpenMCP(openmcpIP string) {
	c := cobrautil.GetOmcpctlConf("/var/lib/omcpctl/config.yaml")

	util.CmdExec("umount -l /mnt")
	defer util.CmdExec("umount -l /mnt")

	util.CmdExec("mount -t nfs " + c.NfsServer + ":/home/nfs/ /mnt")

	memberIP := cobrautil.GetOutboundIP()

	if !fileExists("/mnt/openmcp/" + openmcpIP + "/master") {
		fmt.Println("Failed Register '" + memberIP + "' in OpenMCP Master: " + openmcpIP)
		fmt.Println("=> Not Yet Register OpenMCP.")
		fmt.Println("=> First You Must be Input the Next Command in 'OpenMCP Master Server(" + openmcpIP + ")' : omcpctl register openmcp")
		return
	}

	if memberIP == openmcpIP {
		fmt.Println("Failed Register '" + memberIP + "' in OpenMCP Master: " + openmcpIP)
		fmt.Println("=> Can Not Self regist. [My_IP '" + memberIP + "', OpenMCP_IP '" + openmcpIP + "']")
		return
	}

	// Already Register
	if fileExists("/mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP) {
		fmt.Println("Failed Register '" + memberIP + "' in OpenMCP Master: " + openmcpIP)
		fmt.Println("=> Already Regist")
		return

	} else if fileExists("/mnt/openmcp/" + openmcpIP + "/members/join/" + memberIP) {
		fmt.Println("Failed Register '" + memberIP + "' in OpenMCP Master: " + openmcpIP)
		fmt.Println("=> Already Joined by OpenMCP '" + openmcpIP + "'")
		return

	} else {
		util.CmdExec("mkdir /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP)
		util.CmdExec("mkdir /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP + "/config")
		util.CmdExec("mkdir /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP + "/pki")

		util.CmdExec("cp ~/.kube/config /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP + "/config/config")
		util.CmdExec("cp /etc/kubernetes/pki/etcd/ca.crt /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP + "/pki/ca.crt")
		util.CmdExec("cp /etc/kubernetes/pki/etcd/server.crt /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP + "/pki/server.crt")
		util.CmdExec("cp /etc/kubernetes/pki/etcd/server.key /mnt/openmcp/" + openmcpIP + "/members/unjoin/" + memberIP + "/pki/server.key")

		// SSH Public Key Copy
		util.CmdExec("cat /mnt/ssh/id_rsa.pub >> /root/.ssh/authorized_keys")

		fmt.Println("Success Register '" + memberIP + "' in OpenMCP Master: " + openmcpIP)
		return
	}
}

func init() {
	rootCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
