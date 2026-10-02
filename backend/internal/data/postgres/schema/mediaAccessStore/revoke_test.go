package mediaAccessRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// test normal media access revocation with a valid ID
func TestRevokeMediaAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewMediaAccessRepository(db)
	ctx := context.Background()

	mediaID := createTestMedia(t, db, "Revoke Test Video")
	studentID := uuid.NewString()

	granted, err := repo.GrantMediaAccess(ctx, &models.GrantMediaAccessRequestBody{
		StudentID: studentID,
		MediaID:   mediaID,
	})
	if err != nil {
		t.Fatalf("setup GrantMediaAccess failed: %v", err)
	}

	if err := repo.RevokeMediaAccess(ctx, granted.ID); err != nil {
		t.Fatalf("RevokeMediaAccess failed: %v", err)
	}

	// The student should no longer have access to the media.
	results, err := repo.ListAccessibleMedia(ctx, studentID)
	if err != nil {
		t.Fatalf("ListAccessibleMedia failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 media after revoke, got %d", len(results))
	}
}

// test revoking media access with a nonexistent ID returns an error
func TestRevokeMediaAccess_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaAccessRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	err := repo.RevokeMediaAccess(ctx, uuid.NewString())
	if err == nil {
		t.Fatal("expected an error for a nonexistent id, got nil")
	}
}
