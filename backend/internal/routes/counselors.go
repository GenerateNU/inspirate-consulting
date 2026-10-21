package routes

import (
	"context"
	"inspirate-consulting/internal/auth"
	"net/http"

	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/handlers/counselor"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

// SetUpCounselorRoutes registers the counselor endpoints:
// list pending reviews and list students
func SetUpCounselorRoutes(api huma.API, repository *data.Repository, verifyRole auth.RoleVerifier) {
	counselorHandler := counselor.NewHandler(repository.Counselor)

	// Register GET /counselors/reviews handler.
	huma.Register(api, huma.Operation{
		OperationID: "list-counselor-pending-reviews",
		Middlewares: huma.Middlewares{verifyRole(api, models.CounselorRole)},
		Method:      http.MethodGet,
		Path:        "/counselors/reviews",
		Description: "List the open essay reviews the authenticated counselor needs to complete, oldest first.",
		Tags:        []string{"Counselors"},
	}, func(ctx context.Context, input *models.ListPendingReviewsInput) (*models.ListPendingReviewsOutput, error) {
		reviews, err := counselorHandler.ListPendingReviews(ctx)
		if err != nil {
			return nil, err
		}
		return &models.ListPendingReviewsOutput{Body: reviews}, nil
	})

	// Register GET /counselors/students handler.
	huma.Register(api, huma.Operation{
		OperationID: "list-counselor-students",
		Middlewares: huma.Middlewares{verifyRole(api, models.CounselorRole)},
		Method:      http.MethodGet,
		Path:        "/counselors/students",
		Description: "List the authenticated counselor's students, flagged for falling behind and pending essay reviews.",
		Tags:        []string{"Counselors"},
	}, func(ctx context.Context, input *models.ListCounselorStudentsInput) (*models.ListCounselorStudentsOutput, error) {
		students, err := counselorHandler.ListStudents(ctx, *input)
		if err != nil {
			return nil, err
		}
		return &models.ListCounselorStudentsOutput{Body: students}, nil
	})
}
