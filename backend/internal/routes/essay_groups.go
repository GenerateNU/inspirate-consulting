package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/handlers/essay_groups"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetUpEssayGroupRoutes(api huma.API, repository *data.Repository) {
	essayGroupHandler := essay_groups.NewHandler(repository.EssayGroup)

	// POST route that inserts a row into the essay_groups table
	huma.Register(api, huma.Operation{
		OperationID: "create-essay-group",
		Method:      http.MethodPost,
		Path:        "/essay-groups",
		Description: "Creates an essay group for a student",
		Tags:        []string{"Essay Groups"},
	}, func(ctx context.Context, input *models.CreateEssayGroupInput) (*models.CreateEssayGroupOutput, error) {
		return essayGroupHandler.CreateEssayGroup(ctx, input)
	})

	// GET route that shows all of a students essay groups
	huma.Register(api, huma.Operation{
		OperationID: "get-essay-groups-from-student",
		Method:      http.MethodGet,
		Path:        "/students/{student_id}/essay-groups",
		Description: "Shows all essay groups of a student",
		Tags:        []string{"Essay Groups"},
	}, func(ctx context.Context, input *models.GetEssayGroupsInput) (*models.GetEssayGroupsOutput, error) {
		return essayGroupHandler.GetEssayGroups(ctx, input)
	})
}
