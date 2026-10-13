package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/handlers/extracurricular"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetUpExtracurricularRoutes(api huma.API, repository *data.Repository) {
	extracurricularHandler := extracurricular.NewHandler(repository.Extracurricular)

	huma.Register(api, huma.Operation{
		OperationID: "create-extracurricular",
		Method:      http.MethodPost,
		Path:        "/extracurriculars/{studentID}",
		Description: "Create a new extracurricular for the given student.",
		Tags:        []string{"Extracurriculars"},
	}, extracurricularHandler.CreateExtracurricular)

	huma.Register(api, huma.Operation{
		OperationID: "list-extracurriculars",
		Method:      http.MethodGet,
		Path:        "/extracurriculars",
		Description: "List extracurriculars for the authenticated student.",
		Tags:        []string{"Extracurriculars"},
	}, func(ctx context.Context, _ *models.ListExtracurricularsInput) (*models.ListExtracurricularsOutput, error) {
		list, err := extracurricularHandler.ListExtracurriculars(ctx)
		if err != nil {
			return nil, err
		}
		return &models.ListExtracurricularsOutput{Body: list}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-extracurricular",
		Method:      http.MethodPatch,
		Path:        "/extracurriculars/{id}",
		Description: "Update fields on an extracurricular for the authenticated student.",
		Tags:        []string{"Extracurriculars"},
	}, func(ctx context.Context, input *models.UpdateExtracurricularInput) (*models.UpdateExtracurricularOutput, error) {
		updated, err := extracurricularHandler.UpdateExtracurricular(ctx, input.ID, input)
		if err != nil {
			return nil, err
		}
		return &models.UpdateExtracurricularOutput{Body: *updated}, nil
	})
}
