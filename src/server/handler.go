package server

import (
	"golang.org/x/net/context"
	"initializers"
	"net/http"
)

type Handler interface {
	ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error)
}

func NewHandler(name string, h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		instrumentation := NewInstrumentation(name)
		ctx = context.WithValue(ctx, "instrumentation", instrumentation)

		instrumentation.StartHandler()
		status, err := h.ServeHTTP(ctx, w, r)
		instrumentation.StopHandler(status)

		if err != nil {
			if status == http.StatusNotFound {
				http.NotFound(w, r)
			} else {
				initializers.Error(w, err, status)
			}
		}

		initializers.LogRequest(r, status, err)
	}
}
