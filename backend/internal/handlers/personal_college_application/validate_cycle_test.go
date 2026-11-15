package personalcollegeapplication

import (
	"testing"
	"time"

	"inspirate-consulting/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestValidateApplicationCycle(t *testing.T) {
	t.Parallel()

	deadline := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	college := &models.GlobalCollege{
		ID: 1,
		Cycles: []models.CollegeApplicationCycle{
			{Cycle: models.ApplicationCycleED2, Deadline: &deadline},
			{Cycle: models.ApplicationCycleRolling},
		},
	}

	assert.NoError(t, ValidateApplicationCycle(college, models.ApplicationCycleED2))
	assert.NoError(t, ValidateApplicationCycle(college, models.ApplicationCycleRolling))
	assert.EqualError(t, ValidateApplicationCycle(college, models.ApplicationCycleED), "http error: 400 this college does not offer ED")
	assert.EqualError(t, ValidateApplicationCycle(college, "NOT_REAL"), "http error: 400 invalid application_type")
	assert.Error(t, ValidateApplicationCycle(&models.GlobalCollege{ID: 2}, models.ApplicationCycleRD))
}
