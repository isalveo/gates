package proxy

import (
	"golang.org/x/net/context"
	"initializers"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
	"utils"
)

type ReverseProxyHandler struct{}

func (h *ReverseProxyHandler) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	return h.ReverseProxy(w, r)
}

func (h *ReverseProxyHandler) ReverseProxy(w http.ResponseWriter, r *http.Request) (int, error) {
	name, version, endpoint, err := extractNameVersion(r.URL)
	if err != nil {
		return 404, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: func(network, addr string) (net.Conn, error) {
			return targetConnection(network, endpoint)
		},
		TLSHandshakeTimeout: 10 * time.Second,
	}

	lw := &utils.LoggedResponse{w, r, 200}

	(&httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			r.URL.Host = name + "/" + version
		},
		Transport: transport,
	}).ServeHTTP(lw, r)
	return lw.Status, nil
}

func extractNameVersion(target *url.URL) (name, version, endpoint string, err error) {
	path := target.Path
	if len(path) > 1 && path[0] == '/' {
		path = path[1:]
	}
	tmp := strings.Split(path, "/")
	if len(tmp) < 2 {
		return "", "", "", initializers.StatusNotFound
	}
	name, version = tmp[0], tmp[1]

	endpoint, err = serviceLookup(name, version)
	if err != nil {
		return "", "", "", err
	}

	target.Path = "/" + strings.Join(tmp[2:], "/")
	return name, version, endpoint, nil
}

func targetConnection(network, endpoint string) (net.Conn, error) {
	for i := 0; i <= 5; i++ {
		conn, err := net.Dial(network, endpoint)
		if err != nil {
			continue
		}
		initializers.Logger.Info("PROXIED:" + endpoint)

		return conn, nil
	}

	return nil, initializers.StatusNotFound
}
