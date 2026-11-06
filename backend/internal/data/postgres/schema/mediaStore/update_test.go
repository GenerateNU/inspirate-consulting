package mediaRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// test normal media retitling with a valid ID
func TestUpdateMedia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	created, err := repo.CreateMedia(ctx, &models.CreateMediaRequestBody{
		Title:        "Original Title",
		Description:  "Should keep its description after a retitle",
		LengthInMins: 12,
		S3Key:        "media/retitle.mp4",
	})
	if err != nil {
		t.Fatalf("setup CreateMedia failed: %v", err)
	}

	updated, err := repo.UpdateMedia(ctx, created.ID, "Updated Title")
	if err != nil {
		t.Fatalf("UpdateMedia failed: %v", err)
	}

	if updated.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, updated.ID)
	}
	if updated.Title != "Updated Title" {
		t.Errorf("expected Title %q, got %q", "Updated Title", updated.Title)
	}
	// Only the title should change.
	if updated.Description != created.Description {
		t.Errorf("expected Description %q, got %q", created.Description, updated.Description)
	}
	if updated.LengthInMins != created.LengthInMins {
		t.Errorf("expected LengthInMins %d, got %d", created.LengthInMins, updated.LengthInMins)
	}
	if updated.S3Key != created.S3Key {
		t.Errorf("expected S3Key %q, got %q", created.S3Key, updated.S3Key)
	}

	// The new title should be persisted.
	fetched, err := repo.GetMedia(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetMedia failed: %v", err)
	}
	if fetched.Title != "Updated Title" {
		t.Errorf("expected persisted Title %q, got %q", "Updated Title", fetched.Title)
	}
}

// test updating media with a nonexistent ID returns an error
func TestUpdateMedia_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	_, err := repo.UpdateMedia(ctx, uuid.NewString(), "Updated Title")
	if err == nil {
		t.Fatal("expected an error for a nonexistent id, got nil")
	}
}
