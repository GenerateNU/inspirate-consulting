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
		created, err := globalCollegeHandler.CreateGlobalCollege(ctx, input.Body)
		if err != nil {
			return nil, err
		}
		return &models.CreateGlobalCollegeOutput{Body: *created}, nil
	})

	// Register GET /colleges/{id} handler.
	huma.Register(api, huma.Operation{
		OperationID: "get-global-college",
		Method:      http.MethodGet,
		Path:        "/colleges/{id}",
		Description: "Get a global college by ID.",
		Tags:        []string{"Global Colleges"},
	}, func(ctx context.Context, input *models.GetGlobalCollegeInput) (*models.GetGlobalCollegeOutput, error) {
		college, err := globalCollegeHandler.GetGlobalCollege(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		return &models.GetGlobalCollegeOutput{Body: *college}, nil
	})

	// Register GET /colleges handler.
	huma.Register(api, huma.Operation{
		OperationID: "list-global-colleges",
		Method:      http.MethodGet,
		Path:        "/colleges",
		Description: "List all global colleges.",
		Tags:        []string{"Global Colleges"},
	}, func(ctx context.Context, input *models.ListGlobalCollegesInput) (*models.ListGlobalCollegesOutput, error) {
		colleges, err := globalCollegeHandler.ListGlobalColleges(ctx)
		if err != nil {
			return nil, err
		}
		return &models.ListGlobalCollegesOutput{Body: colleges}, nil
	})

	// Register PUT /colleges/{id} handler.
	huma.Register(api, huma.Operation{
		OperationID: "update-global-college",
		Method:      http.MethodPut,
		Path:        "/colleges/{id}",
		Description: "Update a global college's details.",
		Tags:        []string{"Global Colleges"},
	}, func(ctx context.Context, input *models.UpdateGlobalCollegeInput) (*models.UpdateGlobalCollegeOutput, error) {
		updated, err := globalCollegeHandler.UpdateGlobalCollege(ctx, input.ID, input.Body)
		if err != nil {
			return nil, err
		}
		return &models.UpdateGlobalCollegeOutput{Body: *updated}, nil
	})

	// Register DELETE /colleges/{id} handler.
	huma.Register(api, huma.Operation{
		OperationID: "delete-global-college",
		Method:      http.MethodDelete,
		Path:        "/colleges/{id}",
		Description: "Delete a global college.",
		Tags:        []string{"Global Colleges"},
	}, func(ctx context.Context, input *models.DeleteGlobalCollegeInput) (*models.DeleteGlobalCollegeOutput, error) {
		err := globalCollegeHandler.DeleteGlobalCollege(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		return &models.DeleteGlobalCollegeOutput{}, nil
	})

	// Register GET /colleges/search handler
	huma.Register(api, huma.Operation{
		OperationID: "search-global-colleges",
		Method:      http.MethodGet,
		Path:        "/college-search",
		Description: "Search global colleges by name or location.",
		Tags:        []string{"Global Colleges"},
	}, func(ctx context.Context, input *models.SearchGlobalCollegesInput) (*models.SearchGlobalCollegesOutput, error) {
		results, err := globalCollegeHandler.SearchGlobalColleges(ctx, input.Query)
		if err != nil {
			return nil, err
		}
		return &models.SearchGlobalCollegesOutput{Body: results}, nil
	})
}
