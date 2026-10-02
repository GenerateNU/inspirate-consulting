package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/data"
	essayreview "inspirate-consulting/internal/handlers/essay_review"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetUpEssayReviewRoutes(api huma.API, repository *data.Repository) {
	essayReviewHandler := essayreview.NewHandler(repository.EssayReview)

	huma.Register(api, huma.Operation{
		OperationID: "request-essay-review",
		Method:      http.MethodPost,
		Path:        "/reviews",
		Description: "Spend a student's review balance to open a review on an essay.",
		Tags:        []string{"Essay Reviews"},
	}, func(ctx context.Context, input *models.RequestEssayReviewInput) (*models.RequestEssayReviewOutput, error) {
		requested, err := essayReviewHandler.RequestEssayReview(ctx, input.Body)
		if err != nil {
			return nil, err
		}
		return &models.RequestEssayReviewOutput{Body: *requested}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "refund-essay-review",
		Method:      http.MethodPost,
		Path:        "/reviews/{id}/refund",
		Description: "Refund a review that has not been completed, crediting the balance back to the student.",
		Tags:        []string{"Essay Reviews"},
	}, func(ctx context.Context, input *models.RefundEssayReviewInput) (*models.RefundEssayReviewOutput, error) {
		refunded, err := essayReviewHandler.RefundEssayReview(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		return &models.RefundEssayReviewOutput{Body: *refunded}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "complete-essay-review",
		Method:      http.MethodPost,
		Path:        "/reviews/{id}/complete",
		Description: "Mark a review as having been completed by a counselor.",
		Tags:        []string{"Essay Reviews"},
	}, func(ctx context.Context, input *models.CompleteEssayReviewInput) (*models.CompleteEssayReviewOutput, error) {
		completed, err := essayReviewHandler.CompleteEssayReview(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		return &models.CompleteEssayReviewOutput{Body: *completed}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-essay-review-status",
		Method:      http.MethodGet,
		Path:        "/essays/{essay_id}/review-status",
		Description: "Get the current review status of an essay, from its most recent review request.",
		Tags:        []string{"Essay Reviews"},
	}, func(ctx context.Context, input *models.GetEssayReviewStatusInput) (*models.GetEssayReviewStatusOutput, error) {
		status, err := essayReviewHandler.GetEssayReviewStatus(ctx, input.EssayID)
		if err != nil {
			return nil, err
		}
		return &models.GetEssayReviewStatusOutput{Body: *status}, nil
	})
}
