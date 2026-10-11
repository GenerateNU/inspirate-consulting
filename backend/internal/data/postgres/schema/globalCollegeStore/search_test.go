package globalCollegeRepository

import (
	"context"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

const searchThreshold = 0.2

// query will find by college name
func TestSearchGlobalColleges_MatchesByName(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	_, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Searchable Northeastern University",
		SchoolLocation: "Boston, MA",
	}, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("setup create failed: %v", err)
	}

	results, err := repo.SearchGlobalColleges(ctx, "Northeastern", searchThreshold)
	if err != nil {
		t.Fatalf("SearchGlobalColleges failed: %v", err)
	}

	found := false
	for _, r := range results {
		if r.SchoolName == "Searchable Northeastern University" {
			found = true
		}
	}
	if !found {
		t.Error("expected to find the college by partial name match")
	}
}

// query will find by college location
func TestSearchGlobalColleges_MatchesByLocation(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	_, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Some Unrelated College Name",
		SchoolLocation: "Uniquesburgh, UB",
	}, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("setup create failed: %v", err)
	}

	results, err := repo.SearchGlobalColleges(ctx, "Uniquesburgh", searchThreshold)
	if err != nil {
		t.Fatalf("SearchGlobalColleges failed: %v", err)
	}

	found := false
	for _, r := range results {
		if r.SchoolLocation == "Uniquesburgh, UB" {
			found = true
		}
	}
	if !found {
		t.Error("expected to find the college by location match, confirming OR logic covers both fields")
	}
}

// neither name nor location should match this query
func TestSearchGlobalColleges_NoMatch(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	results, err := repo.SearchGlobalColleges(ctx, "Completely Unrelated Query", searchThreshold)
	if err != nil {
		t.Fatalf("SearchGlobalColleges failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no matches for an unrelated query, got %d", len(results))
	}
}
