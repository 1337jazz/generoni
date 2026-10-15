package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/1337jazz/generoni/internal/adapter/logger"
)

const RequestIDHeader = "X-Request-Id"

// RequestID gives every request a unique ID, puts it in the request context
// and echoes it in the response headers
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := newRequestID()
		w.Header().Set(RequestIDHeader, id)
		next.ServeHTTP(w, r.WithContext(logger.WithRequestID(r.Context(), id)))
	})
}

// newRequestID returns a random 16-byte hex string
func newRequestID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
