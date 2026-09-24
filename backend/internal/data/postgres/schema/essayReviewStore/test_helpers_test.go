package essayReviewRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

func setupTestRepo(t *testing.T) (*EssayReviewRepository, *pgxpool.Pool) {
	t.Helper()

	db := testutils.SetupTestDB(t)
	return NewEssayReviewRepository(db), db
}

func createTestStudent(t *testing.T, db *pgxpool.Pool, reviewBalance int) uuid.UUID {
	t.Helper()
	ctx := context.Background()

	var userID, counselorID, studentID uuid.UUID
	if err := db.QueryRow(ctx,
		`INSERT INTO public.users (name) VALUES ('essay review test') RETURNING id`).Scan(&userID); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	if err := db.QueryRow(ctx,
		`INSERT INTO public.counselor (user_id) VALUES ($1) RETURNING id`, userID).Scan(&counselorID); err != nil {
		t.Fatalf("failed to create test counselor: %v", err)
	}
	if err := db.QueryRow(ctx,
		`INSERT INTO public.student (user_id, year, gpa, review_balance, counselor_id)
		 VALUES ($1, 'senior', 4, $2, $3) RETURNING id`,
		userID, reviewBalance, counselorID).Scan(&studentID); err != nil {
		t.Fatalf("failed to create test student: %v", err)
	}

	return studentID
}

// readBalance returns a student's current review_balance.
func readBalance(t *testing.T, db *pgxpool.Pool, studentID uuid.UUID) int {
	t.Helper()

	var balance int
	if err := db.QueryRow(context.Background(),
		`SELECT review_balance FROM public.student WHERE id = $1`, studentID).Scan(&balance); err != nil {
		t.Fatalf("failed to read balance: %v", err)
	}
	return balance
}

// readLedgerSum returns SUM(subtotal) for a student, which must always equal
// their review_balance.
func readLedgerSum(t *testing.T, db *pgxpool.Pool, studentID uuid.UUID) int {
	t.Helper()

	var sum int
	if err := db.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(subtotal), 0) FROM public.essay_review_transaction WHERE student_id = $1`,
		studentID).Scan(&sum); err != nil {
		t.Fatalf("failed to read ledger sum: %v", err)
	}
	return sum
}

// readStatus returns the generated status column, which is null for rows that
// are not spends.
func readStatus(t *testing.T, db *pgxpool.Pool, id uuid.UUID) *models.ReviewStatus {
	t.Helper()

	var status *models.ReviewStatus
	if err := db.QueryRow(context.Background(),
		`SELECT status FROM public.essay_review_transaction WHERE id = $1`, id).Scan(&status); err != nil {
		t.Fatalf("failed to read status: %v", err)
	}
	return status
}

// spendFor inserts an open spend and returns it.
func spendFor(t *testing.T, repo *EssayReviewRepository, db *pgxpool.Pool, studentID, essayID uuid.UUID, amount int) *models.EssayReviewTransaction {
	t.Helper()

	spend, err := repo.InsertSpend(context.Background(), db, studentID, essayID, amount)
	if err != nil {
		t.Fatalf("failed to insert test spend: %v", err)
	}
	return spend
}
