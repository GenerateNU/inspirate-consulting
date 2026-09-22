package todoItemRepository

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// setupTestRepo connects to the local Supabase Postgres instance.
// Only called from tests that are skipped in short mode
func setupTestRepo(t *testing.T) (*TodoItemRepository, *pgxpool.Pool) {
	t.Helper()
	db, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:54322/postgres")
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	return NewTodoItemRepository(db), db
}

// cleanupTodoItem deletes a todo item by ID so repeated test runs don't
// accumulate rows or interfere with per-student queries.
func cleanupTodoItem(t *testing.T, db *pgxpool.Pool, id string) {
	t.Helper()
	t.Cleanup(func() {
		_, err := db.Exec(context.Background(), `DELETE FROM public.todo_items WHERE id = $1`, id)
		if err != nil {
			t.Logf("cleanup failed for todo item id=%s: %v", id, err)
		}
	})
}
