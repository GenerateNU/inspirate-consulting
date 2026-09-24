package routes

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/data/postgres"
	"inspirate-consulting/internal/errs"
	"os"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	go_json "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type App struct {
	Server *fiber.App
	Repo   *data.Repository
	API    huma.API
}

// Initialize the App union type containing a fiber app and repository.
func InitApp(config config.Config) (*App, error) {
	ctx := context.Background()
	repo := postgres.NewRepository(ctx, config.DB)

	app, humaAPI, err := SetupApp(config, repo)
	if err != nil {
		return nil, err
	}
	return &App{
		Server: app,
		Repo:   repo,
		API:    humaAPI,
	}, nil
}

// Setup the fiber app with the specified configuration and database.
func SetupApp(config config.Config, repo *data.Repository) (*fiber.App, huma.API, error) {
	app := fiber.New(fiber.Config{
		JSONEncoder:  go_json.Marshal,
		JSONDecoder:  go_json.Unmarshal,
		ErrorHandler: errs.ErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(favicon.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(logger.New())

	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000,http://localhost:8080,https://cdn.scalar.com,http://127.0.0.1:8080,http://10.0.2.2:8080,http://localhost:5173,http://localhost"
	}
	splitAllowedOrigins := strings.Split(allowedOrigins, ",")

	app.Use(cors.New(cors.Config{
		AllowOrigins:     splitAllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
	}))
	// Create Huma API with OpenAPI configuration
	humaConfig := huma.DefaultConfig("Inspirate Consulting API", "1.0.0")
	humaConfig.Info.Description = "API for the Inspirate Consulting application"
	humaConfig.Info.Contact = &huma.Contact{
		Name: "Inspirate Consulting Team",
	}
	humaConfig.Servers = []*huma.Server{
		{URL: "http://localhost:8080", Description: "Local development server"},
	}

	humaAPI := humafiber.New(app, humaConfig)

	// Register public routes BEFORE auth middleware
	// routes.SetupAuthRoutes(humaAPI, repo, config)

	// Apply auth middleware — only affects routes registered after this point
	if !config.TestMode {
		humaAPI.UseMiddleware(auth.AuthMiddleware(humaAPI, &config.Supabase))
	}

	// Documentation routes (Huma provides built-in docs at /docs and /openapi.json)
	// setupDocsRoutes(app, "/app/api")

	// Root route
	app.Get("/", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Welcome to Inspirate Consulting!")
	})

	// Register protected Huma endpoints
	if err := setupProtectedHumaRoutes(humaAPI, repo, config); err != nil {
		return nil, nil, err
	}

	return app, humaAPI, nil
}

// Setup protected Huma routes (behind auth middleware)
func setupProtectedHumaRoutes(api huma.API, repo *data.Repository, config config.Config) error {
	// Attach each of the routes to the API here
	SetUpGreetingRoutes(api, repo)
	SetUpGlobalCollegeRoutes(api, repo)
	SetUpPersonalCollegeApplicationRoutes(api, repo)
	SetUpEssayReviewRoutes(api, repo)
	SetUpStudentRoutes(api, repo)
	return nil
}
