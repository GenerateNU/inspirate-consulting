package mediaAccessRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// test fetching the student id for an existing media access row
func TestGetMediaAccessStudentID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewMediaAccessRepository(db)
	ctx := context.Background()

	mediaID := createTestMedia(t, db, "Get Student ID Test Video")
	studentID := uuid.New()

	granted, err := repo.GrantMediaAccess(ctx, &models.GrantMediaAccessRequestBody{
		StudentID: studentID.String(),
		MediaID:   mediaID,
	})
	if err != nil {
		t.Fatalf("setup GrantMediaAccess failed: %v", err)
	}

	got, err := repo.GetMediaAccessStudentID(ctx, granted.ID)
	if err != nil {
		t.Fatalf("GetMediaAccessStudentID failed: %v", err)
	}
	if got != studentID {
		t.Errorf("expected student id %s, got %s", studentID, got)
	}
}

// test fetching the student id with a nonexistent ID returns an error
func TestGetMediaAccessStudentID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaAccessRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	got, err := repo.GetMediaAccessStudentID(ctx, uuid.NewString())
	if err == nil {
		t.Fatal("expected an error for a nonexistent id, got nil")
	}
	if got != uuid.Nil {
		t.Errorf("expected uuid.Nil, got %s", got)
	}
}
