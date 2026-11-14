package activityRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
)

// setupTestRepo creates a fresh, isolated test database (via testcontainers)
// and returns a repository connected to it. Safe for t.Parallel().
func setupTestRepo(t *testing.T) (*ActivityRepository, *pgxpool.Pool) {
	t.Helper()

	db := testutils.SetupTestDB(t)
	return NewActivityRepository(db), db
}

// createTestUser inserts a user directly and returns its id.
func createTestUser(t *testing.T, db *pgxpool.Pool, name string) uuid.UUID {
	t.Helper()

	var userID uuid.UUID
	if err := db.QueryRow(context.Background(),
		`INSERT INTO public.users (name, email) VALUES ($1, $2) RETURNING id`,
		name, uuid.NewString()+"@test.com").Scan(&userID); err != nil {
		t.Fatalf("setup createTestUser(%q) failed: %v", name, err)
	}
	return userID
}

// createTestCounselor inserts a counselor (and its backing user) and returns the counselor id.
func createTestCounselor(t *testing.T, db *pgxpool.Pool) uuid.UUID {
	t.Helper()

	userID := createTestUser(t, db, "activity store test counselor")

	var counselorID uuid.UUID
	if err := db.QueryRow(context.Background(),
		`INSERT INTO public.counselor (user_id) VALUES ($1) RETURNING id`, userID).Scan(&counselorID); err != nil {
		t.Fatalf("setup createTestCounselor failed: %v", err)
	}
	return counselorID
}

// createTestStudent inserts a user and a student row assigned to counselorID,
// returning the student's own id (not the user's id) - Activity.student_id
// references student(id).
func createTestStudent(t *testing.T, db *pgxpool.Pool, counselorID uuid.UUID) uuid.UUID {
	t.Helper()

	userID := createTestUser(t, db, "activity store test student")

	var studentID uuid.UUID
	if err := db.QueryRow(context.Background(),
		`INSERT INTO public.student (user_id, year, gpa, review_balance, counselor_id)
		 VALUES ($1, 'junior', 4, 0, $2) RETURNING id`,
		userID, counselorID).Scan(&studentID); err != nil {
		t.Fatalf("setup createTestStudent failed: %v", err)
	}
	return studentID
}

// createTestEssay inserts a bare essay row and returns its id.
func createTestEssay(t *testing.T, db *pgxpool.Pool, studentID uuid.UUID) uuid.UUID {
	t.Helper()

	var essayID uuid.UUID
	if err := db.QueryRow(context.Background(),
		`INSERT INTO public.essays (student_id, type, link_to_content)
		 VALUES ($1, 'personal-statement', 'https://example.com/essay')
		 RETURNING id`,
		studentID).Scan(&essayID); err != nil {
		t.Fatalf("setup createTestEssay failed: %v", err)
	}
	return essayID
}
