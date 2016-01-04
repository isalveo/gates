package server

import (
	"bytes"
	"fmt"
	"golang.org/x/net/context"
	"initializers"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

type testHandler struct{}

func (h *testHandler) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	status, err := h.Hello(ctx, w, r)
	if err != nil {
		return status, err
	}

	return status, nil
}

func (h *testHandler) Hello(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	fmt.Fprint(w, "hello")
	return 201, nil
}

func TestNewHandler(t *testing.T) {
	http.HandleFunc("/hello", NewHandler("hello", new(testHandler)))

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	http.DefaultServeMux.ServeHTTP(w, r)

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, w.Body)
	if err != nil {
		t.Fatal(err)
	}

	if buf.String() != "hello" {
		t.Fatalf("Http handling failed")
	}
}

func TestMain(m *testing.M) {
	conf, _ := filepath.Abs("../../config/config.json")
	log, _ := filepath.Abs("../../log/h-gatekeeper.log")
	soa, _ := filepath.Abs("../../config/soa.json")

	configPaths := &initializers.Paths{conf, log, soa}
	initializers.Boot(configPaths, "hermes-gatekeeper")

	os.Exit(m.Run())
}
