package accounts

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cex/internal/accounts"
	"cex/pkg/cfg"
)

func main() {
	// 1) Load & validate config
	cfg.Init()

	// 2) Bootstrap the Accounts HTTP server
	e, err := accounts.NewServer()
	if err != nil {
		log.Fatalf("startup failed: %v", err)
	}

	// 3) Run server in background
	srvErr := make(chan error, 1)
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Cfg.Accounts.Port)
		e.Logger.Info("starting accounts service", "addr", addr)
		srvErr <- e.Start(addr)
	}()

	// 4) Wait for SIGINT/SIGTERM or a server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		e.Logger.Info("shutting down", "signal", sig.String())
	case err := <-srvErr:
		e.Logger.Error("server error", "error", err)
	}

	// 5) Graceful shutdown with 10s timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Error("shutdown error", "error", err)
	}
	e.Logger.Info("shutdown complete")
}
