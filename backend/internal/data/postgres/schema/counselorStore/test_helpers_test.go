package counselorstore

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
)

func setupTestRepo(t *testing.T) (*CounselorRepository, *pgxpool.Pool) {
	t.Helper()

	db := testutils.SetupTestDB(t)
	return NewCounselorRepository(db), db
}

// createTestCounselor inserts a user and a counselor row for them.
func createTestCounselor(t *testing.T, db *pgxpool.Pool, name string) {
	t.Helper()
	ctx := context.Background()

	if _, err := db.Exec(ctx,
		`WITH u AS (INSERT INTO public.users (name) VALUES ($1) RETURNING id)
		 INSERT INTO public.counselor (user_id) SELECT id FROM u`, name); err != nil {
		t.Fatalf("failed to create test counselor: %v", err)
	}
}
