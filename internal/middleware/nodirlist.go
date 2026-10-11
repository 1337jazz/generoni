package middleware

import (
	"net/http"
	"strings"
)

// NoDirList returns 404 for directory paths (ending in "/"), so a file server
// never shows a listing of its files.
func NoDirList(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
