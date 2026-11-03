package extracurricular

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// CreateExtracurricular creates a new extracurricular entry using the authenticated student's id from context.
func (h *Handler) CreateExtracurricular(ctx context.Context, input *models.CreateExtracurricularInput) (*models.CreateExtracurricularOutput, error) {
	studentID, err := auth.GetStudentID(ctx)
	if err != nil {
		return nil, errs.Unauthorized()
	}
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, errs.Unauthorized()
	}

	return h.ExtracurricularRepository.CreateExtracurricular(ctx, studentID.String(), userID.String(), *input)
}
