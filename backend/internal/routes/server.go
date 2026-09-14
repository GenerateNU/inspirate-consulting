package routes

import (
	"net/http"

	storage "github.com/GenerateNU/inspirate-consulting/internal/data"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
)

func SetUpProtectedHumaRoutes(api huma.API, store *storage.Store) error {
	SetUpGreetingRoutes(api, store)
	return nil
}

func SetUpApp() (*fiber.App, huma.API, error) {
	// Create new router and API
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://localhost:5174",
			"http://127.0.0.1:5173",
			"http://127.0.0.1:5174",
		},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))
	config := huma.DefaultConfig("My API", "1.0.0")
	config.Servers = []*huma.Server{
		{URL: "http://127.0.0.1:3001"},
	}
	api := humafiber.New(app, config)

	store := &storage.Store{}
	SetUpProtectedHumaRoutes(api, store)

	return app, api, nil
}
