package globalCollegeRepository

import (
	"context"
	"fmt"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

const (
	searchThreshold = 0.2
	testSearchLimit = 10
)

func TestSearchGlobalColleges_MatchesByName(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	_, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Searchable Northeastern University",
		SchoolLocation: "Boston, MA",
	})
	if err != nil {
		t.Fatalf("setup create failed: %v", err)
	}

	results, err := repo.SearchGlobalColleges(ctx, "Northeastern", searchThreshold, testSearchLimit)
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

func TestSearchGlobalColleges_MatchesByLocation(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	_, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Some Unrelated College Name",
		SchoolLocation: "Uniquesburgh, UB",
	})
	if err != nil {
		t.Fatalf("setup create failed: %v", err)
	}

	// Query matches the location, not the name, confirming search checks both.
	results, err := repo.SearchGlobalColleges(ctx, "Uniquesburgh", searchThreshold, testSearchLimit)
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
		t.Error("expected to find the college by location match")
	}
}

func TestSearchGlobalColleges_NoMatch(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	results, err := repo.SearchGlobalColleges(ctx, "Completely Unrelated Query", searchThreshold, testSearchLimit)
	if err != nil {
		t.Fatalf("SearchGlobalColleges failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no matches for an unrelated query, got %d", len(results))
	}
}

// more matching colleges than the limit: only `limit` should come back in relevance order
func TestSearchGlobalColleges_RespectsLimit(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	for i := 1; i <= 5; i++ {
		_, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
			SchoolName:     fmt.Sprintf("Harborview College Campus %d", i),
			SchoolLocation: fmt.Sprintf("City %d, ST", i),
		})
		if err != nil {
			t.Fatalf("setup create %d failed: %v", i, err)
		}
	}
	_, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Harborview College",
		SchoolLocation: "Exact City, EC",
	})
	if err != nil {
		t.Fatalf("setup create exact match failed: %v", err)
	}

	const limit = 3
	results, err := repo.SearchGlobalColleges(ctx, "Harborview College", searchThreshold, limit)
	if err != nil {
		t.Fatalf("SearchGlobalColleges failed: %v", err)
	}

	if len(results) != limit {
		t.Fatalf("expected exactly %d results with 6 matching colleges, got %d", limit, len(results))
	}
	if results[0].SchoolName != "Harborview College" {
		t.Errorf("expected the exact match first, got %q", results[0].SchoolName)
	}
}
