package personalcollegeapplication

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// CreatePersonalCollegeApplication creates a new personal college application entry in the database
func (h *Handler) CreatePersonalCollegeApplication(ctx context.Context, input models.CreatePersonalCollegeApplicationRequestBody) (*models.PersonalCollegeApplication, error) {
	
	// validate requested deadline exists
	college, err := h.GlobalCollegeRepository.GetGlobalCollege(ctx, input.GlobalCollegeID)
	if err != nil {
		return nil, err
	}
	switch input.ApplicationType {
	case "ED":
		if college.EDDeadline == nil {
			return nil, errs.BadRequest("this college does not offer ED")
		}
	case "EA":
		if college.EADeadline == nil {
			return nil, errs.BadRequest("this college does not offer EA")
		}
	case "RD":
		if college.RDDeadline == nil {
			return nil, errs.BadRequest("this college does not offer RD")
		}
	default:
		return nil, errs.BadRequest("invalid application_type")
	}

	studentID := auth.GetStudentID(ctx)

	createdPersonalCollegeApplication, err := h.PersonalCollegeApplicationRepository.CreatePersonalCollegeApplication(ctx, studentID, input)
	if err != nil {
		return nil, err
	}
 
	return createdPersonalCollegeApplication, nil
}
