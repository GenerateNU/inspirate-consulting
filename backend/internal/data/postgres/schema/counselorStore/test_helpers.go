package counselorRepository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// the helpers below insert rows directly via SQL (to avoid cross-package test dependencies)

// inserts a user and a counselor for them, and returns the counselor id.
func createTestCounselor(t *testing.T, db *pgxpool.Pool) uuid.UUID {
	t.Helper()

	var id uuid.UUID
	err := db.QueryRow(
		context.Background(),
<<<<<<< HEAD
		`WITH new_user AS (INSERT INTO public.users (name, email) VALUES ('Test Counselor', gen_random_uuid()::text || '@test.com') RETURNING id)
=======
		`WITH new_user AS (INSERT INTO public.users (name, email) VALUES ('Test Counselor', gen_random_uuid()::text || '@test.invalid') RETURNING id)
>>>>>>> 6f50d7a (fix failing tests (#162))
		INSERT INTO public.counselor (user_id) SELECT id FROM new_user RETURNING id`,
	).Scan(&id)
	if err != nil {
		t.Fatalf("failed to create test counselor: %v", err)
	}

	return id
}

// inserts a user and a student assigned to the given counselor, and returns the student id.
func createTestStudent(t *testing.T, db *pgxpool.Pool, counselorID uuid.UUID, name string, year string) uuid.UUID {
	t.Helper()

	var id uuid.UUID
	err := db.QueryRow(
		context.Background(),
		`WITH new_user AS (INSERT INTO public.users (name, email) VALUES ($1, gen_random_uuid()::text || '@test.com') RETURNING id)
		INSERT INTO public.student (user_id, year, gpa, review_balance, counselor_id)
		SELECT id, $2, 4, 5, $3 FROM new_user RETURNING id`,
		name,
		year,
		counselorID,
	).Scan(&id)
	if err != nil {
		t.Fatalf("failed to create test student: %v", err)
	}

	return id
}

// inserts an essay for the given student and returns its id.
func createTestEssay(t *testing.T, db *pgxpool.Pool, studentID uuid.UUID, essayType string) uuid.UUID {
	t.Helper()

	var id uuid.UUID
	err := db.QueryRow(
		context.Background(),
		`INSERT INTO public.essays (student_id, type, link_to_content) VALUES ($1, $2, 'https://example.com/essay') RETURNING id`,
		studentID,
		essayType,
	).Scan(&id)
	if err != nil {
		t.Fatalf("failed to create test essay: %v", err)
	}

	return id
}

// inserts a review request (spend) for the given essay and returns its id
// a completed review has completed_at set, so it is no longer open
func createTestReview(t *testing.T, db *pgxpool.Pool, studentID uuid.UUID, essayID uuid.UUID, completed bool) uuid.UUID {
	t.Helper()

	var completedAt *time.Time
	if completed {
		now := time.Now()
		completedAt = &now
	}

	var id uuid.UUID
	err := db.QueryRow(
		context.Background(),
		`INSERT INTO public.essay_review_transaction (subtotal, student_id, entry_type, essay_id, completed_at)
		VALUES (-1, $1, 'spend', $2, $3) RETURNING id`,
		studentID,
		essayID,
		completedAt,
	).Scan(&id)
	if err != nil {
		t.Fatalf("failed to create test review: %v", err)
	}

	return id
}

// inserts a to-do item for the given student
func createTestTodoItem(t *testing.T, db *pgxpool.Pool, studentID uuid.UUID, deadline *time.Time, completedAt *time.Time) {
	t.Helper()

	_, err := db.Exec(
		context.Background(),
		`INSERT INTO public.todo_items (student_id, user_id, todo_description, deadline, completed_at)
		VALUES ($1, $1, 'test todo', $2, $3)`,
		studentID,
		deadline,
		completedAt,
	)
	if err != nil {
		t.Fatalf("failed to create test todo item: %v", err)
	}
}
