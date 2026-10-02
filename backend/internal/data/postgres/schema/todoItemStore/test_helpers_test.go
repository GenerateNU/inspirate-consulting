package todoItemRepository

import (
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// setupTestRepo creates a fresh, isolated test database (via testcontainers)
// and returns a repository connected to it. Safe for t.Parallel().
func setupTestRepo(t *testing.T) (*TodoItemRepository, *pgxpool.Pool) {
	t.Helper()
	db := testutils.SetupTestDB(t)
	return NewTodoItemRepository(db), db
}
