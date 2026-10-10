package extracurricular

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// ListExtracurriculars fetches every extracurricular for the authenticated student.
func (h *Handler) ListExtracurriculars(ctx context.Context) ([]models.Extracurricular, error) {
	studentID := auth.GetStudentID(ctx)
	if studentID == "" {
		return nil, errs.BadRequest("student id not found in request context")
	}

	return h.ExtracurricularRepository.ListExtracurriculars(ctx, studentID)
}
