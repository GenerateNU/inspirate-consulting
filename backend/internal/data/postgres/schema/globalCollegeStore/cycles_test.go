package globalCollegeRepository

import (
	"context"
	"testing"
	"time"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

func strPtr(s string) *string { return &s }

// updating a college replaces its full set of cycles: new ones are added, changed ones updated, missing ones removed
func TestUpdateGlobalCollege_ReplacesCycles(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	edDeadline := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	eaDeadline := time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC)
	newEADeadline := time.Date(2026, 11, 20, 0, 0, 0, 0, time.UTC)

	created, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Cycle University",
		SchoolLocation: "Cycle City, CC",
		Cycles: []models.CollegeApplicationCycle{
			{Cycle: models.ApplicationCycleEA, Deadline: &eaDeadline},
			{Cycle: models.ApplicationCycleED, Deadline: &edDeadline},
		},
	})
	if err != nil {
		t.Fatalf("setup CreateGlobalCollege failed: %v", err)
	}
	if len(created.Cycles) != 2 || created.Cycles[0].Cycle != models.ApplicationCycleED {
		t.Fatalf("expected cycles ordered ED, EA; got %+v", created.Cycles)
	}

	updated, err := repo.UpdateGlobalCollege(ctx, created.ID, models.UpdateGlobalCollegeRequestBody{
		SchoolName:     "Cycle University",
		SchoolLocation: "Cycle City, CC",
		Cycles: []models.CollegeApplicationCycle{
			{Cycle: models.ApplicationCycleRolling, RollingStart: strPtr("2026-09-01"), RollingEnd: strPtr("2027-03-01")},
			{Cycle: models.ApplicationCycleEA, Deadline: &newEADeadline},
		},
	})
	if err != nil {
		t.Fatalf("UpdateGlobalCollege failed: %v", err)
	}

	if len(updated.Cycles) != 2 {
		t.Fatalf("expected 2 cycles after update (ED removed), got %+v", updated.Cycles)
	}
	ea, rolling := updated.Cycles[0], updated.Cycles[1]
	if ea.Cycle != models.ApplicationCycleEA || ea.Deadline == nil || !ea.Deadline.Equal(newEADeadline) {
		t.Errorf("expected EA deadline updated to %v, got %+v", newEADeadline, ea)
	}
	if rolling.Cycle != models.ApplicationCycleRolling || rolling.Deadline != nil {
		t.Errorf("expected ROLLING with no deadline, got %+v", rolling)
	}
	if rolling.RollingStart == nil || *rolling.RollingStart != "2026-09-01" ||
		rolling.RollingEnd == nil || *rolling.RollingEnd != "2027-03-01" {
		t.Errorf("expected rolling range 2026-09-01..2027-03-01, got %v..%v", rolling.RollingStart, rolling.RollingEnd)
	}

	// get returns the same cycles
	got, err := repo.GetGlobalCollege(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetGlobalCollege failed: %v", err)
	}
	if len(got.Cycles) != 2 {
		t.Errorf("expected GetGlobalCollege to return 2 cycles, got %+v", got.Cycles)
	}

	// an empty list removes every cycle
	cleared, err := repo.UpdateGlobalCollege(ctx, created.ID, models.UpdateGlobalCollegeRequestBody{
		SchoolName:     "Cycle University",
		SchoolLocation: "Cycle City, CC",
	})
	if err != nil {
		t.Fatalf("UpdateGlobalCollege (clear) failed: %v", err)
	}
	if len(cleared.Cycles) != 0 {
		t.Errorf("expected no cycles after clearing, got %+v", cleared.Cycles)
	}
}

// a bad cycle rolls back the whole write, including the name change
func TestUpdateGlobalCollege_InvalidCycleRollsBack(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	created, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Rollback University",
		SchoolLocation: "Rollback City, RC",
	})
	if err != nil {
		t.Fatalf("setup CreateGlobalCollege failed: %v", err)
	}

	// ED without a deadline violates college_application_cycles_deadline_shape
	_, err = repo.UpdateGlobalCollege(ctx, created.ID, models.UpdateGlobalCollegeRequestBody{
		SchoolName:     "Renamed University",
		SchoolLocation: "Rollback City, RC",
		Cycles:         []models.CollegeApplicationCycle{{Cycle: models.ApplicationCycleED}},
	})
	if err == nil {
		t.Fatal("expected an error for ED without a deadline, got nil")
	}

	got, err := repo.GetGlobalCollege(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetGlobalCollege failed: %v", err)
	}
	if got.SchoolName != "Rollback University" {
		t.Errorf("expected name change to roll back, got %q", got.SchoolName)
	}
}
