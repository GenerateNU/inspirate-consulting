package mediaRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// test normal media deletion with a valid ID
func TestDeleteMedia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	created, err := repo.CreateMedia(ctx, &models.CreateMediaRequestBody{
		Title:        "Deletable Video",
		Description:  "Should be removed",
		LengthInMins: 5,
		S3Key:        "media/deletable.mp4",
	})
	if err != nil {
		t.Fatalf("setup CreateMedia failed: %v", err)
	}

	if err := repo.DeleteMedia(ctx, created.ID); err != nil {
		t.Fatalf("DeleteMedia failed: %v", err)
	}

	// The deleted media should no longer be listed.
	results, err := repo.ListAllMedia(ctx, 10, 0)
	if err != nil {
		t.Fatalf("ListAllMedia failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 media after delete, got %d", len(results))
	}
}

// test deleting media with a nonexistent ID returns an error
func TestDeleteMedia_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	err := repo.DeleteMedia(ctx, uuid.NewString())
	if err == nil {
		t.Fatal("expected an error for a nonexistent id, got nil")
	}
}
