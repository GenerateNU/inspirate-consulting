package routes

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// HealthOutput reports whether the API process is up and serving traffic.
type HealthOutput struct {
	Body struct {
		Status string `json:"status" example:"ok" doc:"Liveness status of the API"`
	}
}

// SetUpHealthRoutes registers the unauthenticated health check the deploy
// workflow polls to decide whether a release is live. The path must stay in
// sync with the skipPaths whitelist in auth.AuthMiddleware.
func SetUpHealthRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "health-check",
		Method:      http.MethodGet,
		Path:        "/api/v1/health",
		Description: "Report whether the API is running.",
		Tags:        []string{"Health"},
	}, func(ctx context.Context, _ *struct{}) (*HealthOutput, error) {
		resp := &HealthOutput{}
		resp.Body.Status = "ok"
		return resp, nil
	})
}
