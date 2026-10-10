package personalcollegeapplication

import (
	"context"

	"inspirate-consulting/internal/auth"
)

// DeletePersonalCollegeApplication deletes a personal college application entry from the database
func (h *Handler) DeletePersonalCollegeApplication(ctx context.Context, applicationID int64) error {
	studentID := auth.GetStudentID(ctx)

	return h.PersonalCollegeApplicationRepository.DeletePersonalCollegeApplication(ctx, studentID, applicationID)
}
