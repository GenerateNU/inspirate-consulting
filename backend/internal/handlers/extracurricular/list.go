package extracurricular

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// ListExtracurriculars fetches every extracurricular for the authenticated student.
func (h *Handler) ListExtracurriculars(ctx context.Context) ([]models.Extracurricular, error) {
	studentID, err := auth.GetStudentID(ctx)
	if err != nil {
		return nil, errs.Unauthorized()
	}

	return h.ExtracurricularRepository.ListExtracurriculars(ctx, studentID.String())
}
