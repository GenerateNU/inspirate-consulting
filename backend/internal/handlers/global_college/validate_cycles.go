package globalcollege

import (
	"fmt"
	"time"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// validateCycles checks the shape of each offered cycle before it reaches the database:
// no cycle listed twice, every non-rolling cycle has a deadline and no rolling dates,
// rolling has no deadline and a start that isn't after its end,
// and deadlines are in order (see validateDeadlineOrder).
func validateCycles(cycles []models.CollegeApplicationCycle) error {
	seen := make(map[string]bool, len(cycles))
	deadlines := make(map[string]time.Time, len(cycles))

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
			deadlines[cycle.Cycle] = *cycle.Deadline
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
	return validateDeadlineOrder(deadlines)
}

// validateDeadlineOrder checks ED < ED2 <= RD and ED/EA/REA < RD, skipping cycles that aren't offered.
// ED and EA/REA aren't compared, since schools often give them the same date.
func validateDeadlineOrder(deadlines map[string]time.Time) error {
	ed, hasED := deadlines[models.ApplicationCycleED]
	ed2, hasED2 := deadlines[models.ApplicationCycleED2]
	rd, hasRD := deadlines[models.ApplicationCycleRD]

	if hasED && hasED2 && !ed.Before(ed2) {
		return errs.BadRequest("ED deadline must be before ED2 deadline")
	}
	if hasED2 && hasRD && ed2.After(rd) {
		return errs.BadRequest("ED2 deadline must be on or before RD deadline")
	}
	if !hasRD {
		return nil
	}
	for _, early := range []string{models.ApplicationCycleED, models.ApplicationCycleEA, models.ApplicationCycleREA} {
		if deadline, ok := deadlines[early]; ok && !deadline.Before(rd) {
			return errs.BadRequest(fmt.Sprintf("%s deadline must be before RD deadline", early))
		}
	}
	return nil
}
