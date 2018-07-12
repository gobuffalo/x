package content_type_helper

import (
	"net/http"
	"strings"
)

// This will check for if the requested path has a known file extention and will
// properly remove it and set the content-type if it is not already set. Header
// Content-Type takes precedence.
func AutoSetContentType(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		path := strings.ToLower(req.URL.Path)
		_, hasCT := req.Header["Content-Type"]
		if strings.HasSuffix(path, ".html") {
			req.URL.Path = path[:len(path)-5]
		} else if strings.HasSuffix(path, ".json") {
			req.URL.Path = path[:len(path)-5]
			if !hasCT {
				req.Header["Content-Type"] = []string{"json"}
			}
		} else if strings.HasSuffix(path, ".xml") {
			req.URL.Path = path[:len(path)-4]
			if !hasCT {
				req.Header["Content-Type"] = []string{"xml"}
			}
		}
	})
}
