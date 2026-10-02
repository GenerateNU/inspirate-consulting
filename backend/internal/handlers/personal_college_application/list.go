package personalcollegeapplication

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

func (h *Handler) ListPersonalCollegeApplicationsByStudentID(ctx context.Context) ([]models.PersonalCollegeApplication, error) {
	studentID := auth.GetStudentID(ctx)
	applications, err := h.PersonalCollegeApplicationRepository.ListPersonalCollegeApplicationsByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	return applications, nil
}