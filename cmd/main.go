package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OmniscienIT/GolangAPI/config"
	"github.com/OmniscienIT/GolangAPI/internal/pkg/database"
	"github.com/OmniscienIT/GolangAPI/internal/pkg/logger"
	"github.com/OmniscienIT/GolangAPI/internal/repository"
	"github.com/OmniscienIT/GolangAPI/internal/service"
	"github.com/OmniscienIT/GolangAPI/internal/transport"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	appLogger := logger.New(cfg.Logger.Level, cfg.Logger.SeqURL)

	appLogger.Info().
		Str("env", cfg.Server.Mode).
		Msg("Starting API...")

	db, err := database.New(cfg.Database)
	if err != nil {
		appLogger.Fatal().
			Err(err).
			Str("db_host", cfg.Database.Host).
			Msg("Failed to connect to database")
	}

	appLogger.Info().
		Str("db_name", cfg.Database.DBName).
		Msg("Database connection established successfully")

	repos := repository.NewRepositories(db, appLogger)

	tokenTTL, err := time.ParseDuration(cfg.Auth.TokenTTL)
	if err != nil {
		appLogger.Warn().
			Err(err).
			Msg("Failed to parse TokenTTL from config, using default 24h")
		tokenTTL = 24 * time.Hour
	}
	services := service.NewServices(service.Deps{
		Repos:     repos,
		TokenTTL:  tokenTTL,
		SignedKey: cfg.Auth.JWTSecret,
		Logger:    appLogger,
	})
	appLogger.Debug().Msg("Services and Repositories initialized")

	handlers := transport.NewHandler(services, appLogger)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: handlers.InitRoutes(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			appLogger.Fatal().
				Err(err).
				Msg("Failed to start server")
		}
	}()

	appLogger.Info().
		Str("port", cfg.Server.Port).
		Msg("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit

	appLogger.Info().
		Str("signal", sig.String()).
		Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error().
			Err(err).
			Msg("Server forced to shutdown")
	}

	appLogger.Info().Msg("Server exited properly")
}
