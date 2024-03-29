package restacular

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestView struct{}

var basicHandler = func(ctx *Context) {
	ctx.Response.Send(200, []byte("Hello world"))
}

var handlerWithParam = func(ctx *Context) {
	id := ctx.Params["id"]
	ctx.Response.Send(200, string(id))
}

func TestAddingRoute(t *testing.T) {
	server := NewServer()
	pattern := "/lets-kill-some-ducks"

	server.Get(pattern, basicHandler)

	foundPattern := server.routes[0].pattern

	if foundPattern != pattern {
		t.Error("Route was not present in the router after adding it")
	}
}

func TestAddingRouteWithTrailingSlash(t *testing.T) {
	server := NewServer()
	pattern := "/lets-kill-some-ducks/"
	expectedPattern := "/lets-kill-some-ducks"

	server.Get(pattern, basicHandler)
	foundPattern := server.routes[0].pattern

	if foundPattern != expectedPattern {
		t.Error("Route was not present in the router after adding it")
	}
	if foundPattern == pattern {
		t.Error("Trailing slash wasn't removed when adding the route")
	}
}

func TestMatchExistingRoute(t *testing.T) {
	server := NewServer()
	pattern := "/ducks/:id"

	server.Get(pattern, basicHandler)

	routeFound, _ := server.matchRequest("GET", "/ducks/0irfer8")

	if routeFound == nil {
		t.Error("Couldn't find a match for the given url")
	}
}

func TestMatchUnexistingRoute(t *testing.T) {
	server := NewServer()
	pattern := "/ducks/:id"

	server.Get(pattern, basicHandler)
	routeFound, _ := server.matchRequest("GET", "/ducks/0irfer8/ducklings")

	if routeFound != nil {
		t.Error("Found a match for the given url when we shouldn't have")
	}
}

// High-level tests -----

func TestServeOk(t *testing.T) {
	server := NewServer()

	server.Get("/ducks/:id", handlerWithParam)

	req, err := http.NewRequest("GET", "/ducks/0irfer8", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Received http code %d instead of 200", w.Code)
	}

	unzippedBody, _ := gzip.NewReader(w.Body)
	buf := new(bytes.Buffer)
	buf.ReadFrom(unzippedBody)
	var body string
	test := buf.String()

	_ = json.Unmarshal([]byte(test), &body)

	if body != "0irfer8" {
		t.Errorf("Got %s as body instead of 0irfer8", body)
	}
}

func TestServeNotExistingRoute(t *testing.T) {
	server := NewServer()

	req, err := http.NewRequest("GET", "/ducks/0irfer8", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Received http code %d instead of 404", w.Code)
	}
}

func TestServeTrailingSlash(t *testing.T) {
	server := NewServer()

	server.Get("/ducks/:id", handlerWithParam)

	req, err := http.NewRequest("GET", "/ducks/0irfer8/", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Received http code %d instead of 404", w.Code)
	}
}
