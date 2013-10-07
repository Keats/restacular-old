package restacular

import (
	"compress/gzip"
	"net/http"
)

type gzipResponseWriter struct {
	http.ResponseWriter
}

func (self *gzipResponseWriter) Write(b []byte) (int, error) {
	gzipWriter := gzip.NewWriter(self.ResponseWriter)
	defer gzipWriter.Close()
	return gzipWriter.Write(b)
}

func gzipWrapper(handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Encoding", "gzip")
		handler(&gzipResponseWriter{resp}, req)
	}
}
