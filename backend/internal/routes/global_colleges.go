package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/data"
	globalcollege "inspirate-consulting/internal/handlers/global_college"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

// SetUpGlobalCollegeRoutes registers the global college endpoints:
// create, get by id, and list.
func SetUpGlobalCollegeRoutes(api huma.API, repository *data.Repository) {
	globalCollegeHandler := globalcollege.NewHandler(repository.GlobalCollege)

	// Register POST /colleges handler.
	huma.Register(api, huma.Operation{
		OperationID: "create-global-college",
		Method:      http.MethodPost,
		Path:        "/colleges",
		Description: "Create a global college with a name, location, and optional EA/ED/RD deadlines.",
		Tags:        []string{"Global Colleges"},
	}, func(ctx context.Context, input *models.CreateGlobalCollegeInput) (*models.CreateGlobalCollegeOutput, error) {
		created, err := globalCollegeHandler.CreateGlobalCollege(ctx, input)
		if err != nil {
			return nil, err
		}
		return created, nil
	})

	// Register GET /colleges/{id} handler.
	huma.Register(api, huma.Operation{
		OperationID: "get-global-college",
		Method:      http.MethodGet,
		Path:        "/colleges/{id}",
		Description: "Get a global college by ID.",
		Tags:        []string{"Global Colleges"},
	}, func(ctx context.Context, input *models.GetGlobalCollegeInput) (*models.GetGlobalCollegeOutput, error) {
		college, err := globalCollegeHandler.GetGlobalCollege(ctx, input)
		if err != nil {
			return nil, err
		}
		return college, nil
	})

	// Register GET /colleges handler.
	huma.Register(api, huma.Operation{
		OperationID: "list-global-colleges",
		Method:      http.MethodGet,
		Path:        "/colleges",
		Description: "List all global colleges.",
		Tags:        []string{"Global Colleges"},
	}, func(ctx context.Context, input *models.ListGlobalCollegesInput) (*models.ListGlobalCollegesOutput, error) {
		colleges, err := globalCollegeHandler.ListGlobalColleges(ctx, input)
		if err != nil {
			return nil, err
		}
		return colleges, nil
	})
}
