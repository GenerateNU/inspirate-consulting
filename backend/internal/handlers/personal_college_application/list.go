package personalcollegeapplication

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

func (h *Handler) ListPersonalCollegeApplicationsByStudentID(ctx context.Context) ([]models.PersonalCollegeApplication, error) {
	studentID, err := auth.GetStudentID(ctx)
	if err != nil {
		return nil, err
	}
	applications, err := h.PersonalCollegeApplicationRepository.ListPersonalCollegeApplicationsByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	return applications, nil
}
