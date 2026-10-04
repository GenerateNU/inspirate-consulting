package main

import (
	"context"
	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/routes"

	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sethvargo/go-envconfig"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize application with config
	app, err := routes.InitApp(*cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Close database connection when main exits
	defer func() {
		slog.Info("Closing database connection")
		if err := app.Repo.Close(); err != nil {
			slog.Error("failed to close database", "error", err)
		}
	}()

	port := cfg.Application.Port

	// Listen for connections with a goroutine
	go func() {
		if err := app.Server.Listen(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for termination signal (SIGINT or SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	slog.Info("Shutting down server")

	// Shutdown server with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := app.Server.ShutdownWithContext(shutdownCtx); err != nil {
		slog.Error("failed to shutdown server gracefully", "error", err)
	}

	slog.Info("Server shutdown complete")
}

func LoadConfig() (*config.Config, error) {
	testMode := os.Getenv("TEST_MODE")

	var cfg config.Config
	// Load configuration from environment variables for production
	err := envconfig.Process(context.Background(), &cfg)
	if err != nil {
		log.Fatalln("Error processing environment variables: ", err)
	}

	cfg.TestMode = testMode == "true"

	return &cfg, nil
}
