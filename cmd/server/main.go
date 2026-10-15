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

	"github.com/1337jazz/generoni/internal/adapter/config/environment"
	"github.com/1337jazz/generoni/internal/adapter/container"
	httpadapter "github.com/1337jazz/generoni/internal/adapter/http"
	"github.com/1337jazz/generoni/internal/adapter/logger"
)

func main() {

	// Bootstrap the config and logger
	cfg, err := environment.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid configuration: %v\n", err)
		os.Exit(1)
	}
	logger.Setup(cfg.LogLevel)

	// Run with graceful shutdown
	slog.Info("starting")
	if err := run(cfg); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func run(cfg environment.Config) error {

	// Configure the shutdown on SIGTERM and SIGINT
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Bootstrap dependencies
	cntr, err := container.New(ctx, cfg)
	if err != nil {
		return err
	}
	defer cntr.DB.Close()

	// Configure the server
	srv := httpadapter.NewServer(cfg, cntr)

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
		stop()
	}

	slog.Info("shutting down")

	// Shutdown the server, forcing a timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}
