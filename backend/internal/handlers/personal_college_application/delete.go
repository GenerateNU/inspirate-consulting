package personalcollegeapplication

import (
	"context"
	"log/slog"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

// DeletePersonalCollegeApplication deletes a personal college application entry from the database
func (h *Handler) DeletePersonalCollegeApplication(ctx context.Context, applicationID int64) error {
	studentID, err := auth.GetStudentID(ctx)
	if err != nil {
		return err
	}

	collegeID, err := h.PersonalCollegeApplicationRepository.DeletePersonalCollegeApplication(ctx, studentID, applicationID)
	if err != nil {
		return err
	}

	if err := h.ActivityRepository.CreateActivity(ctx, studentID, models.CreateActivityBody{
		Action:          models.ActionRemovedSchool,
		GlobalCollegeID: &collegeID,
	}); err != nil {
		slog.Error("failed to log removed_school activity", "student_id", studentID, "global_college_id", collegeID, "error", err)
	}

	return nil
}
