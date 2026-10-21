package personalcollegeapplication

import (
	"context"

	"inspirate-consulting/internal/auth"
)

// DeletePersonalCollegeApplication deletes a personal college application entry from the database
func (h *Handler) DeletePersonalCollegeApplication(ctx context.Context, applicationID int64) error {
	studentID, err := auth.GetStudentID(ctx)
	if err != nil {
		return err
	}

	return h.PersonalCollegeApplicationRepository.DeletePersonalCollegeApplication(ctx, studentID, applicationID)
}
