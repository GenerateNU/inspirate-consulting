package extracurricular

import (
	"context"
	"log/slog"

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

	created, err := h.ExtracurricularRepository.CreateExtracurricular(ctx, studentID.String(), userID.String(), *input)
	if err != nil {
		return nil, err
	}

	if err := h.ActivityRepository.CreateActivity(ctx, studentID, models.CreateActivityBody{
		Action:            models.ActionAddedExtracurricular,
		ExtracurricularID: &created.Body.ID,
	}); err != nil {
		slog.Error("failed to log added_extracurricular activity", "student_id", studentID, "extracurricular_id", created.Body.ID, "error", err)
	}

	return created, nil
}
