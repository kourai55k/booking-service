// This file is used to test and simulate the application's workflow during local development.
// It uses an in-memory repository instead of a real database,
// along with other components for quick testing and debugging.
// This file is **not intended for production use**.
//
// Make sure to exclude this file from production builds and deployments.
// It should be used solely for development and testing purposes.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kourai55k/booking-service/internal/config"
	"github.com/kourai55k/booking-service/internal/data"
	"github.com/kourai55k/booking-service/internal/service"
	"github.com/kourai55k/booking-service/internal/transport/handlers/http/router"
	"github.com/kourai55k/booking-service/internal/transport/handlers/http/userHandler"
	prettyslog "github.com/kourai55k/booking-service/pkg/prettySlog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// setup logger
	log := setupLogger(cfg.Env)
	log.Debug("config loaded", "config", cfg)
	log.Debug("logger initialized")

	// DI
	userRepo := data.NewInMemoryUserRepo()
	userService := service.NewUserService(userRepo)
	httpUserHandler := userHandler.NewUserHandler(userService, log)
	r := router.NewRouter(httpUserHandler)
	// TODO: use config file to configure server
	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	log.Debug("dependencies injected")

	// Channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("ListenAndServe error:", "err", err.Error())
			stop <- os.Interrupt
			return
		}
	}()
	log.Debug("server started", "addr", server.Addr)
	log.Info("app started")

	// Block until we receive a termination signal
	<-stop

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Error("server shutdown error:", "err", err.Error())
	}

	log.Debug("server stopped gracefully")
	log.Info("app stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := prettyslog.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
