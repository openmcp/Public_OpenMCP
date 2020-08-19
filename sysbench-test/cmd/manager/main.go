/*
Copyright 2018 The Multicluster-Controller Authors.

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

package main

import (
	"fmt"
	"net/http"
	"openmcp/openmcp/util"
)
type HttpManager struct{
	HTTPServer_IP string
	HTTPServer_PORT string

	CpuBenchStatus string
	MemoryBenchStatus string
	DiskBenchStatus string
	NetworkBenchStatus string

}


var START_CPU_CMD string = "sysbench-cpu --test=cpu --cpu-max-prime=10000 run"
var START_MEMORY_CMD string = "sysbench-memory --test=memory --memory-block-size=100G --memory-total-size=1000G run"
var START_DISK_CMD string = "sysbench-disk --test=fileio --file-total-size=100G prepare"
var START_NETWORK_CMD string = ""


var STOP_CPU_CMD string =  "kill -9 `ps -ef | grep 'sysbench-cpu' | awk 'NR==1{print $2}'`"
var STOP_MEMORY_CMD string =  "kill -9 `ps -ef | grep 'sysbench-memory' | awk 'NR==1{print $2}'`"
var STOP_DISK_CMD string =  "kill -9 `ps -ef | grep 'sysbench-disk' | awk 'NR==1{print $2}'`"
var STOP_NETWORK_CMD string =  "kill -9 `ps -ef | grep 'sysbench-network' | awk 'NR==1{print $2}'`"



func (h *HttpManager) help(w http.ResponseWriter, r *http.Request){
	fmt.Println("Connect Service")

	w.Write([]byte("OpenMCP Service Response\n"))
	//w.Write([]byte("http://"+h.HTTPServer_IP+":"+h.HTTPServer_PORT+"/etcd/backup/start?time=5\n"))
	//w.Write([]byte("http://"+h.HTTPServer_IP+":"+h.HTTPServer_PORT+"/etcd/backup/stop\n"))
	//w.Write([]byte("http://"+h.HTTPServer_IP+":"+h.HTTPServer_PORT+"/etcd/backup/status\n"))
	//w.Write([]byte("http://"+h.HTTPServer_IP+":"+h.HTTPServer_PORT+"/etcd/get/snapshot/list\n"))
	//w.Write([]byte("http://"+h.HTTPServer_IP+":"+h.HTTPServer_PORT+"/etcd/restore/cluster?cluster=CLUSTER_NAME&file=SNAPSHOT.db\n"))
	//w.Write([]byte("http://"+h.HTTPServer_IP+":"+h.HTTPServer_PORT+"/etcd/restore/all?file=SNAPSHOT.db\n"))


}
func (h *HttpManager)sysbench_daemon(resourceType string) {
	for {
			cmdStr := ""
			if resourceType == "cpu" && h.CpuBenchStatus == "START"{
				cmdStr = START_CPU_CMD
			} else if resourceType == "memory" && h.MemoryBenchStatus == "START"{
				cmdStr = START_MEMORY_CMD
			} else if resourceType == "disk" && h.DiskBenchStatus == "START"{
				cmdStr = START_DISK_CMD
			} else if resourceType == "network" && h.NetworkBenchStatus == "START"{
				cmdStr = START_NETWORK_CMD
			} else {
				return
			}
			fmt.Println(cmdStr)
			_, err := util.CmdExec(cmdStr)
			if err != nil {
				fmt.Println("Err !", err)
			}
	}

}
func (h * HttpManager) printStatus(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("CPU Status : "+ h.CpuBenchStatus +"\n"))
	w.Write([]byte("Memory Status : "+ h.MemoryBenchStatus +"\n"))
	w.Write([]byte("Disk Status : "+ h.DiskBenchStatus +"\n"))
	w.Write([]byte("Network Status : "+ h.NetworkBenchStatus +"\n"))
}

func (h *HttpManager) cpu_start(w http.ResponseWriter, r *http.Request){
	if h.CpuBenchStatus == "STOP" {
		h.CpuBenchStatus = "START"
		go h.sysbench_daemon("cpu")
	}
	h.printStatus(w, r)
}
func (h *HttpManager) memory_start(w http.ResponseWriter, r *http.Request){

	if h.MemoryBenchStatus == "STOP" {
		h.MemoryBenchStatus = "START"
		go h.sysbench_daemon("memory")
	}
	h.printStatus(w, r)
}
func (h *HttpManager) disk_start(w http.ResponseWriter, r *http.Request){

	if h.DiskBenchStatus == "STOP" {
		h.DiskBenchStatus = "START"
		go h.sysbench_daemon("disk")
	}
	h.printStatus(w, r)
}
func (h *HttpManager) network_start(w http.ResponseWriter, r *http.Request){

	if h.NetworkBenchStatus == "STOP" {
		h.NetworkBenchStatus = "START"
		go h.sysbench_daemon("network")
	}
	h.printStatus(w, r)
}

func (h *HttpManager) cpu_stop(w http.ResponseWriter, r *http.Request){
	if h.CpuBenchStatus == "START" {
		h.CpuBenchStatus = "STOP"
		_, err := util.CmdExec(STOP_CPU_CMD)
		if err != nil {
			fmt.Println("Err !", err)
		}

	}
	h.printStatus(w, r)
}
func (h *HttpManager) memory_stop(w http.ResponseWriter, r *http.Request){
	if h.MemoryBenchStatus == "START" {
		h.MemoryBenchStatus = "STOP"
		_, err := util.CmdExec(STOP_MEMORY_CMD)
		if err != nil {
			fmt.Println("Err !", err)
		}
	}
	h.printStatus(w, r)
}
func (h *HttpManager) disk_stop(w http.ResponseWriter, r *http.Request){
	if h.DiskBenchStatus == "START" {
		h.DiskBenchStatus = "STOP"
		_, err := util.CmdExec(STOP_DISK_CMD)
		if err != nil {
			fmt.Println("Err !", err)
		}
	}
	h.printStatus(w, r)
}
func (h *HttpManager) network_stop(w http.ResponseWriter, r *http.Request){
	if h.NetworkBenchStatus == "START" {
		h.NetworkBenchStatus = "STOP"
		_, err := util.CmdExec(STOP_NETWORK_CMD)
		if err != nil {
			fmt.Println("Err !", err)
		}
	}
	h.printStatus(w, r)
}


func main() {
	//HTTPServer_IP :=   os.Getenv("SERVER_IP")
	//HTTPServer_PORT := os.Getenv("SERVER_PORT")

	HTTPServer_PORT := "8080"

	httpManager := &HttpManager{
		HTTPServer_IP:  "",
		HTTPServer_PORT: HTTPServer_PORT,

		CpuBenchStatus:   "STOP",
		MemoryBenchStatus:   "STOP",
		NetworkBenchStatus:   "STOP",
		DiskBenchStatus:   "STOP",
	}

	handler := http.NewServeMux()

	handler.HandleFunc("/", httpManager.help)
	handler.HandleFunc("/cpu/start", httpManager.cpu_start)
	handler.HandleFunc("/memory/start", httpManager.memory_start)
	handler.HandleFunc("/disk/start", httpManager.disk_start)
	handler.HandleFunc("/network/start", httpManager.network_start)
	handler.HandleFunc("/cpu/stop", httpManager.cpu_stop)
	handler.HandleFunc("/memory/stop", httpManager.memory_stop)
	handler.HandleFunc("/disk/stop", httpManager.disk_stop)
	handler.HandleFunc("/network/stop", httpManager.network_stop)



	fmt.Println("Sysbench-test Server Start ! " + ":" + HTTPServer_PORT)
	server := &http.Server{Addr: ":" + HTTPServer_PORT, Handler: handler}
	err := server.ListenAndServe()

	fmt.Println(err)
}
