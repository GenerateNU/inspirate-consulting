package studentRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestAssignCounselor(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	_, oldCounselorID := createNamedTestCounselor(t, db, "Old Counselor")
	_, newCounselorID := createNamedTestCounselor(t, db, "New Counselor")
	studentUserID := createAssignedTestStudent(t, db, "Assign Student", oldCounselorID)

	updated, err := repo.AssignCounselor(ctx, studentUserID, newCounselorID)
	if err != nil {
		t.Fatalf("AssignCounselor failed: %v", err)
	}

	if updated.UserID != studentUserID {
		t.Errorf("expected UserID %v, got %v", studentUserID, updated.UserID)
	}
	if updated.CounselorID == nil || *updated.CounselorID != newCounselorID {
		t.Errorf("expected CounselorID %v, got %v", newCounselorID, updated.CounselorID)
	}
	if got := counselorOf(t, db, studentUserID); got != newCounselorID {
		t.Errorf("expected stored counselor %v, got %v", newCounselorID, got)
	}
}

func TestAssignCounselor_UnknownStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	_, counselorID := createNamedTestCounselor(t, db, "Counselor")

	updated, err := repo.AssignCounselor(ctx, uuid.New(), counselorID)
	if err == nil {
		t.Fatalf("expected an error for an unknown student, got %+v", updated)
	}
	if updated != nil {
		t.Errorf("expected no student on error, got %+v", updated)
	}
}

func TestAssignCounselor_UnknownCounselor(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	_, counselorID := createNamedTestCounselor(t, db, "Counselor")
	studentUserID := createAssignedTestStudent(t, db, "Assign Student", counselorID)

	updated, err := repo.AssignCounselor(ctx, studentUserID, uuid.New())
	if err == nil {
		t.Fatalf("expected an error for an unknown counselor, got %+v", updated)
	}
	if updated != nil {
		t.Errorf("expected no student on error, got %+v", updated)
	}
	if got := counselorOf(t, db, studentUserID); got != counselorID {
		t.Errorf("expected the student to keep counselor %v, got %v", counselorID, got)
	}
}

// A counselor's user id has no student row, so there is no student to assign.
func TestAssignCounselor_CounselorToCounselor(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorUserID, _ := createNamedTestCounselor(t, db, "Assignee Counselor")
	_, otherCounselorID := createNamedTestCounselor(t, db, "Other Counselor")

	updated, err := repo.AssignCounselor(ctx, counselorUserID, otherCounselorID)
	if err == nil {
		t.Fatalf("expected an error when assigning a counselor to a counselor, got %+v", updated)
	}
	if updated != nil {
		t.Errorf("expected no student on error, got %+v", updated)
	}

	var studentRows int
	if err := db.QueryRow(ctx,
		`SELECT count(*) FROM public.student WHERE user_id = $1`, counselorUserID).Scan(&studentRows); err != nil {
		t.Fatalf("failed to count student rows: %v", err)
	}
	if studentRows != 0 {
		t.Errorf("expected no student row for the counselor, got %d", studentRows)
	}
}
