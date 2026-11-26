package mediaAccessRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// test started then completed sets both timestamps, and repeat calls keep the original values
func TestRecordMediaProgress_Idempotent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewMediaAccessRepository(db)
	ctx := context.Background()

	studentID := uuid.NewString()
	mediaID := createTestMedia(t, db, "Progress Video")
	if _, err := repo.GrantMediaAccess(ctx, &models.GrantMediaAccessRequestBody{StudentID: studentID, MediaID: mediaID}); err != nil {
		t.Fatalf("setup GrantMediaAccess failed: %v", err)
	}

	started, err := repo.RecordMediaProgress(ctx, studentID, mediaID, false)
	if err != nil {
		t.Fatalf("RecordMediaProgress(started) failed: %v", err)
	}
	if started.StartedAt == nil {
		t.Fatal("expected StartedAt to be set")
	}
	if started.CompletedAt != nil {
		t.Errorf("expected CompletedAt to be nil after started, got %v", *started.CompletedAt)
	}

	startedAgain, err := repo.RecordMediaProgress(ctx, studentID, mediaID, false)
	if err != nil {
		t.Fatalf("repeat RecordMediaProgress(started) failed: %v", err)
	}
	if !startedAgain.StartedAt.Equal(*started.StartedAt) {
		t.Errorf("expected StartedAt to stay %v, got %v", *started.StartedAt, *startedAgain.StartedAt)
	}

	completed, err := repo.RecordMediaProgress(ctx, studentID, mediaID, true)
	if err != nil {
		t.Fatalf("RecordMediaProgress(completed) failed: %v", err)
	}
	if completed.CompletedAt == nil {
		t.Fatal("expected CompletedAt to be set")
	}
	if !completed.StartedAt.Equal(*started.StartedAt) {
		t.Errorf("expected StartedAt to stay %v after completing, got %v", *started.StartedAt, *completed.StartedAt)
	}

	completedAgain, err := repo.RecordMediaProgress(ctx, studentID, mediaID, true)
	if err != nil {
		t.Fatalf("repeat RecordMediaProgress(completed) failed: %v", err)
	}
	if !completedAgain.CompletedAt.Equal(*completed.CompletedAt) {
		t.Errorf("expected CompletedAt to stay %v, got %v", *completed.CompletedAt, *completedAgain.CompletedAt)
	}

	// A late "started" (e.g. reopening the video) must not undo completion
	reopened, err := repo.RecordMediaProgress(ctx, studentID, mediaID, false)
	if err != nil {
		t.Fatalf("RecordMediaProgress(started) after completion failed: %v", err)
	}
	if reopened.CompletedAt == nil || !reopened.CompletedAt.Equal(*completed.CompletedAt) {
		t.Errorf("expected CompletedAt to stay %v after reopening, got %v", *completed.CompletedAt, reopened.CompletedAt)
	}
}

// test completing without a prior started call also sets StartedAt
func TestRecordMediaProgress_CompleteWithoutStart(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewMediaAccessRepository(db)
	ctx := context.Background()

	studentID := uuid.NewString()
	mediaID := createTestMedia(t, db, "Skip Ahead Video")
	if _, err := repo.GrantMediaAccess(ctx, &models.GrantMediaAccessRequestBody{StudentID: studentID, MediaID: mediaID}); err != nil {
		t.Fatalf("setup GrantMediaAccess failed: %v", err)
	}

	completed, err := repo.RecordMediaProgress(ctx, studentID, mediaID, true)
	if err != nil {
		t.Fatalf("RecordMediaProgress(completed) failed: %v", err)
	}
	if completed.StartedAt == nil || completed.CompletedAt == nil {
		t.Errorf("expected both StartedAt and CompletedAt to be set, got %v / %v", completed.StartedAt, completed.CompletedAt)
	}
}

// test recording progress without access returns an error and touches nothing
func TestRecordMediaProgress_NoAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewMediaAccessRepository(db)
	ctx := context.Background()

	grantedStudent := uuid.NewString()
	otherStudent := uuid.NewString()
	mediaID := createTestMedia(t, db, "Someone Else's Video")
	if _, err := repo.GrantMediaAccess(ctx, &models.GrantMediaAccessRequestBody{StudentID: grantedStudent, MediaID: mediaID}); err != nil {
		t.Fatalf("setup GrantMediaAccess failed: %v", err)
	}

	if _, err := repo.RecordMediaProgress(ctx, otherStudent, mediaID, true); err == nil {
		t.Fatal("expected an error for a student without access, got nil")
	}

	// The granted student's row must be untouched
	results, err := repo.ListAccessibleMedia(ctx, grantedStudent, "", 0.2, 10, 0)
	if err != nil {
		t.Fatalf("ListAccessibleMedia failed: %v", err)
	}
	if len(results) != 1 || results[0].StartedAt != nil || results[0].CompletedAt != nil {
		t.Errorf("expected granted student's progress to be untouched, got %+v", results)
	}
}
