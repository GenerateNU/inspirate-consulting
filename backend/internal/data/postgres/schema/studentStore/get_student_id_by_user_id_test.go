package studentRepository

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"

	"inspirate-consulting/internal/errs"
)

func TestGetStudentIDByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	_, counselorID := createNamedTestCounselor(t, db, "Counselor")
	studentUserID := createAssignedTestStudent(t, db, "Lookup Student", counselorID)

	var wantID uuid.UUID
	if err := db.QueryRow(ctx,
		`SELECT id FROM public.student WHERE user_id = $1`, studentUserID).Scan(&wantID); err != nil {
		t.Fatalf("failed to read student id: %v", err)
	}

	got, err := repo.GetStudentIDByUserID(ctx, studentUserID)
	if err != nil {
		t.Fatalf("GetStudentIDByUserID failed: %v", err)
	}
	if got != wantID {
		t.Errorf("expected student id %v, got %v", wantID, got)
	}
}

func TestGetStudentIDByUserID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	got, err := repo.GetStudentIDByUserID(ctx, uuid.New())
	if err == nil {
		t.Fatalf("expected an error for an unknown user, got %v", got)
	}
	if got != uuid.Nil {
		t.Errorf("expected uuid.Nil on error, got %v", got)
	}

	var httpErr errs.HTTPError
	if !errors.As(err, &httpErr) || httpErr.Code != http.StatusNotFound {
		t.Errorf("expected a not found error, got %v", err)
	}
}
