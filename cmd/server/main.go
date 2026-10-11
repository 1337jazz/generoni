package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/1337jazz/generoni/internal/config"
	"github.com/1337jazz/generoni/internal/handlers"
	"github.com/1337jazz/generoni/internal/logging"
	"github.com/1337jazz/generoni/internal/middleware"
	"github.com/1337jazz/generoni/web"
	"github.com/justinas/alice"
)

func main() {

	// Bootstrap the config and logger
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid configuration: %v\n", err)
		os.Exit(1)
	}
	logging.Init(cfg.LogLevel)

	// Run with graceful shutdown
	if err := run(cfg); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func run(cfg config.Config) error {

	// Configure the server
	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           routes(cfg),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
	}

	// Configure the shutdown on SIGTERM and SIGINT
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Configure the listener and start the server
	listener, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return err
	}
	slog.Info("listening", "addr", listener.Addr().String())

	errc := make(chan error, 1)
	go func() {
		if err := srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errc <- err
		}
	}()

	// Wait for shutdown/error
	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}

	slog.Info("shutting down")

	// Shutdown the server, forcing a timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}

// routes builds the handler for the whole application. The global chain
// wraps every route, and the limited chain adds rate limiting to the routes
// that need it
func routes(cfg config.Config) http.Handler {
	limited := alice.New(middleware.RateLimit(cfg.RateLimit, cfg.RateLimitWindow))
	static := alice.New(middleware.NoDirList, middleware.CacheControl(cfg.StaticCacheControl))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.Handle("GET /static/", static.Then(http.FileServerFS(web.Static)))
	mux.Handle("GET /demo", limited.ThenFunc(handlers.Demo))
	mux.Handle("GET /panic", limited.ThenFunc(handlers.Panic))
	mux.Handle("GET /{$}", limited.ThenFunc(handlers.Home))
	mux.Handle("GET /", limited.ThenFunc(handlers.NotFound))

	global := alice.New(
		middleware.Recover,
		middleware.RequestID,
		middleware.RequestLog,
		middleware.SecureHeaders,
	)
	return global.Then(mux)
}
