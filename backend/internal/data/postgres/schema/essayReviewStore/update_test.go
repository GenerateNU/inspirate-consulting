package essayReviewRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func TestSetStudentBalance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := createTestStudent(t, db, 2)

	if err := repo.SetStudentBalance(ctx, db, studentID, 9); err != nil {
		t.Fatalf("SetStudentBalance failed: %v", err)
	}

	if balance := readBalance(t, db, studentID); balance != 9 {
		t.Errorf("expected balance 9, got %d", balance)
	}
}

func TestSetStudentBalanceNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)

	err := repo.SetStudentBalance(context.Background(), db, uuid.New(), 1)
	if err == nil {
		t.Fatal("expected an error for an unknown student")
	}
	if httpErr, ok := err.(errs.HTTPError); !ok || httpErr.Code != 404 {
		t.Errorf("expected a 404 HTTPError, got %v", err)
	}
}

// The delta is applied by the database so a concurrent writer cannot be lost.
func TestAdjustStudentBalance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := createTestStudent(t, db, 5)

	if err := repo.AdjustStudentBalance(ctx, db, studentID, -2); err != nil {
		t.Fatalf("AdjustStudentBalance failed: %v", err)
	}
	if balance := readBalance(t, db, studentID); balance != 3 {
		t.Errorf("expected balance 3 after debit, got %d", balance)
	}

	if err := repo.AdjustStudentBalance(ctx, db, studentID, 4); err != nil {
		t.Fatalf("AdjustStudentBalance failed: %v", err)
	}
	if balance := readBalance(t, db, studentID); balance != 7 {
		t.Errorf("expected balance 7 after credit, got %d", balance)
	}
}

func TestAdjustStudentBalanceNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)

	err := repo.AdjustStudentBalance(context.Background(), db, uuid.New(), 1)
	if err == nil {
		t.Fatal("expected an error for an unknown student")
	}
	if httpErr, ok := err.(errs.HTTPError); !ok || httpErr.Code != 404 {
		t.Errorf("expected a 404 HTTPError, got %v", err)
	}
}

// Completing must recompute the generated status column.
func TestMarkCompleted(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := createTestStudent(t, db, 5)
	spend := spendFor(t, repo, db, studentID, uuid.New(), 1)

	completed, err := repo.MarkCompleted(ctx, db, spend.ID)
	if err != nil {
		t.Fatalf("MarkCompleted failed: %v", err)
	}

	if completed.CompletedAt == nil {
		t.Error("expected CompletedAt to be set")
	}
	if completed.Status == nil || *completed.Status != models.ReviewStatusCompleted {
		t.Errorf("expected status completed, got %v", completed.Status)
	}
	if completed.Subtotal != spend.Subtotal {
		t.Errorf("completing must not change the amount: got %d, want %d", completed.Subtotal, spend.Subtotal)
	}
	if balance := readBalance(t, db, studentID); balance != 5 {
		t.Errorf("expected balance untouched at 5, got %d", balance)
	}
}

// Linking a reversal must recompute status to refunded.
func TestLinkRefund(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := createTestStudent(t, db, 5)
	charge := spendFor(t, repo, db, studentID, uuid.New(), 1)

	reversal, err := repo.InsertRefund(ctx, db, charge)
	if err != nil {
		t.Fatalf("InsertRefund failed: %v", err)
	}

	refunded, err := repo.LinkRefund(ctx, db, charge.ID, reversal.ID)
	if err != nil {
		t.Fatalf("LinkRefund failed: %v", err)
	}

	if refunded.Refund == nil || *refunded.Refund != reversal.ID {
		t.Errorf("expected refund pointer %v, got %v", reversal.ID, refunded.Refund)
	}
	if refunded.Status == nil || *refunded.Status != models.ReviewStatusRefunded {
		t.Errorf("expected status refunded, got %v", refunded.Status)
	}
	if refunded.CompletedAt != nil {
		t.Error("expected CompletedAt to stay nil")
	}
}

// The unique constraint on refund stops one reversal serving two charges.
func TestLinkRefundRejectsReusedReversal(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := createTestStudent(t, db, 9)
	first := spendFor(t, repo, db, studentID, uuid.New(), 1)
	second := spendFor(t, repo, db, studentID, uuid.New(), 1)

	reversal, err := repo.InsertRefund(ctx, db, first)
	if err != nil {
		t.Fatalf("InsertRefund failed: %v", err)
	}
	if _, err := repo.LinkRefund(ctx, db, first.ID, reversal.ID); err != nil {
		t.Fatalf("LinkRefund failed: %v", err)
	}

	if _, err := repo.LinkRefund(ctx, db, second.ID, reversal.ID); err == nil {
		t.Fatal("expected the unique constraint to reject a reused reversal")
	}
}

func TestMarkCompletedNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)

	_, err := repo.MarkCompleted(context.Background(), db, uuid.New())
	if err == nil {
		t.Fatal("expected an error for an unknown transaction")
	}
	if httpErr, ok := err.(errs.HTTPError); !ok || httpErr.Code != 404 {
		t.Errorf("expected a 404 HTTPError, got %v", err)
	}
}
