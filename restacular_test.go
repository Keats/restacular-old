package restacular

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestView struct{}

var basicHandler = func(resp *Response, r *http.Request) {
	resp.Write([]byte("Hello world"))
}

var handlerWithParam = func(resp *Response, r *http.Request) {
	id := r.URL.Query().Get(":id")
	resp.Write([]byte(id))
}

func TestAddingRoute(t *testing.T) {
	app := NewApplication()
	pattern := "/lets-kill-some-ducks"

	app.SetRoutes(
		Route{"GET", pattern, basicHandler},
	)
	foundPattern := app.routes[0].Pattern

	if foundPattern != pattern {
		t.Error("Route was not present in the router after adding it")
	}
}

func TestAddingRouteWithTrailingSlash(t *testing.T) {
	app := NewApplication()
	pattern := "/lets-kill-some-ducks/"
	expectedPattern := "/lets-kill-some-ducks"

	app.SetRoutes(
		Route{"GET", pattern, basicHandler},
	)
	foundPattern := app.routes[0].Pattern

	if foundPattern != expectedPattern {
		t.Error("Route was not present in the router after adding it")
	}
	if foundPattern == pattern {
		t.Error("Trailing slash wasn't removed when adding the route")
	}
}

func TestMatchExistingRoute(t *testing.T) {
	app := NewApplication()
	pattern := "/ducks/:id"

	app.SetRoutes(
		Route{"GET", pattern, basicHandler},
	)

	routeFound := app.matchRequest("GET", "/ducks/0irfer8")

	if routeFound == nil {
		t.Error("Couldn't find a match for the given url")
	}
}

func TestMatchUnexistingRoute(t *testing.T) {
	app := NewApplication()
	pattern := "/ducks/:id"

	app.SetRoutes(
		Route{"GET", pattern, basicHandler},
	)
	routeFound := app.matchRequest("GET", "/ducks/0irfer8/ducklings")

	if routeFound != nil {
		t.Error("Found a match for the given url when we shouldn't have")
	}
}

// High-level tests -----

func TestServeOk(t *testing.T) {
	app := NewApplication()

	app.SetRoutes(
		Route{"GET", "/ducks/:id", handlerWithParam},
	)

	req, err := http.NewRequest("GET", "/ducks/0irfer8", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Received http code %d instead of 200", w.Code)
	}

	unzippedBody, _ := gzip.NewReader(w.Body)
	buf := new(bytes.Buffer)
	buf.ReadFrom(unzippedBody)
	body := buf.String()

	if body != "0irfer8" {
		t.Errorf("Got %s as body instead of 0irfer8", body)
	}
}

func TestServeNotExistingRoute(t *testing.T) {
	app := NewApplication()

	req, err := http.NewRequest("GET", "/ducks/0irfer8", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Received http code %d instead of 404", w.Code)
	}
}

func TestServeTrailingSlash(t *testing.T) {
	app := NewApplication()

	app.SetRoutes(
		Route{"GET", "/ducks/:id", handlerWithParam},
	)

	req, err := http.NewRequest("GET", "/ducks/0irfer8/", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Received http code %d instead of 404", w.Code)
	}
}
