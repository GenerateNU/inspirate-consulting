package personalcollegeapplication

import (
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func ValidateApplicationDeadline(college *models.GlobalCollege, deadline string) error {

	// validate requested deadline exists
	switch deadline {
	case "ED":
		if college.EDDeadline == nil {
			return errs.BadRequest("this college does not offer ED")
		}
	case "EA":
		if college.EADeadline == nil {
			return errs.BadRequest("this college does not offer EA")
		}
	case "RD":
		if college.RDDeadline == nil {
			return errs.BadRequest("this college does not offer RD")
		}
	default:
		return errs.BadRequest("invalid application_type")
	}
	return nil
}
