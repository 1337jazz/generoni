package handlers

import "net/http"

// Demo sits behind the rate limited chain. Throwaway, delete it along with
// Panic once there are real routes.
func Demo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("rate limited\n"))
}

// Panic shows that Recover turns a panic into the 500 page instead of a
// dropped connection.
func Panic(w http.ResponseWriter, r *http.Request) {
	panic("demo panic")
}
