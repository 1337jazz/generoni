package clientip

import (
	"net"
	"net/http"
	"strings"
)

// FromRequest returns the client IP: the right-most X-Forwarded-For value, or
// RemoteAddr if the header is not set.
// NOTE: We are assuming exactly one trusted proxy sits in front of the app, and
// that it appends the address it saw to X-Forwarded-For. The right-most value is
// the one that proxy added, so a client can't forge it. Anything to its left
// is client-controlled and is ignored.
// Never expose the app directly, because then the client controls the header.
// Recommendation: use Caddy or your cloud provider's load balancer.
func FromRequest(r *http.Request) (string, error) {
	if values := r.Header.Values("X-Forwarded-For"); len(values) > 0 {
		parts := strings.Split(strings.Join(values, ","), ",")
		if ip := strings.TrimSpace(parts[len(parts)-1]); ip != "" {
			return ip, nil
		}
	}
	// NOTE: If we get here the proxy didn't get the IP at all. That shouldn't happen
	// if you're using a proper proxy, but we'll have to fall back to the
	// RemoteAddr, which is probably the Docker IP...this needs to be addressed - by you!
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	return ip, err
}
