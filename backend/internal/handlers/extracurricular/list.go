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
		// When auth is not available, allow handlers or routes to inject
		// a student_id into context via auth.WithStudentID. If still empty,
		// return an error and ask callers to provide `student_id` as a
		// query parameter when testing.
		return nil, errs.BadRequest("student id not found in request context")
	}

	return h.ExtracurricularRepository.ListExtracurriculars(ctx, studentID)
}
