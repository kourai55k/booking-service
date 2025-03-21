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
	"github.com/kourai55k/booking-service/internal/data/postgres"
	"github.com/kourai55k/booking-service/internal/domain/models"
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

	log.Info("starting application")

	// test
	pgPool, err := postgres.ConnectPool(context.Background(), cfg.PostgresConnString)
	if err != nil {
		log.Error("failed to connect to database", "err", err.Error())
	} else {
		log.Info("connected to database successfully")
	}

	userRepo := postgres.NewUserRepo(pgPool)
	err = userRepo.CreateUserTable()
	if err != nil {
		log.Error("failed to create user table", "err", err.Error())
	} else {
		log.Info("created user table successfully")
	}
	userService := service.NewUserService(userRepo)
	httpUserHandler := userHandler.NewUserHandler(userService, log)
	r := router.NewRouter(httpUserHandler)

	// test users
	userRepo.CreateUser(&models.User{Name: "name1", Login: "login1", HashPass: "hashpass1"})
	userRepo.CreateUser(&models.User{Name: "name2", Login: "login2", HashPass: "hashpass2"})

	log.Info("created test users")

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		log.Info("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Info("ListenAndServe error:", "err", err.Error())
			stop <- os.Interrupt
			return
		}
	}()

	// Block until we receive a termination signal
	<-stop
	log.Info("Shutting down server...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server shutdown error:", "err", err.Error())
	}

	log.Info("Server stopped gracefully")
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
