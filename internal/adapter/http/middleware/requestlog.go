package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// RequestLog logs one line per request once it has finished. It must wrap
// Recover, so a recovered panic is logged with the 500 that Recover wrote.
func RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &recorder{ResponseWriter: w, status: http.StatusOK}

		defer func() {
			slog.LogAttrs(r.Context(), slog.LevelInfo, "request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rec.status),
				slog.Int64("bytes", rec.written),
				slog.Float64("duration_ms", float64(time.Since(start).Microseconds())/1000),
			)
		}()

		next.ServeHTTP(rec, r)
	})
}

// recorder remembers the status and size of a response. Unwrap lets
// http.ResponseController reach the real ResponseWriter, so flushing and
// hijacking keep working.
type recorder struct {
	http.ResponseWriter
	status  int
	written int64
}

func (r *recorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *recorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.written += int64(n)
	return n, err
}

func (r *recorder) Unwrap() http.ResponseWriter { return r.ResponseWriter }
