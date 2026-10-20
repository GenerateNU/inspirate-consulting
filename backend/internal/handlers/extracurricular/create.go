package extracurricular

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// CreateExtracurricular creates a new extracurricular entry using the authenticated student's id from context.
func (h *Handler) CreateExtracurricular(ctx context.Context, input *models.CreateExtracurricularInput) (*models.CreateExtracurricularOutput, error) {
	studentID, err := auth.ParseStudentID(input.StudentID)
	if err != nil {
		return nil, err
	}
	if err := h.Authorizer.CanAccessStudent(ctx, studentID); err != nil {
		return nil, err
	}
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, errs.Unauthorized()
	}

	return h.ExtracurricularRepository.CreateExtracurricular(ctx, userID.String(), *input)
}
