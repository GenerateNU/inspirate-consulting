package globalcollege

import (
	"testing"
	"time"

	"inspirate-consulting/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestValidateCycles(t *testing.T) {
	t.Parallel()

	deadline := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	start, end, badDate := "2026-09-01", "2027-03-01", "09/01/2026"

	tests := []struct {
		name    string
		cycles  []models.CollegeApplicationCycle
		wantErr bool
	}{
		{name: "no cycles", cycles: nil},
		{name: "all deadline cycles", cycles: []models.CollegeApplicationCycle{
			{Cycle: "ED", Deadline: &deadline}, {Cycle: "ED2", Deadline: &deadline}, {Cycle: "REA", Deadline: &deadline},
			{Cycle: "EA", Deadline: &deadline}, {Cycle: "RD", Deadline: &deadline},
		}},
		{name: "rolling with range", cycles: []models.CollegeApplicationCycle{{Cycle: "ROLLING", RollingStart: &start, RollingEnd: &end}}},
		{name: "rolling with no dates", cycles: []models.CollegeApplicationCycle{{Cycle: "ROLLING"}}},
		{name: "duplicate cycle", cycles: []models.CollegeApplicationCycle{{Cycle: "ED", Deadline: &deadline}, {Cycle: "ED", Deadline: &deadline}}, wantErr: true},
		{name: "ED without deadline", cycles: []models.CollegeApplicationCycle{{Cycle: "ED"}}, wantErr: true},
		{name: "EA with rolling dates", cycles: []models.CollegeApplicationCycle{{Cycle: "EA", Deadline: &deadline, RollingStart: &start}}, wantErr: true},
		{name: "rolling with deadline", cycles: []models.CollegeApplicationCycle{{Cycle: "ROLLING", Deadline: &deadline}}, wantErr: true},
		{name: "rolling start after end", cycles: []models.CollegeApplicationCycle{{Cycle: "ROLLING", RollingStart: &end, RollingEnd: &start}}, wantErr: true},
		{name: "rolling bad date", cycles: []models.CollegeApplicationCycle{{Cycle: "ROLLING", RollingStart: &badDate, RollingEnd: &end}}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validateCycles(tt.cycles)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
