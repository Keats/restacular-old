package restacular

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// User added route
type Route struct {
	HttpMethod string
	Pattern    string
	Func       func(*Response, *http.Request)
}

// The main object, only cares about route
type Application struct {
	routes []Route
}

// Internal type to get the result from a match
// Params will be added to the request by the dispatch method
type matchingRoute struct {
	route  Route
	params map[string]string
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) SetRoutes(routes ...Route) {
	for i := range routes {
		// if pattern is not the root slash, remove trailing slash
		if len(routes[i].Pattern) > 1 && strings.HasSuffix(routes[i].Pattern, "/") {
			routes[i].Pattern = strings.TrimRight(routes[i].Pattern, "/")
		}
		app.routes = append(app.routes, routes[i])
	}
}

// Here for debug purpose
func (app *Application) ShowRoutes() {
	for i := range app.routes {
		fmt.Printf("%v\n", app.routes[i])
	}
}

// Finds if there is a match and get the params at the same time
func (app *Application) matchRequest(method, url string) *matchingRoute {
	for i := range app.routes {
		// Easy check on the method first
		if method != app.routes[i].HttpMethod {
			continue
		}

		// Getting each piece of the URL, no need to split on root / though
		currentSplitted := strings.Split(app.routes[i].Pattern[1:], "/")
		requestSplitted := strings.Split(url[1:], "/")

		if len(currentSplitted) != len(requestSplitted) {
			continue
		}

		matches := false
		params := make(map[string]string)

		for i := range currentSplitted {
			if strings.HasPrefix(currentSplitted[i], ":") {
				params[currentSplitted[i]] = requestSplitted[i]
				continue
			}
			if currentSplitted[i] != requestSplitted[i] {
				break
			}
			matches = true
		}

		if matches {
			return &matchingRoute{app.routes[i], params}
		}
	}

	return nil
}

// Returns an anonymous function that actually calls the user-defined
// view after setting the params in the URL.RawQuery (or a 404)
func (app *Application) dispatch() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		requestPath := req.URL.Path
		requestMethod := req.Method
		log.Println(requestMethod, requestPath)

		routeFound := app.matchRequest(requestMethod, requestPath)

		if routeFound != nil {
			values := req.URL.Query()

			for i := range routeFound.params {
				values.Add(i, routeFound.params[i])
			}

			req.URL.RawQuery = url.Values(values).Encode() + "&" + req.URL.RawQuery
			routeFound.route.Func(&Response{resp}, req)
		} else {
			http.NotFound(Response{resp}, req)
		}
	}

}

// Basic implementation to satisfy interface, easy to add wrapper around like Gzip
func (app *Application) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	handler := gzipWrapper(app.dispatch())
	handler(resp, req)
}
