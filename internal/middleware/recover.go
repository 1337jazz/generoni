package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/1337jazz/generoni/internal/handlers"
)

// Recover turns a panic into the 500 page. Without it net/http drops the
// connection and the client sees nothing.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			// net/http uses this panic to mean "drop the connection on
			// purpose", so pass it on instead of answering.
			if rec == http.ErrAbortHandler {
				panic(rec)
			}

			slog.LogAttrs(r.Context(), slog.LevelError, "panic",
				slog.Any("value", rec),
				slog.String("request_id", RequestIDFrom(r.Context())),
				slog.String("stack", string(debug.Stack())),
			)

			// Show the 500 page, but never leak the panic to the client.
			handlers.ServerError(w, r)
		}()
		next.ServeHTTP(w, r)
	})
}
