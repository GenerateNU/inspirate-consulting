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

	input := &models.CreateMediaRequestBody{
		Title:        "Writing Your Personal Statement",
		Description:  "A walkthrough of the Common App essay",
		LengthInMins: 45,
		SchoolYear:   "junior",
		S3Key:        "media/personal-statement.mp4",
	}

	created, err := repo.CreateMedia(ctx, input)
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
	if created.LengthInMins != input.LengthInMins {
		t.Errorf("expected LengthInMins %d, got %d", input.LengthInMins, created.LengthInMins)
	}
	if created.SchoolYear != "junior" {
		t.Errorf("expected SchoolYear %s, got %v", "junior", created.SchoolYear)
	}
	if created.S3Key != input.S3Key {
		t.Errorf("expected S3Key %q, got %q", input.S3Key, created.S3Key)
	}
}

