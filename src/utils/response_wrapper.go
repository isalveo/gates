package utils

import (
	"net/http"
)

// wrapping for logging purposes
type LoggedResponse struct {
	W      http.ResponseWriter
	R      *http.Request
	Status int
}

func (w *LoggedResponse) Flush() {
	if wf, ok := w.W.(http.Flusher); ok {
		wf.Flush()
	}
}

func (w *LoggedResponse) Header() http.Header         { return w.W.Header() }
func (w *LoggedResponse) Write(d []byte) (int, error) { return w.W.Write(d) }

func (w *LoggedResponse) WriteHeader(status int) {
	w.Status = status
	w.W.WriteHeader(status)
}
