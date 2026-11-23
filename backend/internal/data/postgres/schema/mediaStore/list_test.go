package mediaRepository

import (
	"context"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// test list function returns all media when the page is large enough to hold them
func TestListAllMedia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	mediaA, err := repo.CreateMedia(ctx, &models.CreateMediaRequestBody{
		Title:        "List Test Alpha",
		Description:  "First video",
		LengthInMins: ptr(10),
		S3Key:        ptr("media/alpha.mp4"),
	}, nil)
	if err != nil {
		t.Fatalf("setup CreateMedia failed: %v", err)
	}

	mediaB, err := repo.CreateMedia(ctx, &models.CreateMediaRequestBody{
		Title:        "List Test Beta",
		Description:  "Second video",
		LengthInMins: ptr(15),
		S3Key:        ptr("media/beta.mp4"),
	}, nil)
	if err != nil {
		t.Fatalf("setup CreateMedia failed: %v", err)
	}

	results, err := repo.ListAllMedia(ctx, "", 0.2, 10, 0)
	if err != nil {
		t.Fatalf("ListAllMedia failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 media, got %d", len(results))
	}

	foundIDs := map[string]bool{}
	for _, r := range results {
		foundIDs[r.ID] = true
	}
	if !foundIDs[mediaA.ID] || !foundIDs[mediaB.ID] {
		t.Error("expected both created media to be present in results")
	}
}

// test list function respects limit and offset, returning results ordered by title
func TestListAllMedia_Pagination(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	titles := []string{"Page Test A", "Page Test B", "Page Test C"}
	for _, title := range titles {
		if _, err := repo.CreateMedia(ctx, &models.CreateMediaRequestBody{
			Title:        title,
			Description:  "Pagination video",
			LengthInMins: ptr(5),
			S3Key:        ptr("media/page.mp4"),
		}, nil); err != nil {
			t.Fatalf("setup CreateMedia failed: %v", err)
		}
	}

	tests := []struct {
		name           string
		limit          int
		offset         int
		expectedTitles []string
	}{
		{name: "first page", limit: 2, offset: 0, expectedTitles: []string{"Page Test A", "Page Test B"}},
		{name: "second page", limit: 2, offset: 2, expectedTitles: []string{"Page Test C"}},
		{name: "offset past end", limit: 2, offset: 3, expectedTitles: []string{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			results, err := repo.ListAllMedia(ctx, "", 0.2, tc.limit, tc.offset)
			if err != nil {
				t.Fatalf("ListAllMedia failed: %v", err)
			}

			if len(results) != len(tc.expectedTitles) {
				t.Fatalf("expected %d media, got %d", len(tc.expectedTitles), len(results))
			}
			for i, r := range results {
				if r.Title != tc.expectedTitles[i] {
					t.Errorf("result %d: expected title %q, got %q", i, tc.expectedTitles[i], r.Title)
				}
			}
		})
	}
}

// test list function returns an empty result, not an error, when there is no media
func TestListAllMedia_NoMedia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	results, err := repo.ListAllMedia(ctx, "", 0.2, 10, 0)
	if err != nil {
		t.Fatalf("ListAllMedia failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 media, got %d", len(results))
	}
}

// test the search query matches a word in the title or description, case-insensitively, even in a
// long description, and leaves out media that doesn't match
func TestListAllMedia_Search(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	byTitle, err := repo.CreateMedia(ctx, &models.CreateMediaRequestBody{
		Title:       "FAFSA Basics",
		Description: "Financial aid overview",
		S3Key:       ptr("media/fafsa.mp4"),
	}, nil)
	if err != nil {
		t.Fatalf("setup CreateMedia failed: %v", err)
	}
	byDescription, err := repo.CreateMedia(ctx, &models.CreateMediaRequestBody{
		Title:       "Money Matters",
		Description: "Filling out the fafsa step by step, with screenshots of every page of the form",
		S3Key:       ptr("media/money.mp4"),
	}, nil)
	if err != nil {
		t.Fatalf("setup CreateMedia failed: %v", err)
	}
	if _, err := repo.CreateMedia(ctx, &models.CreateMediaRequestBody{
		Title:       "Essay Tips",
		Description: "Unrelated",
		S3Key:       ptr("media/essay.mp4"),
	}, nil); err != nil {
		t.Fatalf("setup CreateMedia failed: %v", err)
	}

	results, err := repo.ListAllMedia(ctx, "FaFsA", 0.4, 10, 0)
	if err != nil {
		t.Fatalf("ListAllMedia failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 matching media, got %d", len(results))
	}
	foundIDs := map[string]bool{results[0].ID: true, results[1].ID: true}
	if !foundIDs[byTitle.ID] || !foundIDs[byDescription.ID] {
		t.Error("expected matches on both title and description")
	}
}

// test a misspelled query still matches, and an exact word match ranks above a close one even when
// it sorts later alphabetically
func TestListAllMedia_SearchTypoAndRanking(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	// "A Fafsy Overview" sorts first alphabetically, but "FAFSA Guide" contains the exact word
	closeMatch, err := repo.CreateMedia(ctx, &models.CreateMediaRequestBody{
		Title:       "A Fafsy Overview",
		Description: "test",
		S3Key:       ptr("media/overview.mp4"),
	}, nil)
	if err != nil {
		t.Fatalf("setup CreateMedia failed: %v", err)
	}
	exact, err := repo.CreateMedia(ctx, &models.CreateMediaRequestBody{
		Title:       "FAFSA Guide",
		Description: "test",
		S3Key:       ptr("media/exact.mp4"),
	}, nil)
	if err != nil {
		t.Fatalf("setup CreateMedia failed: %v", err)
	}

	results, err := repo.ListAllMedia(ctx, "fafsa", 0.4, 10, 0)
	if err != nil {
		t.Fatalf("ListAllMedia failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 matching media, got %d", len(results))
	}
	if results[0].ID != exact.ID || results[1].ID != closeMatch.ID {
		t.Errorf("expected %q before %q, got %q then %q", exact.Title, closeMatch.Title, results[0].Title, results[1].Title)
	}
}
