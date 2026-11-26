package extracurricular

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

// ListExtracurriculars fetches every extracurricular for one student. Students
// get their own; counselors must pass the student's id and be assigned to them.
func (h *Handler) ListExtracurriculars(ctx context.Context, rawStudentID string) ([]models.Extracurricular, error) {
	requested, err := auth.ParseOptionalStudentID(rawStudentID)
	if err != nil {
		return nil, err
	}
	studentID, err := h.Authorizer.ResolveStudentID(ctx, requested)
	if err != nil {
		return nil, err
	}

	return h.ExtracurricularRepository.ListExtracurriculars(ctx, studentID.String())
}
