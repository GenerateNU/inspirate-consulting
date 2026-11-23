package todoItemRepository

import (
	"context"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// setupTestRepo creates a fresh, isolated test database (via testcontainers)
// and returns a repository connected to it. Safe for t.Parallel().
func setupTestRepo(t *testing.T) (*TodoItemRepository, *pgxpool.Pool) {
	t.Helper()
	db := testutils.SetupTestDB(t)
	return NewTodoItemRepository(db), db
}

// createTestCounselor inserts a user and a counselor, returning the counselor id and user id.
func createTestCounselor(t *testing.T, db *pgxpool.Pool) (uuid.UUID, uuid.UUID) {
	t.Helper()
	var counselorID, userID uuid.UUID
	err := db.QueryRow(context.Background(),
		`WITH new_user AS (
			INSERT INTO public.users (name, email)
			VALUES ('Test Counselor', gen_random_uuid()::text || '@test.com')
			RETURNING id
		)
		INSERT INTO public.counselor (user_id)
		SELECT id FROM new_user
		RETURNING id, user_id`,
	).Scan(&counselorID, &userID)
	if err != nil {
		t.Fatalf("failed to create test counselor: %v", err)
	}
	return counselorID, userID
}

// createTestStudentWithCounselor inserts a user and a student assigned to a counselor.
func createTestStudentWithCounselor(t *testing.T, db *pgxpool.Pool, counselorID uuid.UUID, year string) uuid.UUID {
	t.Helper()
	var studentID uuid.UUID
	err := db.QueryRow(context.Background(),
		`WITH new_user AS (
			INSERT INTO public.users (name, email)
			VALUES ('Test Student', gen_random_uuid()::text || '@test.com')
			RETURNING id
		)
		INSERT INTO public.student (user_id, year, gpa, review_balance, counselor_id)
		SELECT id, $1, 4, 0, $2 FROM new_user
		RETURNING id`,
		year, counselorID,
	).Scan(&studentID)
	if err != nil {
		t.Fatalf("failed to create test student: %v", err)
	}
	return studentID
}
