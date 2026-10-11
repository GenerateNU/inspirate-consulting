package main

import (
	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/routes"
	"log"
	"os"
)

func main() {
	// SetupApp prepends PUBLIC_API_URL to the server list when it is set, so
	// clear it here: otherwise the committed spec would depend on whoever ran
	// this command and CI's up-to-date check would fail for unrelated reasons.
	os.Unsetenv("PUBLIC_API_URL")

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
