package mediaRepository

import (
	"context"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// test list function returns all media
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

	results, err := repo.ListAllMedia(ctx)
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

// test list function returns an empty result, not an error, when there is no media
func TestListAllMedia_NoMedia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	results, err := repo.ListAllMedia(ctx)
	if err != nil {
		t.Fatalf("ListAllMedia failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 media, got %d", len(results))
	}
}
