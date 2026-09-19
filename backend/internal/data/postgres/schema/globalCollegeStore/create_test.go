package globalCollegeRepository

import (
	"context"
	"testing"
	"time"

	"inspirate-consulting/internal/models"
)

// test normal global college creation with valid input
func TestCreateGlobalCollege(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	edDeadline := time.Now().Add(30 * 24 * time.Hour).UTC().Truncate(time.Second)

	input := models.CreateGlobalCollegeInput{
		SchoolName: "Test University",
		SchoolLocation: "Test City, TS",
		EADeadline: nil,
		EDDeadline: &edDeadline,
		RDDeadline: nil,
	}

	created, err := repo.CreateGlobalCollege(ctx, input)
	if err != nil {
		t.Fatalf("CreateGlobalCollege failed: %v", err)
	}
	cleanupCollege(t, db, created.ID)

	if created.ID == 0 {
		t.Error("expected a generated ID, got 0")
	}
	if created.SchoolName != input.SchoolName {
		t.Errorf("expected SchoolName %q, got %q", input.SchoolName, created.SchoolName)
	}
	if created.SchoolLocation != input.SchoolLocation {
		t.Errorf("expected SchoolLocation %q, got %q", input.SchoolLocation, created.SchoolLocation)
	}
	if created.EADeadline != nil {
		t.Errorf("expected nil EADeadline, got %v", created.EADeadline)
	}
	if created.RDDeadline != nil {
		t.Errorf("expected nil RDDeadline, got %v", created.RDDeadline)
	}
	if created.EDDeadline == nil {
		t.Fatal("expected non-nil EDDeadline")
	}
	if !created.EDDeadline.Equal(edDeadline) {
		t.Errorf("expected EDDeadline %v, got %v", edDeadline, *created.EDDeadline)
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if created.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestCreateGlobalCollege_AllDeadlinesNull(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	input := models.CreateGlobalCollegeInput{
		SchoolName: "No Deadlines University",
		SchoolLocation: "Nowhere, NA",
	}

	created, err := repo.CreateGlobalCollege(ctx, input)
	if err != nil {
		t.Fatalf("expected create with all-nil deadlines to succeed, got error: %v", err)
	}
	cleanupCollege(t, db, created.ID)

	if created.EADeadline != nil || created.EDDeadline != nil || created.RDDeadline != nil {
		t.Errorf("expected all deadlines nil, got EA=%v ED=%v RD=%v",
			created.EADeadline, created.EDDeadline, created.RDDeadline)
	}
}

// test duplication handling for global college creation 
// (case-insensitive unique index on (school_name, school_location))
func TestCreateGlobalCollege_Duplicate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	input := models.CreateGlobalCollegeInput{
		SchoolName: "Duplicate University",
		SchoolLocation: "Dupe City, DC",
	}

	first, err := repo.CreateGlobalCollege(ctx, input)
	if err != nil {
		t.Fatalf("first CreateGlobalCollege failed: %v", err)
	}
	cleanupCollege(t, db, first.ID)

	// Same name/location, different casing, should still collide.
	dupInput := models.CreateGlobalCollegeInput{
		SchoolName: "duplicate university",
		SchoolLocation: "DUPE CITY, DC",
	}
	_, err = repo.CreateGlobalCollege(ctx, dupInput)
	if err == nil {
		t.Fatal("expected an error creating a case-insensitive duplicate, got nil")
	}
}

