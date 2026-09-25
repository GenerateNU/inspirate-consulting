package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/handlers/student"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetUpStudentRoutes(api huma.API, repository *data.Repository) {
	studentHandler := student.NewHandler(repository.EssayReview)

	huma.Register(api, huma.Operation{
		OperationID: "set-student-review-balance",
		Method:      http.MethodPatch,
		Path:        "/students/{id}/review-balance",
		Description: "Overwrite a student's review balance with an absolute value, recording the change in the review ledger.",
		Tags:        []string{"Students"},
	}, func(ctx context.Context, input *models.SetStudentReviewBalanceInput) (*models.SetStudentReviewBalanceOutput, error) {
		updated, err := studentHandler.SetReviewBalance(ctx, input.ID, input.Body)
		if err != nil {
			return nil, err
		}
		return &models.SetStudentReviewBalanceOutput{Body: *updated}, nil
	})
}
