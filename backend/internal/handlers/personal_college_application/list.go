package personalcollegeapplication

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

// lists a student's college list; students get their own, counselors pass one of their students' ids
func (h *Handler) ListPersonalCollegeApplicationsByStudentID(ctx context.Context, rawStudentID string) ([]models.PersonalCollegeApplication, error) {
	requested, err := auth.ParseOptionalStudentID(rawStudentID)
	if err != nil {
		return nil, err
	}
	studentID, err := h.Authorizer.ResolveStudentID(ctx, requested)
	if err != nil {
		return nil, err
	}
	applications, err := h.PersonalCollegeApplicationRepository.ListPersonalCollegeApplicationsByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	return applications, nil
}
