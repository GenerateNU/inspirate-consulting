package essayReviewRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func TestLockStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := createTestStudent(t, db, 7)

	student, err := repo.LockStudent(ctx, db, studentID)
	if err != nil {
		t.Fatalf("LockStudent failed: %v", err)
	}

	if student.ID != studentID {
		t.Errorf("expected ID %v, got %v", studentID, student.ID)
	}
	if student.ReviewBalance != 7 {
		t.Errorf("expected ReviewBalance 7, got %d", student.ReviewBalance)
	}
	if student.Year != models.Senior {
		t.Errorf("expected Year %q, got %q", models.Senior, student.Year)
	}
	if student.Organization != nil {
		t.Errorf("expected nil Organization, got %q", *student.Organization)
	}
}

func TestLockStudentNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)

	_, err := repo.LockStudent(context.Background(), db, uuid.New())
	if err == nil {
		t.Fatal("expected an error for an unknown student")
	}
	if httpErr, ok := err.(errs.HTTPError); !ok || httpErr.Code != 404 {
		t.Errorf("expected a 404 HTTPError, got %v", err)
	}
}

func TestLockTransaction(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := createTestStudent(t, db, 5)
	spend := spendFor(t, repo, db, studentID, uuid.New(), 1)

	locked, err := repo.LockTransaction(ctx, db, spend.ID)
	if err != nil {
		t.Fatalf("LockTransaction failed: %v", err)
	}

	if locked.ID != spend.ID {
		t.Errorf("expected ID %v, got %v", spend.ID, locked.ID)
	}
	if locked.Subtotal != -1 {
		t.Errorf("expected Subtotal -1, got %d", locked.Subtotal)
	}
	if locked.Status == nil || *locked.Status != models.ReviewStatusOpen {
		t.Errorf("expected status open, got %v", locked.Status)
	}
}

func TestLockTransactionNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)

	_, err := repo.LockTransaction(context.Background(), db, uuid.New())
	if err == nil {
		t.Fatal("expected an error for an unknown transaction")
	}
	if httpErr, ok := err.(errs.HTTPError); !ok || httpErr.Code != 404 {
		t.Errorf("expected a 404 HTTPError, got %v", err)
	}
}

// The open-review lookup drives the one-review-per-essay rule, so it must stop
// matching as soon as a review is resolved.
func TestFindOpenReviewForEssay(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := createTestStudent(t, db, 9)
	essayID := uuid.New()

	found, err := repo.FindOpenReviewForEssay(ctx, db, essayID)
	if err != nil {
		t.Fatalf("FindOpenReviewForEssay failed: %v", err)
	}
	if found != nil {
		t.Fatalf("expected no open review for an untouched essay, got %v", *found)
	}

	spend := spendFor(t, repo, db, studentID, essayID, 1)

	found, err = repo.FindOpenReviewForEssay(ctx, db, essayID)
	if err != nil {
		t.Fatalf("FindOpenReviewForEssay failed: %v", err)
	}
	if found == nil || *found != spend.ID {
		t.Fatalf("expected the open spend %v, got %v", spend.ID, found)
	}

	if _, err := repo.MarkCompleted(ctx, db, spend.ID); err != nil {
		t.Fatalf("MarkCompleted failed: %v", err)
	}

	found, err = repo.FindOpenReviewForEssay(ctx, db, essayID)
	if err != nil {
		t.Fatalf("FindOpenReviewForEssay failed: %v", err)
	}
	if found != nil {
		t.Errorf("expected a completed review to free the essay, got %v", *found)
	}
}

// A reversal row carries the same essay_id, and must never be mistaken for the review.
func TestFindEssayReviewStatusIgnoresReversals(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := createTestStudent(t, db, 9)
	essayID := uuid.New()

	charge := spendFor(t, repo, db, studentID, essayID, 1)
	reversal, err := repo.InsertRefund(ctx, db, charge)
	if err != nil {
		t.Fatalf("InsertRefund failed: %v", err)
	}
	if _, err := repo.LinkRefund(ctx, db, charge.ID, reversal.ID); err != nil {
		t.Fatalf("LinkRefund failed: %v", err)
	}

	status, err := repo.FindEssayReviewStatus(ctx, db, essayID)
	if err != nil {
		t.Fatalf("FindEssayReviewStatus failed: %v", err)
	}
	if status.TransactionID != charge.ID {
		t.Errorf("expected the spend %v, got %v", charge.ID, status.TransactionID)
	}
	if status.Status != models.ReviewStatusRefunded {
		t.Errorf("expected status refunded, got %q", status.Status)
	}
	if status.EssayID != essayID || status.StudentID != studentID {
		t.Errorf("unexpected ids: essay=%v student=%v", status.EssayID, status.StudentID)
	}
	if status.RequestedAt.IsZero() {
		t.Error("expected RequestedAt to be set")
	}
}

// An essay reviewed more than once reports the newest attempt.
func TestFindEssayReviewStatusReturnsLatestSpend(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := createTestStudent(t, db, 9)
	essayID := uuid.New()

	first := spendFor(t, repo, db, studentID, essayID, 1)
	if _, err := repo.MarkCompleted(ctx, db, first.ID); err != nil {
		t.Fatalf("MarkCompleted failed: %v", err)
	}
	second := spendFor(t, repo, db, studentID, essayID, 1)

	status, err := repo.FindEssayReviewStatus(ctx, db, essayID)
	if err != nil {
		t.Fatalf("FindEssayReviewStatus failed: %v", err)
	}
	if status.TransactionID != second.ID {
		t.Errorf("expected the newest spend %v, got %v", second.ID, status.TransactionID)
	}
	if status.Status != models.ReviewStatusOpen {
		t.Errorf("expected status open, got %q", status.Status)
	}
	if status.CompletedAt != nil {
		t.Errorf("expected nil CompletedAt on an open review, got %v", *status.CompletedAt)
	}
}

func TestFindEssayReviewStatusNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)

	_, err := repo.FindEssayReviewStatus(context.Background(), db, uuid.New())
	if err == nil {
		t.Fatal("expected an error for an unreviewed essay")
	}
	if httpErr, ok := err.(errs.HTTPError); !ok || httpErr.Code != 404 {
		t.Errorf("expected a 404 HTTPError, got %v", err)
	}
}
