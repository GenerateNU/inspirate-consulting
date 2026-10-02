package mediaRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// test normal media retrieval with a valid ID
func TestGetMedia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	created, err := repo.CreateMedia(ctx, &models.CreateMediaRequestBody{
		Title:        "Findable Video",
		Description:  "Should be retrievable by id",
		LengthInMins: 20,
		S3Key:        "media/findable.mp4",
	})
	if err != nil {
		t.Fatalf("setup CreateMedia failed: %v", err)
	}

	fetched, err := repo.GetMedia(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetMedia failed: %v", err)
	}

	if fetched.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, fetched.ID)
	}
	if fetched.Title != created.Title {
		t.Errorf("expected Title %q, got %q", created.Title, fetched.Title)
	}
	if fetched.S3Key != created.S3Key {
		t.Errorf("expected S3Key %q, got %q", created.S3Key, fetched.S3Key)
	}
}

// test retrieval of media with a nonexistent ID returns an error
func TestGetMedia_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	_, err := repo.GetMedia(ctx, uuid.NewString())
	if err == nil {
		t.Fatal("expected an error for a nonexistent id, got nil")
	}
}
