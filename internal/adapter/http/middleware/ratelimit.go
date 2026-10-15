package middleware

import (
	"net/http"
	"time"

	"github.com/1337jazz/generoni/internal/adapter/http/clientip"
	"github.com/go-chi/httprate"
)

// RateLimit allows each client IP up to limit requests per window, then
// responds 429 Too Many Requests.
func RateLimit(limit int, window time.Duration) func(http.Handler) http.Handler {
	return httprate.LimitBy(limit, window, clientip.FromRequest)
}
