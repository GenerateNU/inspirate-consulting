package mediaAccessRepository

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// createTestMedia inserts a media row directly via SQL (to avoid a cross-package test dependency) and
// returns its id, for use as a valid media_id foreign key in media access tests.
func createTestMedia(t *testing.T, db *pgxpool.Pool, title string) string {
	t.Helper()

	var id string
	err := db.QueryRow(
		context.Background(),
		`INSERT INTO public.media (title, description, length_in_mins, s3_key) VALUES ($1, $2, $3, $4) RETURNING id`,
		title,
		"test description",
		10,
		"media/test.mp4",
	).Scan(&id)
	if err != nil {
		t.Fatalf("failed to create test media: %v", err)
	}

	return id
}
