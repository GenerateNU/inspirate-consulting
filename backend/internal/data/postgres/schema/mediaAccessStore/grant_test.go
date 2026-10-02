package mediaAccessRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// test normal media access grant with valid input
func TestGrantMediaAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewMediaAccessRepository(db)
	ctx := context.Background()

	mediaID := createTestMedia(t, db, "Grant Test Video")

	input := &models.GrantMediaAccessRequestBody{
		StudentID: uuid.NewString(),
		MediaID:   mediaID,
	}

	created, err := repo.GrantMediaAccess(ctx, input)
	if err != nil {
		t.Fatalf("GrantMediaAccess failed: %v", err)
	}

	if created.ID == "" {
		t.Error("expected a generated ID, got empty string")
	}
	if created.StudentID != input.StudentID {
		t.Errorf("expected StudentID %q, got %q", input.StudentID, created.StudentID)
	}
	if created.MediaID != input.MediaID {
		t.Errorf("expected MediaID %q, got %q", input.MediaID, created.MediaID)
	}
}

// test granting access to a nonexistent media fails, since media_id is a foreign key
func TestGrantMediaAccess_NonexistentMedia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaAccessRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	input := &models.GrantMediaAccessRequestBody{
		StudentID: uuid.NewString(),
		MediaID:   uuid.NewString(),
	}

	_, err := repo.GrantMediaAccess(ctx, input)
	if err == nil {
		t.Fatal("expected an error for a nonexistent media_id, got nil")
	}
}
