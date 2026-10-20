package studentstore

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
)

func setupTestRepo(t *testing.T) (*StudentRepository, *pgxpool.Pool) {
	t.Helper()

	db := testutils.SetupTestDB(t)
	return NewStudentRepository(db), db
}

// createTestCounselor inserts a user and a counselor row for them, returning the
// user's id and the counselor's id.
func createTestCounselor(t *testing.T, db *pgxpool.Pool, name string) (uuid.UUID, uuid.UUID) {
	t.Helper()
	ctx := context.Background()

	var userID, counselorID uuid.UUID
	if err := db.QueryRow(ctx,
		`INSERT INTO public.users (name) VALUES ($1) RETURNING id`, name).Scan(&userID); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	if err := db.QueryRow(ctx,
		`INSERT INTO public.counselor (user_id) VALUES ($1) RETURNING id`, userID).Scan(&counselorID); err != nil {
		t.Fatalf("failed to create test counselor: %v", err)
	}

	return userID, counselorID
}

// createTestStudent inserts a user and a student row assigned to counselorID,
// returning the user's id.
func createTestStudent(t *testing.T, db *pgxpool.Pool, name string, counselorID uuid.UUID) uuid.UUID {
	t.Helper()
	ctx := context.Background()

	var userID uuid.UUID
	if err := db.QueryRow(ctx,
		`INSERT INTO public.users (name) VALUES ($1) RETURNING id`, name).Scan(&userID); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	if _, err := db.Exec(ctx,
		`INSERT INTO public.student (user_id, year, gpa, review_balance, counselor_id)
		 VALUES ($1, 'junior', 4, 0, $2)`,
		userID, counselorID); err != nil {
		t.Fatalf("failed to create test student: %v", err)
	}

	return userID
}

// counselorOf reads the counselor currently assigned to the student with the given user id.
func counselorOf(t *testing.T, db *pgxpool.Pool, studentUserID uuid.UUID) uuid.UUID {
	t.Helper()

	var counselorID uuid.UUID
	if err := db.QueryRow(context.Background(),
		`SELECT counselor_id FROM public.student WHERE user_id = $1`, studentUserID).Scan(&counselorID); err != nil {
		t.Fatalf("failed to read student's counselor: %v", err)
	}
	return counselorID
}
