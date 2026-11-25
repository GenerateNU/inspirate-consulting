package mediaRepository

import (
	"context"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// test normal media creation with valid input
func TestCreateMedia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	schoolYear := models.Junior

	input := &models.CreateMediaRequestBody{
		Title:        "Writing Your Personal Statement",
		Description:  "A walkthrough of the Common App essay",
		LengthInMins: ptr(45),
		SchoolYear:   &schoolYear,
		S3Key:        ptr("media/personal-statement.mp4"),
	}

	created, err := repo.CreateMedia(ctx, input, nil)
	if err != nil {
		t.Fatalf("CreateMedia failed: %v", err)
	}

	if created.ID == "" {
		t.Error("expected a generated ID, got empty string")
	}
	if created.Title != input.Title {
		t.Errorf("expected Title %q, got %q", input.Title, created.Title)
	}
	if created.Description != input.Description {
		t.Errorf("expected Description %q, got %q", input.Description, created.Description)
	}
	if created.LengthInMins == nil || *created.LengthInMins != *input.LengthInMins {
		t.Errorf("expected LengthInMins %d, got %v", *input.LengthInMins, created.LengthInMins)
	}
	if created.SchoolYear == nil {
		t.Errorf("expected SchoolYear %q, got nil", *input.SchoolYear)
	}
	if *created.SchoolYear != *input.SchoolYear {
		t.Errorf("expected SchoolYear %q, got %q", *input.SchoolYear, *created.SchoolYear)
	}
	if created.S3Key == nil || *created.S3Key != *input.S3Key {
		t.Errorf("expected S3Key %q, got %v", *input.S3Key, created.S3Key)
	}
	if created.YoutubeVideoID != nil {
		t.Errorf("expected nil YoutubeVideoID, got %q", *created.YoutubeVideoID)
	}
}

// test creating YouTube-backed media with no s3_key or length
func TestCreateMedia_Youtube(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	created, err := repo.CreateMedia(ctx, &models.CreateMediaRequestBody{
		Title:       "FAFSA Walkthrough",
		Description: "Hosted on YouTube",
	}, ptr("dQw4w9WgXcQ"))
	if err != nil {
		t.Fatalf("CreateMedia failed: %v", err)
	}

	if created.YoutubeVideoID == nil || *created.YoutubeVideoID != "dQw4w9WgXcQ" {
		t.Errorf("expected YoutubeVideoID %q, got %v", "dQw4w9WgXcQ", created.YoutubeVideoID)
	}
	if created.S3Key != nil {
		t.Errorf("expected nil S3Key, got %q", *created.S3Key)
	}
	if created.LengthInMins != nil {
		t.Errorf("expected nil LengthInMins, got %d", *created.LengthInMins)
	}
}
