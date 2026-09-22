package essayReviewRepository

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	dbinterface "inspirate-consulting/internal/data/db-interface"
)

func TestWithTxCommits(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := createTestStudent(t, db, 4)

	err := repo.WithTx(ctx, func(tx dbinterface.QueryInterface) error {
		if err := repo.AdjustStudentBalance(ctx, tx, studentID, -1); err != nil {
			return err
		}
		return repo.InsertAdjustment(ctx, tx, studentID, -1)
	})
	if err != nil {
		t.Fatalf("WithTx failed: %v", err)
	}

	if balance := readBalance(t, db, studentID); balance != 3 {
		t.Errorf("expected balance 3, got %d", balance)
	}
	if sum := readLedgerSum(t, db, studentID); sum != -1 {
		t.Errorf("expected ledger sum -1, got %d", sum)
	}
}

// Everything written before the failure must disappear, which is what stops a
// half-applied refund leaving the ledger out of step with the balance.
func TestWithTxRollsBackEveryWrite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := createTestStudent(t, db, 4)
	boom := errors.New("boom")

	err := repo.WithTx(ctx, func(tx dbinterface.QueryInterface) error {
		if err := repo.AdjustStudentBalance(ctx, tx, studentID, -1); err != nil {
			return err
		}
		if err := repo.InsertAdjustment(ctx, tx, studentID, -1); err != nil {
			return err
		}
		if _, err := repo.InsertSpend(ctx, tx, studentID, uuid.New(), 1); err != nil {
			return err
		}
		return boom
	})

	if !errors.Is(err, boom) {
		t.Fatalf("expected the callback error back, got %v", err)
	}
	if balance := readBalance(t, db, studentID); balance != 4 {
		t.Errorf("expected balance rolled back to 4, got %d", balance)
	}
	if sum := readLedgerSum(t, db, studentID); sum != 0 {
		t.Errorf("expected no ledger rows to survive, got sum %d", sum)
	}

	var rows int
	if err := db.QueryRow(ctx,
		`SELECT count(*) FROM public.essay_review_transaction WHERE student_id = $1`, studentID).Scan(&rows); err != nil {
		t.Fatalf("failed to count rows: %v", err)
	}
	if rows != 0 {
		t.Errorf("expected 0 rows after rollback, got %d", rows)
	}
}

// The invariant the ledger exists to provide.
func TestLedgerReconcilesWithBalance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	// start at zero so every movement goes through the ledger
	studentID := createTestStudent(t, db, 0)

	reconcile := func(step string, want int) {
		t.Helper()
		balance, sum := readBalance(t, db, studentID), readLedgerSum(t, db, studentID)
		if balance != want || sum != want {
			t.Fatalf("%s: balance=%d ledger=%d, want %d for both", step, balance, sum, want)
		}
	}

	reconcile("start", 0)

	// grant 5
	if err := repo.WithTx(ctx, func(tx dbinterface.QueryInterface) error {
		if err := repo.SetStudentBalance(ctx, tx, studentID, 5); err != nil {
			return err
		}
		return repo.InsertAdjustment(ctx, tx, studentID, 5)
	}); err != nil {
		t.Fatalf("grant failed: %v", err)
	}
	reconcile("after grant", 5)

	// spend 2
	var charge = func() uuid.UUID {
		var id uuid.UUID
		if err := repo.WithTx(ctx, func(tx dbinterface.QueryInterface) error {
			if err := repo.AdjustStudentBalance(ctx, tx, studentID, -2); err != nil {
				return err
			}
			spend, err := repo.InsertSpend(ctx, tx, studentID, uuid.New(), 2)
			if err != nil {
				return err
			}
			id = spend.ID
			return nil
		}); err != nil {
			t.Fatalf("spend failed: %v", err)
		}
		return id
	}()
	reconcile("after spend", 3)

	// refund it
	if err := repo.WithTx(ctx, func(tx dbinterface.QueryInterface) error {
		locked, err := repo.LockTransaction(ctx, tx, charge)
		if err != nil {
			return err
		}
		reversal, err := repo.InsertRefund(ctx, tx, locked)
		if err != nil {
			return err
		}
		if _, err := repo.LinkRefund(ctx, tx, locked.ID, reversal.ID); err != nil {
			return err
		}
		return repo.AdjustStudentBalance(ctx, tx, studentID, reversal.Subtotal)
	}); err != nil {
		t.Fatalf("refund failed: %v", err)
	}
	reconcile("after refund", 5)
}
