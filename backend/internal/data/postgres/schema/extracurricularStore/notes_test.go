package extracurricularRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"inspirate-consulting/internal/models"
)

func TestExtracurricularNotes_CreateUpdateDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()
	parent := createTestExtracurricular(t, repo, uuid.NewString(), createTestUser(t, db), minimalCreateInput("Debate", "2025-01-01"))

	created, err := repo.CreateExtracurricularNote(ctx, parent.ID, models.ExtracurricularNoteRequest{Body: "First note"})
	if err != nil {
		t.Fatalf("CreateExtracurricularNote failed: %v", err)
	}
	if created.ID == 0 || created.ExtracurricularID != parent.ID || created.Body != "First note" {
		t.Errorf("created = %+v", *created)
	}

	updated, err := repo.UpdateExtracurricularNote(ctx, parent.ID, created.ID, models.ExtracurricularNoteRequest{Body: "Edited"})
	if err != nil {
		t.Fatalf("UpdateExtracurricularNote failed: %v", err)
	}
	if updated.ID != created.ID || updated.Body != "Edited" || !updated.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("updated = %+v", *updated)
	}
	if updated.ModifiedAt.Before(created.ModifiedAt) {
		t.Errorf("modified_at went backwards: %v -> %v", created.ModifiedAt, updated.ModifiedAt)
	}

	detail, err := repo.GetExtracurricular(ctx, parent.ID)
	if err != nil {
		t.Fatalf("GetExtracurricular failed: %v", err)
	}
	if len(detail.Notes) != 1 || detail.Notes[0].Body != "Edited" {
		t.Errorf("notes on detail = %+v", detail.Notes)
	}

	if err := repo.DeleteExtracurricularNote(ctx, parent.ID, created.ID); err != nil {
		t.Fatalf("DeleteExtracurricularNote failed: %v", err)
	}
	if n := countRows(t, db, "extracurricular_notes", parent.ID); n != 0 {
		t.Errorf("notes remaining = %d, want 0", n)
	}
	assertNotFound(t, repo.DeleteExtracurricularNote(ctx, parent.ID, created.ID))
}

func TestExtracurricularNotes_CreateForMissingParentFails(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	if _, err := repo.CreateExtracurricularNote(context.Background(), 999999, models.ExtracurricularNoteRequest{Body: "orphan"}); err == nil {
		t.Error("expected foreign key error, got nil")
	}
}

// A note can't be edited or deleted through a different extracurricular's id.
func TestExtracurricularNotes_ScopedToParent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()
	userID := createTestUser(t, db)
	owner := createTestExtracurricular(t, repo, uuid.NewString(), userID, minimalCreateInput("Owner", "2025-01-01"))
	other := createTestExtracurricular(t, repo, uuid.NewString(), userID, minimalCreateInput("Other", "2025-01-01"))
	noteID := insertNote(t, db, owner.ID, "2025-02-01T00:00:00Z", "Private")

	_, err := repo.UpdateExtracurricularNote(ctx, other.ID, noteID, models.ExtracurricularNoteRequest{Body: "Hijacked"})
	assertNotFound(t, err)
	assertNotFound(t, repo.DeleteExtracurricularNote(ctx, other.ID, noteID))

	detail, err := repo.GetExtracurricular(ctx, owner.ID)
	if err != nil {
		t.Fatalf("GetExtracurricular failed: %v", err)
	}
	if len(detail.Notes) != 1 || detail.Notes[0].Body != "Private" {
		t.Errorf("owner's note changed: %+v", detail.Notes)
	}
}
