package personalcollegeapplication

import (
	"fmt"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// ValidateApplicationCycle checks that the college offers the cycle the student wants to apply by.
func ValidateApplicationCycle(college *models.GlobalCollege, cycle string) error {
	switch cycle {
	case models.ApplicationCycleED, models.ApplicationCycleED2, models.ApplicationCycleREA,
		models.ApplicationCycleEA, models.ApplicationCycleRD, models.ApplicationCycleRolling:
	default:
		return errs.BadRequest("invalid application_type")
	}

	for _, offered := range college.Cycles {
		if offered.Cycle == cycle {
			return nil
		}
	}
	return errs.BadRequest(fmt.Sprintf("this college does not offer %s", cycle))
}
