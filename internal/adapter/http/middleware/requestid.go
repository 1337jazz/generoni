package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/1337jazz/generoni/internal/adapter/logger"
)

// RequestIDHeader carries the request ID back to the client, so a user can
// quote it in a bug report and it can be found in the logs.
const RequestIDHeader = "X-Request-Id"

// RequestID gives every request a unique ID, puts it in the request context
// and echoes it in the response headers.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := newRequestID()
		w.Header().Set(RequestIDHeader, id)
		next.ServeHTTP(w, r.WithContext(logger.WithRequestID(r.Context(), id)))
	})
}

func newRequestID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
