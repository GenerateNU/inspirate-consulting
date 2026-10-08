package extracurricular

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// CreateExtracurricular creates a new extracurricular entry using the authenticated student's id from context.
func (h *Handler) CreateExtracurricular(ctx context.Context, input *models.CreateExtracurricularInput) (*models.CreateExtracurricularOutput, error) {
	userID := auth.GetUserID(ctx)
	if userID == "" {
		return nil, errs.BadRequest("user id not found in request context")
	}

	return h.ExtracurricularRepository.CreateExtracurricular(ctx, userID, *input)
}