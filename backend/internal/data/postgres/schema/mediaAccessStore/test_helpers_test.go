package mediaAccessRepository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// createTestMedia inserts a media row directly via SQL (to avoid a cross-package test dependency) and
// returns its id, for use as a valid media_id foreign key in media access tests.
func createTestMedia(t *testing.T, db *pgxpool.Pool, title string) string {
	t.Helper()
	return createTestMediaWithKey(t, db, title, "media/test.mp4")
}

// createTestMediaWithKey is createTestMedia with a specific S3 key, for tests that look media up by key.
func createTestMediaWithKey(t *testing.T, db *pgxpool.Pool, title string, s3Key string) string {
	t.Helper()

	var id string
	err := db.QueryRow(
		context.Background(),
		`INSERT INTO public.media (title, description, length_in_mins, s3_key) VALUES ($1, $2, $3, $4) RETURNING id`,
		title,
		"test description",
		10,
		s3Key,
	).Scan(&id)
	if err != nil {
		t.Fatalf("failed to create test media: %v", err)
	}

	return id
}

// createTestMediaTodo inserts a todo item linked to the media for the student, for due date tests.
func createTestMediaTodo(t *testing.T, db *pgxpool.Pool, studentID string, mediaID string, deadline time.Time, completed bool) {
	t.Helper()

	var completedAt *time.Time
	if completed {
		now := time.Now()
		completedAt = &now
	}

	_, err := db.Exec(
		context.Background(),
		`INSERT INTO public.todo_items (student_id, user_id, title, todo_description, deadline, completed_at, media_id, link_type)
		 VALUES ($1, $2, 'Watch video', 'Watch the assigned video', $3, $4, $5, 'module')`,
		studentID,
		uuid.NewString(),
		deadline,
		completedAt,
		mediaID,
	)
	if err != nil {
		t.Fatalf("failed to create test todo: %v", err)
	}
}
