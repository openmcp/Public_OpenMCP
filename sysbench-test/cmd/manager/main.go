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
	bytes2 "bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"openmcp/openmcp/util"
	"time"
	"unsafe"
)
type HttpManager struct{
	HTTPServer_IP string
	HTTPServer_PORT string

	CpuBenchStatus string
	MemoryBenchStatus string
	DiskBenchStatus string
	NetworkBenchStatus string

	CpuLevel string
	MemoryLevel string
	DiskLevel string
	NetworkLevel string

	CpuStopChan chan bool
	CpuStopOkChan chan bool

	MemoryStopChan chan bool
	MemoryStopOkChan chan bool

	DiskStopChan chan bool
	DiskStopOkChan chan bool



}


var START_CPU_CMD_START string = "sysbench-cpu --test=cpu --cpu-max-prime=" // 10000
var START_CPU_CMD_END string =" run"
var START_MEMORY_CMD_START string = "sysbench-memory --test=memory --memory-block-size="
var START_MEMORY_CMD_END string = " --memory-total-size=10G run"
var START_DISK_CMD_START string = "sysbench-disk --test=fileio --file-test-mode=rndrw --file-total-size=" // 100G
var START_DISK_CMD_END1 string = " prepare"
var START_DISK_CMD_END2 string = " run"
var START_DISK_CMD_END3 string = " cleanup"
//var START_NETWORK_CMD string = ""

//var START_CPU_CMD string = "sysbench-cpu --test=cpu --cpu-max-prime=500 run" // 10000
//var START_MEMORY_CMD string = "sysbench-memory --test=memory --memory-block-size=50M --memory-total-size=50M run"
//var START_DISK_CMD1 string = "sysbench --test=fileio --file-total-size=50M --file-test-mode=rndrw prepare" // 100G
//var START_DISK_CMD2 string = "sysbench --test=fileio --file-total-size=50M --file-test-mode=rndrw run" // 100G
//var START_DISK_CMD3 string = "sysbench --test=fileio --file-total-size=50M --file-test-mode=rndrw cleanup" // 100G

//var STOP_CPU_CMD string =  "kill -9 `ps -ef | grep 'sysbench-cpu' | awk 'NR==1{print $2}'`"
//var STOP_MEMORY_CMD string =  "kill -9 `ps -ef | grep 'sysbench-memory' | awk 'NR==1{print $2}'`"
//var STOP_DISK_CMD string =  "kill -9 `ps -ef | grep 'sysbench-disk' | awk 'NR==1{print $2}'`"
//var STOP_NETWORK_CMD string =  "kill -9 `ps -ef | grep 'sysbench-network' | awk 'NR==1{print $2}'`"


func getMsec(level string) int{
	msec := 1000
	//if level == "1" {
	//	msec = 5000
	//} else if level == "2" {
	//	msec = 3000
	//} else if level == "3" {
	//	msec = 1000
	//} else if level == "4" {
	//	msec = 500
	//} else if level == "5" {
	//	msec = 100
	//}
	return msec
}
func getValue(rType, level string) string{
	value := ""
	if rType == "cpu" {
		if level == "1"{
			value = "100"
		} else if level == "2" {
			value = "300"
		}else if level == "3" {
			value = "500"
		}else if level == "4" {
			value = "1000"
		}else if level == "5" {
			value = "1500"
		}
	} else if rType == "memory" {
		if level == "1"{
			value = "100M"
		} else if level == "2" {
			value = "200M"
		}else if level == "3" {
			value = "300M"
		}else if level == "4" {
			value = "500M"
		}else if level == "5" {
			value = "1G"
		}
	} else if rType == "disk" {
		if level == "1"{
			value = "20M"
		} else if level == "2" {
			value = "50M"
		}else if level == "3" {
			value = "100M"
		}else if level == "4" {
			value = "200M"
		}else if level == "5" {
			value = "500M"
		}
	}
	return value
}
func (h *HttpManager) help(w http.ResponseWriter, r *http.Request){
	fmt.Println("Connect Service")

	w.Write([]byte("OpenMCP Service Response\n"))

}
func (h *HttpManager)sysbench_daemon_cpu() {
	msec := getMsec(h.CpuLevel)
	value := getValue("cpu", h.CpuLevel)
	cmdStr := START_CPU_CMD_START + value + START_CPU_CMD_END
	fmt.Println(cmdStr)
	_, err := util.CmdExec(cmdStr)
	if err != nil {
		fmt.Println("Err !", err)
	}
	for {
		select {
		case <-time.After(time.Duration(msec) * time.Millisecond):

			fmt.Println(cmdStr)
			_, err := util.CmdExec(cmdStr)
			if err != nil {
				fmt.Println("Err !", err)
			}
		case <-h.CpuStopChan:
			h.CpuBenchStatus = "STOP"
			fmt.Println("Send Current Status: STOP")
			h.CpuStopOkChan <- true
			return
		}
	}
}

func (h *HttpManager)sysbench_daemon_memory() {
	msec := getMsec(h.MemoryLevel)
	value := getValue("memory", h.MemoryLevel)

	cmdStr := START_MEMORY_CMD_START + value + START_MEMORY_CMD_END
	fmt.Println(cmdStr)
	_, err := util.CmdExec(cmdStr)
	if err != nil {
		fmt.Println("Err !", err)
	}
	for {
		select {
		case <-time.After(time.Duration(msec) * time.Millisecond):

			fmt.Println(cmdStr)
			_, err := util.CmdExec(cmdStr)
			if err != nil {
				fmt.Println("Err !", err)
			}
		case <-h.MemoryStopChan:
			h.MemoryBenchStatus = "STOP"
			fmt.Println("Send Current Status: STOP")
			h.MemoryStopOkChan <- true
			return
		}
	}
}
func (h *HttpManager)sysbench_daemon_disk() {
	msec := getMsec(h.DiskLevel)
	value := getValue("disk", h.DiskLevel)

	cmdStr := START_DISK_CMD_START + value

	cmdStr2 := cmdStr + START_DISK_CMD_END1
	fmt.Println(cmdStr2)
	_, err := util.CmdExec(cmdStr2)
	if err != nil {
		fmt.Println("Err !", err)
	}
	cmdStr3 := cmdStr + START_DISK_CMD_END2
	fmt.Println(cmdStr3)
	_, err = util.CmdExec(cmdStr3)
	if err != nil {
		fmt.Println("Err !", err)
	}
	cmdStr4 := cmdStr + START_DISK_CMD_END3
	fmt.Println(cmdStr4)
	_, err = util.CmdExec(cmdStr4)
	if err != nil {
		fmt.Println("Err !", err)
	}

	for {
		select {
		case <-time.After(time.Duration(msec) * time.Millisecond):
			fmt.Println(cmdStr2)
			_, err := util.CmdExec(cmdStr2)
			if err != nil {
				fmt.Println("Err !", err)
			}
			fmt.Println(cmdStr3)
			_, err = util.CmdExec(cmdStr3)
			if err != nil {
				fmt.Println("Err !", err)
			}
			fmt.Println(cmdStr4)
			_, err = util.CmdExec(cmdStr4)
			if err != nil {
				fmt.Println("Err !", err)
			}
		case <-h.DiskStopChan:
			h.DiskBenchStatus = "STOP"
			fmt.Println("Send Current Status: STOP")
			h.DiskStopOkChan <- true
			return
		}
	}
}


func (h *HttpManager) networkrx_load(){
	fmt.Println("RX START!")
	serverIP := "http://10.0.3.12:9090"
	urlPath := "/networkrx"


	for {
		if h.NetworkBenchStatus == "STOP" {
			break
		}

		//create request object
		req, err := http.NewRequest("GET", serverIP+urlPath, nil)

		if err != nil {
			fmt.Println(err)
			h.NetworkBenchStatus = "STOP"
			break
			//panic(err)
		}
		//create client object
		client := &http.Client{}

		//send request
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
			h.NetworkBenchStatus = "STOP"
			break
			//panic(err)
		}
		defer resp.Body.Close()

		//receive result
		bytes, _ := ioutil.ReadAll(resp.Body)
		str := string(bytes)
		fmt.Println(str)
	}

}
func (h *HttpManager) networktx_load(){
	fmt.Println("TX START!")
	serverIP := "http://10.0.3.12:9090"
	urlPath := "/networktx"

	a := ""
	for i := 0 ; i < 10000 ; i++ {
		a = a+"aaaaaaaaaaaaaaaaaaaaa"
	}
	fmt.Println(unsafe.Sizeof(a))

	for {
		if h.NetworkBenchStatus == "STOP" {
			break
		}

		//create request object
		req, err := http.NewRequest("POST", serverIP+urlPath, bytes2.NewBufferString(a))

		if err != nil {
			fmt.Println(err)
			h.NetworkBenchStatus = "STOP"
			break
			//panic(err)
		}

		req.Header.Set("Content-Type", "text/plain")
		//create client object
		client := &http.Client{}

		//send request
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
			h.NetworkBenchStatus = "STOP"
			break
			//panic(err)
		}
		defer resp.Body.Close()

		//receive result
		bytes, _ := ioutil.ReadAll(resp.Body)
		str := string(bytes)
		fmt.Println(str)
	}

}

func (h * HttpManager) printStatus(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("CPU Status : "+ h.CpuBenchStatus))
	if h.CpuBenchStatus == "START" {
		w.Write([]byte(" / Level : "+ h.CpuLevel))
	}
	w.Write([]byte("\n"))

	w.Write([]byte("Memory Status : "+ h.MemoryBenchStatus))
	if h.MemoryBenchStatus == "START" {
		w.Write([]byte(" / Level : "+ h.MemoryLevel))
	}
	w.Write([]byte("\n"))

	w.Write([]byte("Disk Status : "+ h.DiskBenchStatus))
	if h.DiskBenchStatus == "START" {
		w.Write([]byte(" / Level : "+ h.DiskLevel))
	}
	w.Write([]byte("\n"))

	w.Write([]byte("Network Status : "+ h.NetworkBenchStatus))
	if h.NetworkBenchStatus == "START" {
		w.Write([]byte(" / Level : "+ h.NetworkLevel))
	}
	w.Write([]byte("\n"))

}

//start
func (h *HttpManager) cpu_start(w http.ResponseWriter, r *http.Request){
	levels, ok := r.URL.Query()["v"]

	if !ok || len(levels[0]) < 1 {
		w.Write([]byte("Url Param 'v' is missing"))
		return
	}
	level := levels[0]

	if h.CpuBenchStatus == "STOP" {
		fmt.Println("SYSBENCH DAEMON CPU START")
		h.CpuBenchStatus = "START"
		h.CpuLevel = level
		go h.sysbench_daemon_cpu()
	} else {
		w.Write([]byte("Failed CPU Status : "+h.CpuBenchStatus +", First STOP CPU\n"))
	}
	h.printStatus(w, r)
}
func (h *HttpManager) memory_start(w http.ResponseWriter, r *http.Request){
	levels, ok := r.URL.Query()["v"]

	if !ok || len(levels[0]) < 1 {
		w.Write([]byte("Url Param 'v' is missing"))
		return
	}
	level := levels[0]

	if h.MemoryBenchStatus == "STOP" {
		h.MemoryBenchStatus = "START"
		h.MemoryLevel = level
		go h.sysbench_daemon_memory()
	}else {
		w.Write([]byte("Failed Memory Status : "+h.MemoryBenchStatus +", First STOP Memory\n"))
	}
	h.printStatus(w, r)
}
func (h *HttpManager) disk_start(w http.ResponseWriter, r *http.Request){
	levels, ok := r.URL.Query()["v"]

	if !ok || len(levels[0]) < 1 {
		w.Write([]byte("Url Param 'v' is missing"))
		return
	}
	level := levels[0]

	if h.DiskBenchStatus == "STOP" {
		h.DiskBenchStatus = "START"
		h.DiskLevel = level
		go h.sysbench_daemon_disk()
	}else {
		w.Write([]byte("Failed Disk Status : "+h.DiskBenchStatus +", First STOP Disk\n"))
	}
	h.printStatus(w, r)
}
func (h *HttpManager) networkrx_start(w http.ResponseWriter, r *http.Request){
//receive large data
	levels, ok := r.URL.Query()["v"]

	if !ok || len(levels[0]) < 1 {
		w.Write([]byte("Url Param 'v' is missing"))
		return
	}
	level := levels[0]

	if h.NetworkBenchStatus == "STOP" {
		h.NetworkBenchStatus = "START"
		h.NetworkLevel = level
		go h.networkrx_load()
	} else {
		w.Write([]byte("Failed Network Status : "+h.NetworkBenchStatus +", First STOP Network\n"))
	}
	h.printStatus(w, r)
}
func (h *HttpManager) networktx_start(w http.ResponseWriter, r *http.Request){
//transmit large data
	levels, ok := r.URL.Query()["v"]

	if !ok || len(levels[0]) < 1 {
		w.Write([]byte("Url Param 'v' is missing"))
		return
	} else {
		w.Write([]byte("Failed Network Status : "+h.NetworkBenchStatus +", First STOP Network\n"))
	}
	level := levels[0]

	if h.NetworkBenchStatus == "STOP" {
		h.NetworkBenchStatus = "START"
		h.NetworkLevel = level
		go h.networktx_load()
	} else {

	}
	h.printStatus(w, r)
}

//stop
func (h *HttpManager) cpu_stop(w http.ResponseWriter, r *http.Request){
	fmt.Println("CPU STOP Called")
	if h.CpuBenchStatus == "START" {
		fmt.Println("CPU Current Status: START")
		h.CpuStopChan <- true
		fmt.Println("Send CPU STOP Channel")
		<- h.CpuStopOkChan
		fmt.Println("Recv CPU STOP Channel")

	}
	h.printStatus(w, r)
}
func (h *HttpManager) memory_stop(w http.ResponseWriter, r *http.Request){
	fmt.Println("Memory STOP Called")
	if h.MemoryBenchStatus == "START" {
		fmt.Println("Memory Current Status: START")
		h.MemoryStopChan <- true
		fmt.Println("Send Memory STOP Channel")
		<- h.MemoryStopOkChan
		fmt.Println("Recv Memory STOP Channel")

	}
	h.printStatus(w, r)
}
func (h *HttpManager) disk_stop(w http.ResponseWriter, r *http.Request){
	fmt.Println("Disk STOP Called")
	if h.DiskBenchStatus == "START" {
		fmt.Println("Disk Current Status: START")
		h.DiskStopChan <- true
		fmt.Println("Send Disk STOP Channel")
		<- h.DiskStopOkChan
		fmt.Println("Recv Disk STOP Channel")

	}
	h.printStatus(w, r)
}
func (h *HttpManager) networkrx_stop(w http.ResponseWriter, r *http.Request){
	if h.NetworkBenchStatus == "START" {
		h.NetworkBenchStatus = "STOP"
	}
	h.printStatus(w, r)
}
func (h *HttpManager) networktx_stop(w http.ResponseWriter, r *http.Request){
	if h.NetworkBenchStatus == "START" {
		h.NetworkBenchStatus = "STOP"

	}
	h.printStatus(w, r)
}

func main() {
	//HTTPServer_IP :=   os.Getenv("SERVER_IP")
	//HTTPServer_PORT := os.Getenv("SERVER_PORT")

	HTTPServer_PORT := "7070"

	httpManager := &HttpManager{
		HTTPServer_IP:  "",
		HTTPServer_PORT: HTTPServer_PORT,

		CpuBenchStatus:   "STOP",
		MemoryBenchStatus:   "STOP",
		NetworkBenchStatus:   "STOP",
		DiskBenchStatus:   "STOP",

		CpuStopChan: make(chan bool),
		CpuStopOkChan: make(chan bool),

		MemoryStopChan: make(chan bool),
		MemoryStopOkChan: make(chan bool),

		DiskStopChan: make(chan bool),
		DiskStopOkChan: make(chan bool),

	}

	handler := http.NewServeMux()

	handler.HandleFunc("/", httpManager.help)
	handler.HandleFunc("/cpu/start", httpManager.cpu_start)
	handler.HandleFunc("/memory/start", httpManager.memory_start)
	handler.HandleFunc("/disk/start", httpManager.disk_start)
	handler.HandleFunc("/networkrx/start", httpManager.networkrx_start)
	handler.HandleFunc("/networktx/start", httpManager.networktx_start)
	handler.HandleFunc("/cpu/stop", httpManager.cpu_stop)
	handler.HandleFunc("/memory/stop", httpManager.memory_stop)
	handler.HandleFunc("/disk/stop", httpManager.disk_stop)
	handler.HandleFunc("/networkrx/stop", httpManager.networkrx_stop)
	handler.HandleFunc("/networktx/stop", httpManager.networktx_stop)



	fmt.Println("Sysbench-test Server Start ! " + ":" + HTTPServer_PORT)
	server := &http.Server{Addr: ":" + HTTPServer_PORT, Handler: handler}
	err := server.ListenAndServe()

	fmt.Println(err)
}
