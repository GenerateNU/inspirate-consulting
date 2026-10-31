package globalcollege

import (
	"context"
	"fmt"
	"net/http"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// range from 0.0 to 1.0, inclusive, where 0.0 is a complete mismatch and 1.0 is a perfect match
const nearDuplicateThreshold = 0.45

// returns a 409 if another college is similar to name/location
func (h *Handler) checkNearDuplicate(ctx context.Context, name, location string, excludeID *int64) error {
	similar, err := h.GlobalCollegeRepository.FindNearDuplicateGlobalColleges(ctx, name, location, nearDuplicateThreshold)
	if err != nil {
		return err
	}

	// exclude self
	match := firstOtherCollege(similar, excludeID)
	if match == nil {
		return nil
	}

	return errs.HTTPError{
		Code: http.StatusConflict,
		Message: fmt.Sprintf(
			"%q is very similar to an existing college: %q. Resubmit with confirm_duplicate=true to proceed anyway.",
			name, match.SchoolName,
		),
	}
}

// returns the first college whose ID isn't excludeID, or nil if there isn't one.
func firstOtherCollege(colleges []models.GlobalCollege, excludeID *int64) *models.GlobalCollege {
	for i := range colleges {
		if excludeID == nil || colleges[i].ID != *excludeID {
			return &colleges[i]
		}
	}
	return nil
}
