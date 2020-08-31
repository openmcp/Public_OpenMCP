package httphandler

import (
	"net/http"
)

func (h *HttpManager)MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
