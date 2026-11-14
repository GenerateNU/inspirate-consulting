package essayGroupRepository

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

func TestGetEssayGroupStudentID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayGroupRepository(db)
	ctx := context.Background()

	studentID := uuid.New()
	created, err := repo.CreateEssayGroup(ctx, studentID, models.CreateEssayGroupBody{
		Name: "Common App",
	})
	if err != nil {
		t.Fatalf("setup CreateEssayGroup failed: %v", err)
	}

	got, err := repo.GetEssayGroupStudentID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetEssayGroupStudentID failed: %v", err)
	}

	if got != studentID {
		t.Errorf("expected StudentID %v, got %v", studentID, got)
	}
}

func TestGetEssayGroupStudentID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayGroupRepository(db)
	ctx := context.Background()

	got, err := repo.GetEssayGroupStudentID(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected an error getting the student ID of a nonexistent group, got nil")
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
