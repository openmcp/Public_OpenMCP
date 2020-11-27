
package main
import (
	"fmt"
	"time"
)
type HttpManager struct {
	CpuStopChan chan bool
}
func (h *HttpManager)func1() {
	<-h.CpuStopChan
	fmt.Println("Hi!")
}
func main() {
	httpManager := &HttpManager{
	//	CpuStopChan: make(chan bool),
	}
	go httpManager.func1()
	time.Sleep(3 * time.Second)
	httpManager.CpuStopChan <- true
	time.Sleep(1 * time.Second)
	fmt.Println("main Exit")
}
