package globalCollegeRepository

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// setupTestRepo connects to the local Supabase Postgres instance.
// Only called from tests that are skipped in short mode
func setupTestRepo(t *testing.T) (*GlobalCollegeRepository, *pgxpool.Pool) {
	t.Helper()
	db, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:54322/postgres")
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	return NewGlobalCollegeRepository(db), db
}

// cleanupCollege deletes a college by ID so repeated test runs don't
// accumulate rows or collide with the unique name/location index.
func cleanupCollege(t *testing.T, db *pgxpool.Pool, id int64) {
	t.Helper()
	t.Cleanup(func() {
		_, err := db.Exec(context.Background(), `DELETE FROM public.global_colleges WHERE id = $1`, id)
		if err != nil {
			t.Logf("cleanup failed for college id=%d: %v", id, err)
		}
	})
}