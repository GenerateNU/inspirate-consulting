package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/handlers/essays"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetUpEssayRoutes(api huma.API, repository *data.Repository) {
	essayHandler := essays.NewHandler(repository.Essay)

	// GET route that shows all of a students essay
	huma.Register(api, huma.Operation{
		OperationID: "get-essays-from-student",
		Method:      http.MethodGet,
		Path:        "/students/{student_id}/essays",
		Description: "Shows all essays of a student",
		Tags:        []string{"Essays"},
	}, func(ctx context.Context, input *models.GetEssaysFromStudentInput) (*models.GetEssaysFromStudentOutput, error) {
		return essayHandler.GetEssaysFromStudent(ctx, input)
	})

	// POST route that inserts a row into the essays table
	huma.Register(api, huma.Operation{
		OperationID: "create-essay",
		Method:      http.MethodPost,
		Path:        "/essays",
		Description: "Creates an essay with the logged in student",
		Tags:        []string{"Essays"},
	}, func(ctx context.Context, input *models.CreateEssayInput) (*models.CreateEssayOutput, error) {
		return essayHandler.CreateEssay(ctx, input)
	})

	// PATCH route that updates a students essay status
	huma.Register(api, huma.Operation{
		OperationID: "update-essay-status",
		Method:      http.MethodPatch,
		Path:        "/essays/{essay_id}/status",
		Description: "updates the status of a students essay",
		Tags:        []string{"Essays"},
	}, func(ctx context.Context, input *models.UpdateStatusInput) (*models.UpdateStatusOutput, error) {
		return essayHandler.UpdateStatus(ctx, input)
	})
}
