package proxy

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
	"strings"
	"testing"
)

func TestReverseProxy404(t *testing.T) {

	ctx := context.Background()
	h := new(ReverseProxyHandler)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", fmt.Sprintf("/v1/test"), nil)
	if err != nil {
		t.Fatal(err)
	}
	status, err := h.ServeHTTP(ctx, w, r)

	if status != 404 {
		t.Fatalf("Response code should eq ", status)
	}
}

func TestReverseProxy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"proxy":"ok"}`)
	}))
	defer ts.Close()

	m := getField(initializers.Registry(), "hermes")
	m["v1"] = strings.Replace(ts.URL, "http://", "", 1)

	ctx := context.Background()
	h := new(ReverseProxyHandler)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", fmt.Sprintf("/hermes/v1"), nil)
	if err != nil {
		t.Fatal(err)
	}
	status, err := h.ServeHTTP(ctx, w, r)

	if status != 200 {
		t.Fatalf("Response code should eq 200 but was", status)
	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, w.Body)
	if err != nil {
		t.Fatal(err)
	}

	if buf.String() != `{"proxy":"ok"}` {
		t.Fatalf("Unexpected response", buf.String())
	}
}

func TestExtractNameVersion(t *testing.T) {
	r, err := http.NewRequest("POST", fmt.Sprintf("/hermes/v1"), nil)
	if err != nil {
		t.Fatal(err)
	}

	name, version, endpoint, e := extractNameVersion(r.URL)
	if e != nil {
		t.Fatalf("Error occured", e)
	}

	if name != "hermes" {
		t.Fatalf("Service name is incorrect:", name)
	}

	if version != "v1" {
		t.Fatalf("Service version is incorrect:", version)
	}

	if endpoint == "" {
		t.Fatalf("Endpoint is incorrect")
	}

	r, err = http.NewRequest("POST", fmt.Sprintf("/yoda/v1"), nil)
	if err != nil {
		t.Fatal(err)
	}

	name, version, endpoint, e = extractNameVersion(r.URL)
	if e != nil {
		t.Fatalf("Error occured", e)
	}

	if name != "yoda" {
		t.Fatalf("Service name is incorrect:", name)
	}

	if version != "v1" {
		t.Fatalf("Service version is incorrect:", version)
	}

	if endpoint == "" {
		t.Fatalf("Endpoint is incorrect.")
	}

	r, err = http.NewRequest("POST", fmt.Sprintf("/not-found/v1"), nil)
	if err != nil {
		t.Fatal(err)
	}

	name, version, endpoint, e = extractNameVersion(r.URL)
	if !strings.Contains(e.Error(), "Not Found") {
		t.Fatalf("Should be Not Found but was", e.Error())
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
