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
// list pending reviews, list students, update counselor, fetch all counselors, and get a counselor
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

	// Register PATCH /counselors handler.
	huma.Register(api, huma.Operation{
		OperationID: "update-counselor",
		Middlewares: huma.Middlewares{verifyRole(api, models.CounselorRole)},
		Method:      http.MethodPatch,
		Path:        "/counselors",
		Description: "Update the authenticated counselor's custom link and account expiration. Omit a field to keep it, send null to clear it.",
		Tags:        []string{"Counselors"},
	}, func(ctx context.Context, input *models.UpdateCounselorInput) (*models.UpdateCounselorOutput, error) {
		updatedCounselor, err := counselorHandler.UpdateCounselor(ctx, input.Body)
		if err != nil {
			return nil, err
		}
		return &models.UpdateCounselorOutput{Body: updatedCounselor}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "fetch-counselors",
		Middlewares: huma.Middlewares{verifyRole(api, models.CounselorRole)},
		Method:      http.MethodGet,
		Path:        "/counselors/fetchall",
		Description: "Fetch a page of counselors, optionally filtered by user_id and name",
		Tags:        []string{"Counselors"},
	}, func(ctx context.Context, input *models.FetchCounselorsInput) (*models.FetchCounselorsOutput, error) {
		page, err := counselorHandler.ListCounselors(ctx, input)
		if err != nil {
			return nil, err
		}
		return &models.FetchCounselorsOutput{Body: page}, nil
	})

	// Registered after the static /counselors/* routes so they aren't matched as an ID
	huma.Register(api, huma.Operation{
		OperationID: "get-counselor",
		Method:      http.MethodGet,
		Path:        "/counselors/{id}",
		Description: "Get a counselor's profile",
		Tags:        []string{"Counselors"},
	}, func(ctx context.Context, input *models.GetCounselorInput) (*models.GetCounselorOutput, error) {
		fetchedCounselor, err := counselorHandler.GetCounselor(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		return &models.GetCounselorOutput{Body: fetchedCounselor}, nil
	})
}
