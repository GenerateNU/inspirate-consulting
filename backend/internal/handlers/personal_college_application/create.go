package personalcollegeapplication

import (
	"context"
	"log/slog"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

// CreatePersonalCollegeApplication creates a new personal college application entry in the database
func (h *Handler) CreatePersonalCollegeApplication(ctx context.Context, input models.CreatePersonalCollegeApplicationRequestBody) (*models.PersonalCollegeApplication, error) {
	// validate the college offers the requested cycle
	college, err := h.GlobalCollegeRepository.GetGlobalCollege(ctx, input.GlobalCollegeID)
	if err != nil {
		return nil, err
	}
	err = ValidateApplicationCycle(college, input.ApplicationType)
	if err != nil {
		return nil, err
	}

	studentID, err := auth.GetStudentID(ctx)
	if err != nil {
		return nil, err
	}

	createdPersonalCollegeApplication, err := h.PersonalCollegeApplicationRepository.CreatePersonalCollegeApplication(ctx, studentID, input)
	if err != nil {
		return nil, err
	}

	if err := h.ActivityRepository.CreateActivity(ctx, studentID, models.CreateActivityBody{
		Action:          models.ActionAddedSchool,
		GlobalCollegeID: &createdPersonalCollegeApplication.GlobalCollegeID,
	}); err != nil {
		slog.Error("failed to log added_school activity", "student_id", studentID, "global_college_id", createdPersonalCollegeApplication.GlobalCollegeID, "error", err)
	}

	return createdPersonalCollegeApplication, nil
}
