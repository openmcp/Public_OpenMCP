package httphandler

import (
	"net/http"
	"strings"
)

func (h *HttpManager) RouteHandler(w http.ResponseWriter, r *http.Request) {

	splitUrl := strings.Split(r.URL.Path, "/")

	if splitUrl[1] == "apis" || splitUrl[1] == "api" {
		h.ApiHandler(w, r)
	} else if splitUrl[1] == "metrics" {
		h.MetricsHandler(w, r)
	}

}