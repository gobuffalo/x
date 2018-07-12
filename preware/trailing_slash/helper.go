package trailing_slash

import (
	"net/http"
	"strings"
)

// This will check for if the requested path has a trailing slash ('/') and remove it.
func RemoveTrailingSlash(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		if strings.HasSuffix(path, "/") && path != "/" {
			req.URL.Path = path[:len(path)-1]
		}
	})
}
