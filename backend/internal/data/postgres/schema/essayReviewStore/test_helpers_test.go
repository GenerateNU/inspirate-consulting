package essayReviewRepository

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"inspirate-consulting/internal/models"
)

// setupTestRepo connects to the local Supabase Postgres instance.
// Only called from tests that are skipped in short mode
func setupTestRepo(t *testing.T) (*EssayReviewRepository, *pgxpool.Pool) {
	t.Helper()

	connStr := os.Getenv("TEST_DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:54322/postgres"
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		t.Fatalf("failed to parse test db config: %v", err)
	}
	// Mirror db-connection.go so tests exercise the protocol the app uses.
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	t.Cleanup(db.Close)

	return NewEssayReviewRepository(db), db
}

// createTestStudent inserts the users -> counselor -> student chain a ledger
// row needs, and removes all of it when the test finishes.
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

	t.Cleanup(func() {
		// the self-referencing refund column has to be cleared before delete
		if _, err := db.Exec(ctx,
			`UPDATE public.essay_review_transaction SET refund = NULL WHERE student_id = $1`, studentID); err != nil {
			t.Logf("cleanup failed unlinking refunds for student=%s: %v", studentID, err)
		}
		for _, stmt := range []struct {
			sql string
			arg uuid.UUID
		}{
			{`DELETE FROM public.essay_review_transaction WHERE student_id = $1`, studentID},
			{`DELETE FROM public.student WHERE id = $1`, studentID},
			{`DELETE FROM public.counselor WHERE id = $1`, counselorID},
			{`DELETE FROM public.users WHERE id = $1`, userID},
		} {
			if _, err := db.Exec(ctx, stmt.sql, stmt.arg); err != nil {
				t.Logf("cleanup failed for %s: %v", stmt.arg, err)
			}
		}
	})

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
