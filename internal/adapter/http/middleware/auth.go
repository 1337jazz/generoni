package middleware

import "net/http"

// Add your auth logic here. In the set up it runs after the RateLimit, so failed attempts are still limited,
// but it's up to you to decide if that's a good idea
func Auth(next http.Handler) http.Handler {
	return next
}
