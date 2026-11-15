package globalcollege

import (
	"fmt"
	"time"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// validateCycles checks the shape of each offered cycle before it reaches the database:
// no cycle listed twice, every non-rolling cycle has a deadline and no rolling dates,
// and rolling has no deadline and a start that isn't after its end.
func validateCycles(cycles []models.CollegeApplicationCycle) error {
	seen := make(map[string]bool, len(cycles))

	for _, cycle := range cycles {
		if seen[cycle.Cycle] {
			return errs.BadRequest(fmt.Sprintf("cycle %s is listed more than once", cycle.Cycle))
		}
		seen[cycle.Cycle] = true

		if cycle.Cycle != models.ApplicationCycleRolling {
			if cycle.Deadline == nil {
				return errs.BadRequest(fmt.Sprintf("cycle %s needs a deadline", cycle.Cycle))
			}
			if cycle.RollingStart != nil || cycle.RollingEnd != nil {
				return errs.BadRequest(fmt.Sprintf("only ROLLING can have rolling_start/rolling_end, not %s", cycle.Cycle))
			}
			continue
		}

		if cycle.Deadline != nil {
			return errs.BadRequest("ROLLING has no deadline; use rolling_start/rolling_end instead")
		}
		if cycle.RollingStart != nil && cycle.RollingEnd != nil {
			start, startErr := time.Parse(time.DateOnly, *cycle.RollingStart)
			end, endErr := time.Parse(time.DateOnly, *cycle.RollingEnd)
			if startErr != nil || endErr != nil {
				return errs.BadRequest("rolling_start and rolling_end must be YYYY-MM-DD dates")
			}
			if start.After(end) {
				return errs.BadRequest("rolling_start must be on or before rolling_end")
			}
		}
	}
	return nil
}
