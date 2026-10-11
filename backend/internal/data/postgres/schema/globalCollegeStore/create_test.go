package globalCollegeRepository

import (
	"context"
	"testing"
	"time"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

const nearDuplicateThreshold = 0.45

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
		EADeadline:     nil,
		EDDeadline:     &edDeadline,
		RDDeadline:     nil,
	}

	output, err := repo.CreateGlobalCollege(ctx, input, nearDuplicateThreshold)
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

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	input := models.CreateGlobalCollegeRequestBody{
		SchoolName:     "No Deadlines University",
		SchoolLocation: "Nowhere, NA",
	}

	output, err := repo.CreateGlobalCollege(ctx, input, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("expected create with all-nil deadlines to succeed, got error: %v", err)
	}
	created := output

	if created.EADeadline != nil || created.EDDeadline != nil || created.RDDeadline != nil {
		t.Errorf("expected all deadlines nil, got EA=%v ED=%v RD=%v",
			created.EADeadline, created.EDDeadline, created.RDDeadline)
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

	_, err := repo.CreateGlobalCollege(ctx, input, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("first CreateGlobalCollege failed: %v", err)
	}

	// Same name/location, different casing, should still collide.
	dupInput := models.CreateGlobalCollegeRequestBody{
		SchoolName:     "duplicate university",
		SchoolLocation: "DUPE CITY, DC",
	}
	_, err = repo.CreateGlobalCollege(ctx, dupInput, nearDuplicateThreshold)
	if err == nil {
		t.Fatal("expected an error creating a case-insensitive duplicate, got nil")
	}
}

// test a near-duplicate (not exact) is rejected by the fuzzy check
func TestCreateGlobalCollege_RejectsNearDuplicate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	_, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Northeastern University - Boston",
		SchoolLocation: "Boston, MA",
	}, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("first CreateGlobalCollege failed: %v", err)
	}

	_, err = repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Northeastern Boston Campus",
		SchoolLocation: "Boston, MA",
	}, nearDuplicateThreshold)
	if err == nil {
		t.Fatal("expected an error creating a near-duplicate college, got nil")
	}
}

// test that ConfirmDuplicate allows creation despite a detected near-duplicate
func TestCreateGlobalCollege_ConfirmDuplicateOverridesNearDuplicateCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	_, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Northeastern University - Boston",
		SchoolLocation: "Boston, MA",
	}, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("first CreateGlobalCollege failed: %v", err)
	}

	created, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:       "Northeastern Boston Campus",
		SchoolLocation:   "Boston, MA",
		ConfirmDuplicate: true,
	}, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("expected ConfirmDuplicate to allow creation despite the near-duplicate warning, got error: %v", err)
	}
	if created.SchoolName != "Northeastern Boston Campus" {
		t.Errorf("expected SchoolName 'Northeastern Boston Campus', got %q", created.SchoolName)
	}
}

// test that a similarly-named college in a genuinely DIFFERENT location is NOT flagged
func TestCreateGlobalCollege_SimilarNameDifferentLocationNotFlagged(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	_, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Western State University",
		SchoolLocation: "Colorado City, CO",
	}, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("first CreateGlobalCollege failed: %v", err)
	}

	// should be allowed because the location is different, even though the name is similar
	created, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Western State University",
		SchoolLocation: "Somewhere Else, SE",
	}, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("expected no near-duplicate rejection for a different location, got error: %v", err)
	}
	if created.SchoolLocation != "Somewhere Else, SE" {
		t.Errorf("expected SchoolLocation 'Somewhere Else, SE', got %q", created.SchoolLocation)
	}
}
