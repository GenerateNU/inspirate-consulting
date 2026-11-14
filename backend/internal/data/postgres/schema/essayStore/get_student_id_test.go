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

func TestGetEssayStudentID(t *testing.T) {
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

	got, err := repo.GetEssayStudentID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetEssayStudentID failed: %v", err)
	}

	if got != studentID {
		t.Errorf("expected StudentID %v, got %v", studentID, got)
	}
}

func TestGetEssayStudentID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayRepository(db)
	ctx := context.Background()

	got, err := repo.GetEssayStudentID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected an error getting the student ID of a nonexistent essay, got nil")
	}
	if got != uuid.Nil {
		t.Errorf("expected the nil UUID on error, got %v", got)
	}

	var httpErr errs.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected an errs.HTTPError, got %T: %v", err, err)
	}
	if httpErr.GetStatus() != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, httpErr.GetStatus())
	}
}
