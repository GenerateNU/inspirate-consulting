package routes

import (
	"context"
	"inspirate-consulting/internal/auth"
	"net/http"

	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/handlers/student"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetUpStudentRoutes(api huma.API, repository *data.Repository, verifyRole auth.RoleVerifier, authorizer *auth.Authorizer) {
	studentHandler := student.NewHandler(repository.EssayReview, repository.Student, authorizer)

	huma.Register(api, huma.Operation{
		OperationID: "set-student-review-balance",
		Middlewares: huma.Middlewares{verifyRole(api, models.CounselorRole)},
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

	huma.Register(api, huma.Operation{
		OperationID: "get-student",
		Method:      http.MethodGet,
		Path:        "/students/{id}",
		Description: "Get a student",
		Tags:        []string{"Students"},
	}, func(ctx context.Context, input *models.GetStudentInput) (*models.GetStudentOutput, error) {
		fetchedStudent, err := studentHandler.GetStudent(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		return &models.GetStudentOutput{Body: fetchedStudent}, nil
	})
}
