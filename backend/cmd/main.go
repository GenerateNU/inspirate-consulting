package main

import (
	"context"
	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/routes"
	"inspirate-consulting/internal/supabase"

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

	// Use a concrete struct for envconfig since SupabaseInterface can't be loaded from env vars
	type envConf struct {
		Application config.Application
		DB          config.DB
		Supabase    supabase.Supabase
		S3          config.S3
	}
	var env envConf
	err := envconfig.Process(context.Background(), &env)
	if err != nil {
		log.Fatalln("Error processing environment variables: ", err)
	}

	return &config.Config{
		Application: env.Application,
		DB:          env.DB,
		Supabase:    &env.Supabase,
		S3:          env.S3,
		TestMode:    testMode == "true",
	}, nil
}
