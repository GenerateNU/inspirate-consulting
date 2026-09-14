package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/GenerateNU/inspirate-consulting/internal/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.yaml.in/yaml/v3"
)

func main() {
	// Load config

	// Start the fiber app with the config

	app, api, err := routes.SetUpApp()

	// Root route
	app.Get("/", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Welcome to Inspirate Consulting!")
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to setup app: %v\n", err)
		os.Exit(1)
	}

	// Get OpenAPI spec
	openAPI := api.OpenAPI()

	// Create api directory if it doesn't exist
	apiDir := "api"
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create api directory: %v\n", err)
		os.Exit(1)
	}

	// Write YAML file
	yamlPath := filepath.Join(apiDir, "openapi.yaml")
	yamlFile, err := os.Create(yamlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create YAML file: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := yamlFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close YAML file: %v\n", err)
		}
	}()

	encoder := yaml.NewEncoder(yamlFile)
	encoder.SetIndent(2)
	if err := encoder.Encode(openAPI); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to encode OpenAPI spec: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ OpenAPI spec generated: %s\n", yamlPath)

	connString := os.Getenv("SUPABASE_DATABASE_URL")
	if connString == "" {
		connString = os.Getenv("DATABASE_URL")
	}
	if connString == "" {
		connString = "postgresql://postgres:postgres@127.0.0.1:54322/postgres"
	}

	ctx := context.Background()

	// Create a connection pool
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	// Verify connectivity
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	log.Println("Successfully connected to Supabase via pgxpool")

	log.Fatal(app.Listen(":3001"))

	// Start the server!
	// http.ListenAndServe("127.0.0.1:8888", router)
}
