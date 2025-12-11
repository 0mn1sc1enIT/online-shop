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
	// 1. Config
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// 2. Logger
	appLogger := logger.New(cfg.Logger.Level, cfg.Logger.SeqURL)
	appLogger.Info().Msg("Starting API...")

	// 3. Database
	db, err := database.New(cfg.Database)
	if err != nil {
		appLogger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// 4. Repositories
	repos := repository.NewRepositories(db)

	// 5. Services
	tokenTTL, _ := time.ParseDuration(cfg.Auth.TokenTTL)
	services := service.NewServices(service.Deps{
		Repos:     repos,
		TokenTTL:  tokenTTL,
		SignedKey: cfg.Auth.JWTSecret,
	})

	// 6. Handlers (Transport)
	handlers := transport.NewHandler(services, appLogger)

	// 7. Server Setup
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: handlers.InitRoutes(),
	}

	// Запуск сервера в горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			appLogger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	appLogger.Info().Msgf("Server started on port %s", cfg.Server.Port)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	appLogger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error().Err(err).Msg("Server forced to shutdown")
	}

	appLogger.Info().Msg("Server exited properly")
}
