package essayReviewRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"inspirate-consulting/internal/models"
)

func TestInsertSpend(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := createTestStudent(t, db, 5)
	essayID := uuid.New()

	created, err := repo.InsertSpend(ctx, db, studentID, essayID, 3)
	if err != nil {
		t.Fatalf("InsertSpend failed: %v", err)
	}

	if created.Subtotal != -3 {
		t.Errorf("expected Subtotal -3, got %d", created.Subtotal)
	}
	if created.EntryType != models.ReviewEntrySpend {
		t.Errorf("expected EntryType %q, got %q", models.ReviewEntrySpend, created.EntryType)
	}
	if created.StudentID != studentID {
		t.Errorf("expected StudentID %v, got %v", studentID, created.StudentID)
	}
	if created.EssayID == nil || *created.EssayID != essayID {
		t.Errorf("expected EssayID %v, got %v", essayID, created.EssayID)
	}
	if created.Status == nil || *created.Status != models.ReviewStatusOpen {
		t.Errorf("expected status open, got %v", created.Status)
	}
	if created.CompletedAt != nil || created.Refund != nil {
		t.Errorf("expected a fresh spend to be unresolved, got completed=%v refund=%v", created.CompletedAt, created.Refund)
	}
	if created.ActorID != nil {
		t.Errorf("expected ActorID nil until auth supplies it, got %v", created.ActorID)
	}
	if created.CreatedAt.IsZero() || created.UpdatedAt.IsZero() {
		t.Error("expected CreatedAt and UpdatedAt to be set")
	}
}

// A reversal negates the charge so the pair sums to zero.
func TestInsertRefund(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := createTestStudent(t, db, 5)
	essayID := uuid.New()
	charge := spendFor(t, repo, db, studentID, essayID, 2)

	reversal, err := repo.InsertRefund(ctx, db, charge)
	if err != nil {
		t.Fatalf("InsertRefund failed: %v", err)
	}

	if reversal.Subtotal != 2 {
		t.Errorf("expected Subtotal 2, got %d", reversal.Subtotal)
	}
	if charge.Subtotal+reversal.Subtotal != 0 {
		t.Errorf("charge and reversal should sum to zero, got %d", charge.Subtotal+reversal.Subtotal)
	}
	if reversal.EntryType != models.ReviewEntryRefund {
		t.Errorf("expected EntryType %q, got %q", models.ReviewEntryRefund, reversal.EntryType)
	}
	if reversal.StudentID != studentID {
		t.Errorf("expected StudentID %v, got %v", studentID, reversal.StudentID)
	}
	if reversal.EssayID == nil || *reversal.EssayID != essayID {
		t.Errorf("expected the charge's essay %v, got %v", essayID, reversal.EssayID)
	}
	// only spends have a lifecycle
	if reversal.Status != nil {
		t.Errorf("expected nil status on a reversal, got %q", *reversal.Status)
	}
}

// An adjustment belongs to no essay.
func TestInsertAdjustment(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := createTestStudent(t, db, 0)

	if err := repo.InsertAdjustment(ctx, db, studentID, 5); err != nil {
		t.Fatalf("InsertAdjustment failed: %v", err)
	}

	var id uuid.UUID
	var subtotal int
	var entryType models.ReviewEntryType
	var essayID *uuid.UUID
	if err := db.QueryRow(ctx,
		`SELECT id, subtotal, entry_type, essay_id FROM public.essay_review_transaction WHERE student_id = $1`,
		studentID).Scan(&id, &subtotal, &entryType, &essayID); err != nil {
		t.Fatalf("failed to read the adjustment back: %v", err)
	}

	if subtotal != 5 {
		t.Errorf("expected Subtotal 5, got %d", subtotal)
	}
	if entryType != models.ReviewEntryAdjustment {
		t.Errorf("expected EntryType %q, got %q", models.ReviewEntryAdjustment, entryType)
	}
	if essayID != nil {
		t.Errorf("expected nil EssayID on an adjustment, got %v", *essayID)
	}
	if status := readStatus(t, db, id); status != nil {
		t.Errorf("expected nil status on an adjustment, got %q", *status)
	}
}

// A negative adjustment records an admin lowering a balance.
func TestInsertAdjustmentNegative(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := createTestStudent(t, db, 0)

	if err := repo.InsertAdjustment(ctx, db, studentID, -3); err != nil {
		t.Fatalf("InsertAdjustment failed: %v", err)
	}

	if sum := readLedgerSum(t, db, studentID); sum != -3 {
		t.Errorf("expected ledger sum -3, got %d", sum)
	}
}
