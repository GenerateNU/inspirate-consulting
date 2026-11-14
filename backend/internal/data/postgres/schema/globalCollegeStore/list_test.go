package globalCollegeRepository

import (
	"context"
	"testing"
	"time"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// test list function returns all global colleges
// intentionally is not run in parallel, because it depends on the state of the database and could be affected by other tests
// assumes initial database state is empty
func TestListGlobalColleges(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	beforeOutput, err := repo.ListGlobalColleges(ctx)
	if err != nil {
		t.Fatalf("ListGlobalColleges failed: %v", err)
	}
	beforeCount := len(beforeOutput)

	eaDeadline := time.Now().Add(10 * 24 * time.Hour).UTC().Truncate(time.Second)
	edDeadline := time.Now().Add(20 * 24 * time.Hour).UTC().Truncate(time.Second)
	rdDeadline := time.Now().Add(90 * 24 * time.Hour).UTC().Truncate(time.Second)

	inputs := []models.CreateGlobalCollegeRequestBody{
		{
			SchoolName:     "Alpha University",
			SchoolLocation: "Alpha City, AA",
			Cycles:         []models.CollegeApplicationCycle{{Cycle: models.ApplicationCycleEA, Deadline: &eaDeadline}},
		},
		{
			SchoolName:     "Beta College",
			SchoolLocation: "Beta Town, BB",
			Cycles:         []models.CollegeApplicationCycle{{Cycle: models.ApplicationCycleED, Deadline: &edDeadline}},
		},
		{
			SchoolName:     "Gamma Institute",
			SchoolLocation: "Gamma Village, CC",
			Cycles:         []models.CollegeApplicationCycle{{Cycle: models.ApplicationCycleRD, Deadline: &rdDeadline}},
		},
		{
			SchoolName:     "Delta State",
			SchoolLocation: "Delta City, DD",
		},
	}

	created := make([]models.GlobalCollege, 0, len(inputs))
	for _, input := range inputs {
		output, err := repo.CreateGlobalCollege(ctx, input)
		if err != nil {
			t.Fatalf("setup CreateGlobalCollege(%q) failed: %v", input.SchoolName, err)
		}
		created = append(created, *output)
	}

	afterOutput, err := repo.ListGlobalColleges(ctx)
	if err != nil {
		t.Fatalf("ListGlobalColleges failed: %v", err)
	}
	after := afterOutput

	if len(after) != beforeCount+len(inputs) {
		t.Errorf("expected %d colleges after creating %d new ones, got %d",
			beforeCount+len(inputs), len(inputs), len(after))
	}

	// Every created college should be present in the list, with matching fields.
	byID := make(map[int64]models.GlobalCollege, len(after))
	for _, c := range after {
		byID[c.ID] = c
	}

	for _, want := range created {
		got, ok := byID[want.ID]
		if !ok {
			t.Errorf("expected college %q (id=%d) to appear in ListGlobalColleges result", want.SchoolName, want.ID)
			continue
		}
		if got.SchoolName != want.SchoolName {
			t.Errorf("id=%d: expected SchoolName %q, got %q", want.ID, want.SchoolName, got.SchoolName)
		}
		if got.SchoolLocation != want.SchoolLocation {
			t.Errorf("id=%d: expected SchoolLocation %q, got %q", want.ID, want.SchoolLocation, got.SchoolLocation)
		}
	}

	// Spot-check distinct fields were maintained
	alpha := byID[created[0].ID]
	if len(alpha.Cycles) != 1 || alpha.Cycles[0].Cycle != models.ApplicationCycleEA ||
		alpha.Cycles[0].Deadline == nil || !alpha.Cycles[0].Deadline.Equal(eaDeadline) {
		t.Errorf("expected Alpha University to offer only EA on %v, got %+v", eaDeadline, alpha.Cycles)
	}

	delta := byID[created[3].ID]
	if len(delta.Cycles) != 0 {
		t.Errorf("expected Delta State to offer no cycles, got %+v", delta.Cycles)
	}
}
