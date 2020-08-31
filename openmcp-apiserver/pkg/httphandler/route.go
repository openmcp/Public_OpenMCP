package httphandler

import (
	"fmt"
	"net/http"
	"strings"
)
func PopLeftSlice(splitUrl []string) []string {
	splitUrl = append(splitUrl[:0], splitUrl[1:]...)
	return splitUrl
}

func (h *HttpManager) RouteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("r.URL.Path : ", r.URL.Path)
	splitUrl := strings.Split(r.URL.Path, "/")
	splitUrl = PopLeftSlice(splitUrl)

	if splitUrl[0] == "apis" || splitUrl[0] == "api" {
		h.ApiHandler(w, r)
	} else if splitUrl[0] == "metrics" {

		h.MetricsHandler(w, r, splitUrl)
	}

}