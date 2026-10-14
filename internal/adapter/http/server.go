package http

import (
	"encoding/json"
	"net/http"

	"github.com/1337jazz/generoni/internal/adapter/config/environment"
	"github.com/1337jazz/generoni/internal/adapter/container"
	"github.com/1337jazz/generoni/internal/adapter/http/middleware"
	"github.com/justinas/alice"
)

type Server struct {
	*http.Server
	cntr *container.Container
}

func NewServer(cfg environment.Config, cntr *container.Container) *Server {
	s := &Server{cntr: cntr}
	s.Server = &http.Server{
		Addr:              cfg.Addr,
		Handler:           s.routes(cfg),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
	}
	return s
}

// routes builds the handler for the whole application. The global chain
// wraps every route, and the limited chain adds rate limiting to the routes
// that need it
func (s *Server) routes(cfg environment.Config) http.Handler {
	itemHandler := s.cntr.ItemHandler

	limited := alice.New(middleware.RateLimit(cfg.RateLimit, cfg.RateLimitWindow))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.healthHandler)
	mux.HandleFunc("GET /panic", s.panicHandler)
	mux.Handle("GET /items", limited.ThenFunc(itemHandler.List))
	mux.Handle("POST /items", limited.ThenFunc(itemHandler.Create))
	mux.Handle("GET /items/{id}", limited.ThenFunc(itemHandler.Get))

	global := alice.New(
		middleware.RequestID,
		middleware.RequestLog,
		middleware.Recover,
		middleware.Timeout(cfg.RequestTimeout),
		middleware.SecureHeaders,
		http.NewCrossOriginProtection().Handler,
	)
	return global.Then(mux)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("test panic")
}
