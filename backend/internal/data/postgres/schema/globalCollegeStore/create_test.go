package globalCollegeRepository

import (
	"context"
	"testing"
	"time"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// test normal global college creation with valid input
func TestCreateGlobalCollege(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	edDeadline := time.Now().Add(30 * 24 * time.Hour).UTC().Truncate(time.Second)

	input := models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Test University",
		SchoolLocation: "Test City, TS",
		Cycles: []models.CollegeApplicationCycle{
			{Cycle: models.ApplicationCycleED, Deadline: &edDeadline},
		},
	}

	output, err := repo.CreateGlobalCollege(ctx, input)
	if err != nil {
		t.Fatalf("CreateGlobalCollege failed: %v", err)
	}
	created := output

	if created.ID == 0 {
		t.Error("expected a generated ID, got 0")
	}
	if created.SchoolName != input.SchoolName {
		t.Errorf("expected SchoolName %q, got %q", input.SchoolName, created.SchoolName)
	}
	if created.SchoolLocation != input.SchoolLocation {
		t.Errorf("expected SchoolLocation %q, got %q", input.SchoolLocation, created.SchoolLocation)
	}
	if len(created.Cycles) != 1 {
		t.Fatalf("expected 1 cycle, got %d", len(created.Cycles))
	}
	if created.Cycles[0].Cycle != models.ApplicationCycleED {
		t.Errorf("expected cycle ED, got %q", created.Cycles[0].Cycle)
	}
	if created.Cycles[0].Deadline == nil || !created.Cycles[0].Deadline.Equal(edDeadline) {
		t.Errorf("expected ED deadline %v, got %v", edDeadline, created.Cycles[0].Deadline)
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if created.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestCreateGlobalCollege_NoCycles(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	input := models.CreateGlobalCollegeRequestBody{
		SchoolName:     "No Deadlines University",
		SchoolLocation: "Nowhere, NA",
	}

	output, err := repo.CreateGlobalCollege(ctx, input)
	if err != nil {
		t.Fatalf("expected create with no cycles to succeed, got error: %v", err)
	}
	created := output

	if created.Cycles == nil || len(created.Cycles) != 0 {
		t.Errorf("expected an empty (non-nil) cycles list, got %v", created.Cycles)
	}
}

// test duplicate global college creation with same school_name and school_location fails
func TestCreateGlobalCollege_Duplicate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	input := models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Duplicate University",
		SchoolLocation: "Dupe City, DC",
	}

	_, err := repo.CreateGlobalCollege(ctx, input)
	if err != nil {
		t.Fatalf("first CreateGlobalCollege failed: %v", err)
	}

	// Same name/location, different casing, should still collide.
	dupInput := models.CreateGlobalCollegeRequestBody{
		SchoolName:     "duplicate university",
		SchoolLocation: "DUPE CITY, DC",
	}
	_, err = repo.CreateGlobalCollege(ctx, dupInput)
	if err == nil {
		t.Fatal("expected an error creating a case-insensitive duplicate, got nil")
	}
}
