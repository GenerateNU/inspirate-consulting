package personalcollegeapplication

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

// updates the rank of a personal college application entry in the database
// handles auto-update of other applications' ranks to maintain a valid ranking order
func (h *Handler) UpdatePersonalCollegeApplicationRank(ctx context.Context, application_id int64, input models.UpdatePersonalCollegeApplicationRankRequestBody) ([]models.PersonalCollegeApplication, error) {

	studentID := auth.GetStudentID(ctx)

	updatedPersonalCollegeApplications, err := h.PersonalCollegeApplicationRepository.UpdateApplicationRank(ctx, studentID, application_id, input.Rank)
	if err != nil {
		return nil, err
	}

	return updatedPersonalCollegeApplications, nil
}
