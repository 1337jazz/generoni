package http

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/1337jazz/generoni/web"
)

// spa serves the frontend from web/dist. Unknown paths get index.html so client-side routes work
func spa() http.Handler {
	files, _ := fs.Sub(web.Dist, "dist")
	server := http.FileServerFS(files)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fs.Stat(files, strings.TrimPrefix(r.URL.Path, "/")); err != nil {
			r.URL.Path = "/"
		}
		server.ServeHTTP(w, r)
	})
}
