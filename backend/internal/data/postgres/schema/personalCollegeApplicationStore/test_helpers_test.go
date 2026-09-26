package personalCollegeApplicationRepository

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

	return id
}
