package globalCollegeRepository

import (
	"context"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// similar name and location should be flagged as a near-duplicate
func TestFindNearDuplicateGlobalColleges_MatchesSimilarNameAndLocation(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	_, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Northeastern University - Boston",
		SchoolLocation: "Boston, MA",
	}, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("setup create failed: %v", err)
	}

	results, err := repo.FindNearDuplicateGlobalColleges(ctx, "Northeastern Boston Campus", "Boston, MA", nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("FindNearDuplicateGlobalColleges failed: %v", err)
	}

	if len(results) == 0 {
		t.Error("expected the known near-duplicate pair to match when both name AND location are similar")
	}
}

// similar-ish name and location should not be flagged based on set threshold
func TestFindNearDuplicateGlobalColleges_DoesNotMatchDifferentNameAndLocation(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	_, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Boston College",
		SchoolLocation: "Boston, MA",
	}, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("setup create failed: %v", err)
	}

	results, err := repo.FindNearDuplicateGlobalColleges(ctx, "Boston University", "Boston, MA", nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("FindNearDuplicateGlobalColleges failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected NO match for similar name but different location (AND logic), but got %d results", len(results))
	}
}

// same name different location should NOT be flagged as a near-duplicate
func TestFindNearDuplicateGlobalColleges_DoesNotMatchSimilarNameDifferentLocation(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	_, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Western State University",
		SchoolLocation: "Colorado City, CO",
	}, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("setup create failed: %v", err)
	}

	results, err := repo.FindNearDuplicateGlobalColleges(ctx, "Western State University", "Faraway City, FC", nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("FindNearDuplicateGlobalColleges failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected NO match for similar name but different location (AND logic), but got %d results", len(results))
	}
}

// same location different name = not a duplicate
func TestFindNearDuplicateGlobalColleges_DoesNotMatchSimilarLocationDifferentName(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	_, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Totally Unrelated College",
		SchoolLocation: "Shared City, SC",
	}, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("setup create failed: %v", err)
	}

	results, err := repo.FindNearDuplicateGlobalColleges(ctx, "Completely Different School", "Shared City, SC", nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("FindNearDuplicateGlobalColleges failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected NO match for similar location but unrelated name (AND logic), but got %d results", len(results))
	}
}
