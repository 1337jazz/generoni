package clientip

import (
	"net"
	"net/http"
)

// FromRequest reads X-Forwarded-For and falls back to RemoteAddr if it's not set.
// NOTE: We are assuming you're using a proper proxy and that the proxy is setting
// the X-Forwarded-For header itself, and dropping any values it receives from the client.
// Recommendation: use Caddy or your cloud provider's load balancer.
func FromRequest(r *http.Request) (string, error) {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip, nil
	}
	// NOTE: If we get here the proxy didn't get the IP at all. That shouldn't happen
	// if you're using a proper proxy, but we'll have to fall back to the
	// RemoteAddr, which is probably the Docker IP...this needs to be addressed - by you!
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	return ip, err
}
