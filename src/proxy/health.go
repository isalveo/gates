package proxy

import (
	"fmt"
	"golang.org/x/net/context"
	"net/http"
)

type HealthCheckHandler struct{}

func (h *HealthCheckHandler) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	return h.HealthCheck(w, r)
}

func (h *HealthCheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) (int, error) {
	_, err := fmt.Fprint(w, "ok")
	if err != nil {
		return 500, err
	}

	return 200, nil
}
