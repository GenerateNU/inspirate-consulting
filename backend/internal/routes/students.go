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
		OperationID: "fetch-students",
		Middlewares: huma.Middlewares{verifyRole(api, models.CounselorRole)},
		Method:      http.MethodGet,
		Path:        "/students",
		Description: "Fetch a page of students, optionally filtered by user_id and name",
		Tags:        []string{"Students"},
	}, func(ctx context.Context, input *models.FetchStudentsInput) (*models.FetchStudentsOutput, error) {
		page, err := studentHandler.ListStudents(ctx, input)
		if err != nil {
			return nil, err
		}
		return &models.FetchStudentsOutput{Body: page}, nil
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

	huma.Register(api, huma.Operation{
		OperationID: "assign-counselor",
		Middlewares: huma.Middlewares{verifyRole(api, models.CounselorRole)},
		Method:      http.MethodPatch,
		Path:        "/students/assign",
		Description: "Assign a counselor to a particular student",
		Tags:        []string{"Counselors"},
	}, func(ctx context.Context, input *models.AssignCounselorInput) (*models.AssignCounselorOutput, error) {
		updated, err := studentHandler.AssignCounselor(ctx, input)
		if err != nil {
			return nil, err
		}
		return &models.AssignCounselorOutput{Body: *updated}, nil
	})
}
