package restacular

import (
	"log"
	"net/http"
	"strings"
)

type Context struct {
	Request  *Request
	Response Response
	Params   map[string]string
}

type route struct {
	httpMethod string
	pattern    string
	handler    func(*Context)
}

type Server struct {
	routes []route
	debug  bool
}

func NewServer() *Server {
	return &Server{}
}

func newContext(request *Request, response Response, params map[string]string) *Context {
	return &Context{request, response, params}
}

func (server *Server) addRoute(method string, pattern string, handler func(*Context)) {
	if len(pattern) > 1 && strings.HasSuffix(pattern, "/") {
		pattern = strings.TrimRight(pattern, "/")
	}
	server.routes = append(server.routes, route{method, pattern, handler})
}

func (server *Server) Get(pattern string, handler func(*Context)) {
	server.addRoute("GET", pattern, handler)
}

func (server *Server) Post(pattern string, handler func(*Context)) {
	server.addRoute("POST", pattern, handler)
}

func (server *Server) Put(pattern string, handler func(*Context)) {
	server.addRoute("PUT", pattern, handler)
}

func (server *Server) Delete(pattern string, handler func(*Context)) {
	server.addRoute("DELETE", pattern, handler)
}

func (server *Server) matchRequest(method, url string) (*route, map[string]string) {
	params := make(map[string]string)

	for _, route := range server.routes {
		if method != route.httpMethod {
			continue
		}

		currentSplitted := strings.Split(route.pattern[1:], "/")
		requestSplitted := strings.Split(url[1:], "/")

		if len(currentSplitted) != len(requestSplitted) {
			continue
		}

		matches := false

		for i := range currentSplitted {
			if strings.HasPrefix(currentSplitted[i], ":") {
				params[currentSplitted[i][1:]] = requestSplitted[i]
				continue
			}
			if currentSplitted[i] != requestSplitted[i] {
				break
			}
			matches = true
		}

		if matches {
			return &route, params
		}
	}
	return nil, params
}

func (server *Server) dispatch() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		requestPath := req.URL.Path
		requestMethod := req.Method

		log.Println(requestMethod, requestPath)

		routeFound, params := server.matchRequest(requestMethod, requestPath)

		if routeFound != nil {
			context := newContext(&Request{*req}, Response{resp}, params)
			routeFound.handler(context)
			return
		}

		http.NotFound(Response{resp}, req)
		return

	}

}

func (server *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	handler := gzipWrapper(server.dispatch())
	handler(resp, req)
}
