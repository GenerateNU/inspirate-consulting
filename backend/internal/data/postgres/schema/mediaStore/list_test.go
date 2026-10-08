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
		LengthInMins: 10,
		S3Key:        "media/alpha.mp4",
	})
	if err != nil {
		t.Fatalf("setup CreateMedia failed: %v", err)
	}

	mediaB, err := repo.CreateMedia(ctx, &models.CreateMediaRequestBody{
		Title:        "List Test Beta",
		Description:  "Second video",
		LengthInMins: 15,
		S3Key:        "media/beta.mp4",
	})
	if err != nil {
		t.Fatalf("setup CreateMedia failed: %v", err)
	}

	results, err := repo.ListAllMedia(ctx, 10, 0)
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
			LengthInMins: 5,
			S3Key:        "media/page.mp4",
		}); err != nil {
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
			results, err := repo.ListAllMedia(ctx, tc.limit, tc.offset)
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

	results, err := repo.ListAllMedia(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListAllMedia failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 media, got %d", len(results))
	}
}
