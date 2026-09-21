package personalCollegeApplicationRepository

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// setupTestRepo connects to the local Supabase Postgres instance.
// Only called from tests that are skipped in short mode
func setupTestRepo(t *testing.T) (*PersonalCollegeApplicationRepository, *pgxpool.Pool) {
	t.Helper()

	connStr := os.Getenv("TEST_DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:54322/postgres"
	}

	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	return NewPersonalCollegeApplicationRepository(db), db
}

// createTestGlobalCollege inserts a global college directly via SQL (to avoid a cross-package test dependency) and 
// returns its id, for use as a valid global_college_id foreign key in personal college application tests.
func createTestGlobalCollege(t *testing.T, db *pgxpool.Pool, schoolName string, schoolLocation string) int64 {
	t.Helper()

	var id int64
	err := db.QueryRow(
		context.Background(),
		`INSERT INTO public.global_colleges (school_name, school_location) VALUES ($1, $2) RETURNING id`,
		schoolName,
		schoolLocation,
	).Scan(&id)
	if err != nil {
		t.Fatalf("failed to create test global college: %v", err)
	}

	t.Cleanup(func() {
		_, err := db.Exec(context.Background(), `DELETE FROM public.global_colleges WHERE id = $1`, id)
		if err != nil {
			t.Logf("cleanup failed for global college id=%d: %v", id, err)
		}
	})

	return id
}

// cleanupApplication deletes a personal college application by ID after the test completes.
func cleanupApplication(t *testing.T, db *pgxpool.Pool, id int64) {
	t.Helper()
	t.Cleanup(func() {
		_, err := db.Exec(context.Background(), `DELETE FROM public.personal_college_applications WHERE id = $1`, id)
		if err != nil {
			t.Logf("cleanup failed for application id=%d: %v", id, err)
		}
	})
}