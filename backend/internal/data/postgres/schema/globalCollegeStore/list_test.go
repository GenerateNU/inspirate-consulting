package globalCollegeRepository

import (
	"context"
	"testing"
	"time"

	"inspirate-consulting/internal/models"
)

// test list function returns all global colleges
// intentionally is not run in parallel, because it depends on the state of the database and could be affected by other tests
// assumes initial database state is empty
func TestListGlobalColleges(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	before, err := repo.ListGlobalColleges(ctx)
	if err != nil {
		t.Fatalf("ListGlobalColleges failed: %v", err)
	}
	beforeCount := len(before)

	eaDeadline := time.Now().Add(10 * 24 * time.Hour).UTC().Truncate(time.Second)
	edDeadline := time.Now().Add(20 * 24 * time.Hour).UTC().Truncate(time.Second)
	rdDeadline := time.Now().Add(90 * 24 * time.Hour).UTC().Truncate(time.Second)

	inputs := []models.CreateGlobalCollegeInput{
		{
			SchoolName:     "Alpha University",
			SchoolLocation: "Alpha City, AA",
			EADeadline:     &eaDeadline,
		},
		{
			SchoolName:     "Beta College",
			SchoolLocation: "Beta Town, BB",
			EDDeadline:     &edDeadline,
		},
		{
			SchoolName:     "Gamma Institute",
			SchoolLocation: "Gamma Village, CC",
			RDDeadline:     &rdDeadline,
		},
		{
			SchoolName:     "Delta State",
			SchoolLocation: "Delta City, DD",
		},
	}

	created := make([]*models.GlobalCollege, 0, len(inputs))
	for _, input := range inputs {
		c, err := repo.CreateGlobalCollege(ctx, input)
		if err != nil {
			t.Fatalf("setup CreateGlobalCollege(%q) failed: %v", input.SchoolName, err)
		}
		cleanupCollege(t, db, c.ID)
		created = append(created, c)
	}

	after, err := repo.ListGlobalColleges(ctx)
	if err != nil {
		t.Fatalf("ListGlobalColleges failed: %v", err)
	}

	if len(after) != beforeCount+len(inputs) {
		t.Errorf("expected %d colleges after creating %d new ones, got %d",
			beforeCount+len(inputs), len(inputs), len(after))
	}

	// Every created college should be present in the list, with matching fields.
	byID := make(map[int64]*models.GlobalCollege, len(after))
	for i := range after {
		byID[after[i].ID] = after[i]
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
	if alpha.EADeadline == nil || !alpha.EADeadline.Equal(eaDeadline) {
		t.Errorf("expected Alpha University EADeadline %v, got %v", eaDeadline, alpha.EADeadline)
	}
	if alpha.EDDeadline != nil || alpha.RDDeadline != nil {
		t.Errorf("expected Alpha University ED/RD deadlines nil, got ED=%v RD=%v", alpha.EDDeadline, alpha.RDDeadline)
	}

	delta := byID[created[3].ID]
	if delta.EADeadline != nil || delta.EDDeadline != nil || delta.RDDeadline != nil {
		t.Errorf("expected Delta State to have all nil deadlines, got EA=%v ED=%v RD=%v",
			delta.EADeadline, delta.EDDeadline, delta.RDDeadline)
	}
}
