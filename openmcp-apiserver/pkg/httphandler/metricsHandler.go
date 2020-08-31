package httphandler

import (
	"fmt"
	"net/http"
)


func (h *HttpManager)MetricsHandler(w http.ResponseWriter, r *http.Request, splitUrl []string) {

	ns := ""
	fmt.Println(splitUrl)

	if splitUrl[0] == "namespaces" {
		splitUrl = PopLeftSlice(splitUrl)
		fmt.Println(splitUrl)
		ns = splitUrl[0]
		splitUrl = PopLeftSlice(splitUrl)
		fmt.Println(splitUrl)
	}

	Node_Or_Pod := splitUrl[0]
	splitUrl = PopLeftSlice(splitUrl)
	fmt.Println(splitUrl)

	metric := splitUrl[0]
	splitUrl = PopLeftSlice(splitUrl)
	fmt.Println(splitUrl)

	fmt.Println(ns, Node_Or_Pod, metric)




	w.Write([]byte("Hello"))
}
