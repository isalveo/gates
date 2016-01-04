package proxy

import (
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthCheck(t *testing.T) {

	ctx := context.Background()
	h := new(HealthCheckHandler)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("get", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	h.ServeHTTP(ctx, w, r)

	if w.Code != 200 {
		t.Fatalf("Code was not 200", w.Code)
	}

	if p, err := ioutil.ReadAll(w.Body); err != nil {
		t.Fatalf("Error occured", err)
	} else if !strings.Contains(string(p), "ok") {
		t.Fatalf("Response should be ok: %s", p)
	}
}
