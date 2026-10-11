package middleware

import (
	"net"
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

// RateLimit allows each client IP up to limit requests per window, then
// responds 429 Too Many Requests.
func RateLimit(limit int, window time.Duration) func(http.Handler) http.Handler {
	return httprate.LimitBy(limit, window, clientIP)
}

// clientIP reads X-Forwarded-For and falls back to RemoteAddr if it's not set.
// NOTE: We are assuming you're using a proper proxy and that the proxy is setting
// the X-Forwarded-For header itself, and dropping any values it receives from the client.
// Recommendation: use Caddy or your cloud provider's load balancer.
func clientIP(r *http.Request) (string, error) {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip, nil
	}
	// NOTE: If we get here the proxy didn't get the IP at all. That shouldn't happen
	// if you're using a proper proxy, but we'll have to fall back to the
	// RemoteAddr, which is probably the Docker IP...this needs to be addressed - by you!
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	return ip, err
}
