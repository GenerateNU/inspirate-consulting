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
		LengthInMins: ptr(12),
		S3Key:        ptr("media/retitle.mp4"),
	}, nil)
	if err != nil {
		t.Fatalf("setup CreateMedia failed: %v", err)
	}

	updated, err := repo.UpdateMedia(ctx, created.ID, models.UpdateMediaFields{Title: ptr("Updated Title")})
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
	if updated.LengthInMins == nil || *updated.LengthInMins != *created.LengthInMins {
		t.Errorf("expected LengthInMins %d, got %v", *created.LengthInMins, updated.LengthInMins)
	}
	if updated.S3Key == nil || *updated.S3Key != *created.S3Key {
		t.Errorf("expected S3Key %q, got %v", *created.S3Key, updated.S3Key)
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

// test that setting a YouTube id on S3-backed media swaps the source and keeps the title
func TestUpdateMedia_SwitchToYoutube(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	created, err := repo.CreateMedia(ctx, &models.CreateMediaRequestBody{
		Title:        "Switching Source",
		Description:  "Starts on S3",
		LengthInMins: ptr(8),
		S3Key:        ptr("media/switch.mp4"),
	}, nil)
	if err != nil {
		t.Fatalf("setup CreateMedia failed: %v", err)
	}

	updated, err := repo.UpdateMedia(ctx, created.ID, models.UpdateMediaFields{YoutubeVideoID: ptr("dQw4w9WgXcQ")})
	if err != nil {
		t.Fatalf("UpdateMedia failed: %v", err)
	}

	if updated.Title != created.Title {
		t.Errorf("expected Title %q to be unchanged, got %q", created.Title, updated.Title)
	}
	if updated.YoutubeVideoID == nil || *updated.YoutubeVideoID != "dQw4w9WgXcQ" {
		t.Errorf("expected YoutubeVideoID %q, got %v", "dQw4w9WgXcQ", updated.YoutubeVideoID)
	}
	if updated.S3Key != nil {
		t.Errorf("expected S3Key to be cleared, got %q", *updated.S3Key)
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

	_, err := repo.UpdateMedia(ctx, uuid.NewString(), models.UpdateMediaFields{Title: ptr("Updated Title")})
	if err == nil {
		t.Fatal("expected an error for a nonexistent id, got nil")
	}
}

// test updating the description and switching a YouTube video back to S3
func TestUpdateMedia_DescriptionAndSwitchToS3(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	created, err := repo.CreateMedia(ctx, &models.CreateMediaRequestBody{
		Title:       "Starts on YouTube",
		Description: "Old description",
	}, ptr("dQw4w9WgXcQ"))
	if err != nil {
		t.Fatalf("setup CreateMedia failed: %v", err)
	}

	updated, err := repo.UpdateMedia(ctx, created.ID, models.UpdateMediaFields{
		Description: ptr("New description"),
		S3Key:       ptr("media/now-on-s3.mp4"),
	})
	if err != nil {
		t.Fatalf("UpdateMedia failed: %v", err)
	}

	if updated.Title != created.Title {
		t.Errorf("expected Title %q to be unchanged, got %q", created.Title, updated.Title)
	}
	if updated.Description != "New description" {
		t.Errorf("expected Description %q, got %q", "New description", updated.Description)
	}
	if updated.S3Key == nil || *updated.S3Key != "media/now-on-s3.mp4" {
		t.Errorf("expected S3Key %q, got %v", "media/now-on-s3.mp4", updated.S3Key)
	}
	if updated.YoutubeVideoID != nil {
		t.Errorf("expected YoutubeVideoID to be cleared, got %q", *updated.YoutubeVideoID)
	}
}
