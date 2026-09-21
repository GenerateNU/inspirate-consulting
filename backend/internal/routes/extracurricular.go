package routes

import (
	"context"
	"net/http"
	"os"

	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/handlers/extracurricular"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/auth"

	"github.com/danielgtaylor/huma/v2"
)

func SetUpExtracurricularRoutes(api huma.API, repository *data.Repository) {
	h := extracurricular.NewHandler(repository.Extracurricular)

	// POST /extracurriculars - create
	testMode := os.Getenv("TEST_MODE") == "true"

	type createInput struct {
		Body models.CreateExtracurricularRequest `body:""`
	}
	huma.Register(api, huma.Operation{
		OperationID: "create-extracurricular",
		Method:      http.MethodPost,
		Path:        "/extracurriculars",
		Description: "Create a new extracurricular for the authenticated student.",
		Tags:        []string{"Extracurriculars"},
	}, func(ctx context.Context, input *createInput) (*models.CreateExtracurricularOutput, error) {
		// adapt bound body into the handler's expected input type
		createIn := models.CreateExtracurricularInput{Body: input.Body}
		return h.CreateExtracurricular(ctx, &createIn)
	})

	// GET /extracurriculars - list for student
	// Accept an optional `student_id` query parameter for local testing.
	type listInput struct {
		StudentID string `query:"student_id"`
	}
	huma.Register(api, huma.Operation{
		OperationID: "list-extracurriculars",
		Method:      http.MethodGet,
		Path:        "/extracurriculars",
		Description: "List extracurriculars for the authenticated student.",
		Tags:        []string{"Extracurriculars"},
	}, func(ctx context.Context, input *listInput) (*models.ListExtracurricularsOutput, error) {
		if testMode && input != nil && input.StudentID != "" {
			ctx = auth.WithStudentID(ctx, input.StudentID)
		}
		list, err := h.ListExtracurriculars(ctx)
		if err != nil {
			return nil, err
		}
		return &models.ListExtracurricularsOutput{Body: list}, nil
	})

	// PATCH /extracurriculars/{id} - update
	type updateInput struct {
		ID   int64                              `path:"id"`
		Body models.UpdateExtracurricularRequest `body:""`
	}

	huma.Register(api, huma.Operation{
		OperationID: "update-extracurricular",
		Method:      http.MethodPatch,
		Path:        "/extracurriculars/{id}",
		Description: "Update fields on an extracurricular for the authenticated student.",
		Tags:        []string{"Extracurriculars"},
	}, func(ctx context.Context, input *updateInput) (*models.UpdateExtracurricularOutput, error) {
		updateIn := models.UpdateExtracurricularInput{ID: input.ID, Body: input.Body}
		updated, err := h.UpdateExtracurricular(ctx, input.ID, &updateIn)
		if err != nil {
			return nil, err
		}
		return &models.UpdateExtracurricularOutput{Body: *updated}, nil
	})
}
