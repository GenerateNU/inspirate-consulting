package personalcollegeapplication

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

// updates the global college ID, category, and/or deadline of a personal college application entry in the database
func (h *Handler) UpdatePersonalCollegeApplication(ctx context.Context, id int64, input models.UpdatePersonalCollegeApplicationRequestBody) (*models.PersonalCollegeApplication, error) {

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

	updatedPersonalCollegeApplication, err := h.PersonalCollegeApplicationRepository.UpdatePersonalCollegeApplication(ctx, studentID, id, input)
	if err != nil {
		return nil, err
	}

	return updatedPersonalCollegeApplication, nil
}
