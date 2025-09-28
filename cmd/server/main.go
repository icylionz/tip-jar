package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tipjar/internal/config"
	"tipjar/internal/database"
	"tipjar/internal/handlers"
	"tipjar/internal/auth"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize database
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize auth
	authService := auth.NewService(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.SessionSecret)

	// Initialize handlers (handles static files and all routes)
	h := handlers.New(db, authService, cfg)
	h.RegisterRoutes(e)

	slog.Info("Server starting", "port", cfg.Port)

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", cfg.Port)); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited")
}