package essayRepository

import (
	"context"
	"errors"
	"net/http"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
)

func TestDeleteEssay(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayRepository(db)
	ctx := context.Background()

	studentID := uuid.New()
	created := createEssay(t, repo, studentID, models.CreateEssayBody{
		Type:          "personal-statement",
		LinkToContent: "https://docs.google.com/document/d/abc123",
	})

	if err := repo.DeleteEssay(ctx, created.ID, studentID); err != nil {
		t.Fatalf("DeleteEssay failed: %v", err)
	}

	essays, err := repo.GetEssaysFromStudent(ctx, studentID)
	if err != nil {
		t.Fatalf("GetEssaysFromStudent failed: %v", err)
	}
	if len(essays) != 0 {
		t.Errorf("expected the essay to be deleted, got %d essays", len(essays))
	}
}

func TestDeleteEssay_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayRepository(db)
	ctx := context.Background()

	err := repo.DeleteEssay(ctx, uuid.New(), uuid.New())
	if err == nil {
		t.Fatal("expected an error deleting a nonexistent essay, got nil")
	}

	var httpErr errs.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected an errs.HTTPError, got %T: %v", err, err)
	}
	if httpErr.GetStatus() != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, httpErr.GetStatus())
	}
}

// a student must not be able to delete another student's essay
func TestDeleteEssay_WrongStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayRepository(db)
	ctx := context.Background()

	ownerID := uuid.New()
	created := createEssay(t, repo, ownerID, models.CreateEssayBody{
		Type:          "personal-statement",
		LinkToContent: "https://docs.google.com/document/d/abc123",
	})

	err := repo.DeleteEssay(ctx, created.ID, uuid.New())
	if err == nil {
		t.Fatal("expected an error deleting another student's essay, got nil")
	}

	var httpErr errs.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected an errs.HTTPError, got %T: %v", err, err)
	}
	if httpErr.GetStatus() != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, httpErr.GetStatus())
	}

	// The owner's essay must survive the attempt.
	essays, err := repo.GetEssaysFromStudent(ctx, ownerID)
	if err != nil {
		t.Fatalf("GetEssaysFromStudent failed: %v", err)
	}
	if len(essays) != 1 {
		t.Errorf("expected the owner's essay to still exist, got %d essays", len(essays))
	}
}
