package handlers

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/1337jazz/generoni/web"
)

// Each page gets its own template set (base layout + that page), so pages can
// all define "content" without clashing. Parsed once at startup, so a broken
// template panics here instead of on the first request.
var (
	homePage     = page("home.html")
	notFoundPage = page("404.html")
	errorPage    = page("500.html")
)

func page(name string) *template.Template {
	return template.Must(template.ParseFS(web.Templates, "templates/base.html", "templates/"+name))
}

func Home(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusOK, homePage, nil)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusNotFound, notFoundPage, nil)
}

func ServerError(w http.ResponseWriter, r *http.Request) {
	render(w, r, http.StatusInternalServerError, errorPage, nil)
}

// render executes into a buffer first, so a template error can still become a
// clean 500 instead of a half written page.
func render(w http.ResponseWriter, r *http.Request, status int, t *template.Template, data any) {
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "base", data); err != nil {
		slog.ErrorContext(r.Context(), "render failed", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(status)
	_, _ = buf.WriteTo(w)
}
