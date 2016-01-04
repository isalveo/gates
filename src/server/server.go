package server

import (
	"net/http"
	"proxy"
)

func init() {
	http.HandleFunc("/health", NewHandler("health", new(proxy.HealthCheckHandler)))
	http.HandleFunc("/", NewHandler("proxy", new(proxy.ReverseProxyHandler)))
}
