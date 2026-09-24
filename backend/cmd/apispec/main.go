package main

import (
	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/routes"
	"log"
	"os"
)

func main() {
	cfg := config.Config{
		TestMode: true,
	}

	_, humaAPI, err := routes.SetupApp(cfg, &data.Repository{})
	if err != nil {
		log.Fatalf("failed to set up app: %v", err)
	}

	b, err := humaAPI.OpenAPI().YAML()
	if err != nil {
		log.Fatalf("failed to generate spec: %v", err)
	}

	if err := os.WriteFile("internal/api/openapi.yaml", b, 0644); err != nil {
		log.Fatalf("failed to write spec: %v", err)
	}
}
