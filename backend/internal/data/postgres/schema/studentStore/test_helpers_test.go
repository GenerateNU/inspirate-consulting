package studentRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// setupTestRepo creates a fresh, isolated test database (via testcontainers)
// and returns a repository connected to it. Safe for t.Parallel().
func setupTestRepo(t *testing.T) (*StudentRepository, *pgxpool.Pool) {
	t.Helper()

	db := testutils.SetupTestDB(t)
	return NewStudentRepository(db), db
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

	userID := createTestUser(t, db, "student store test counselor")

	var counselorID uuid.UUID
	if err := db.QueryRow(context.Background(),
		`INSERT INTO public.counselor (user_id) VALUES ($1) RETURNING id`, userID).Scan(&counselorID); err != nil {
		t.Fatalf("setup createTestCounselor failed: %v", err)
	}
	return counselorID
}

// createTestStudent inserts a student directly, since the store has no create method,
// and returns the row as it was written so tests can compare against it.
func createTestStudent(t *testing.T, db *pgxpool.Pool, student models.Student) models.Student {
	t.Helper()

	if err := db.QueryRow(context.Background(),
		`INSERT INTO public.student (user_id, year, organization, gpa, review_balance, counselor_id, student_interest)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		student.UserID, student.Year, student.Organization, student.GPA, student.ReviewBalance, student.CounselorID, student.Interest,
	).Scan(&student.ID); err != nil {
		t.Fatalf("setup createTestStudent failed: %v", err)
	}
	return student
}

// createNamedTestCounselor inserts a counselor (and its backing user) with the
// given name, returning the user's id and the counselor's id.
func createNamedTestCounselor(t *testing.T, db *pgxpool.Pool, name string) (uuid.UUID, uuid.UUID) {
	t.Helper()

	userID := createTestUser(t, db, name)

	var counselorID uuid.UUID
	if err := db.QueryRow(context.Background(),
		`INSERT INTO public.counselor (user_id) VALUES ($1) RETURNING id`, userID).Scan(&counselorID); err != nil {
		t.Fatalf("failed to create test counselor: %v", err)
	}
	return userID, counselorID
}

// createAssignedTestStudent inserts a user and a student row assigned to counselorID,
// returning the user's id.
func createAssignedTestStudent(t *testing.T, db *pgxpool.Pool, name string, counselorID uuid.UUID) uuid.UUID {
	t.Helper()

	userID := createTestUser(t, db, name)
	if _, err := db.Exec(context.Background(),
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
