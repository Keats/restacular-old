package restacular

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Route struct {
	HttpMethod string
	Pattern    string
	Handler    func(Response, *http.Request)
}

type Router struct {
	routes []Route
}

type matchingRoute struct {
	route  Route
	params map[string]string
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) ShowRoutes() {
	for i := range r.routes {
		fmt.Printf("%v\n", r.routes[i])
	}
}

func (r *Router) AddRoutes(routes ...Route) {
	for i := range routes {
		if len(routes[i].Pattern) > 1 && strings.HasSuffix(routes[i].Pattern, "/") {
			routes[i].Pattern = strings.TrimRight(routes[i].Pattern, "/")
		}
		r.routes = append(r.routes, routes[i])
	}
}

func (r *Router) match(method, url string) *matchingRoute {
	for i := range r.routes {
		if method != r.routes[i].HttpMethod {
			continue
		}

		currentSplitted := strings.Split(r.routes[i].Pattern[1:], "/")
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
			return &matchingRoute{r.routes[i], params}
		}
	}

	return nil
}

func (r *Router) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	requestPath := req.URL.Path
	requestMethod := req.Method

	routeFound := r.match(requestMethod, requestPath)
	if routeFound != nil {
		values := req.URL.Query()
		for i := range routeFound.params {
			values.Add(i, routeFound.params[i])
		}
		req.URL.RawQuery = url.Values(values).Encode() + "&" + req.URL.RawQuery
		routeFound.route.Handler(Response{resp}, req)
	} else {
		http.NotFound(Response{resp}, req)
	}
}
