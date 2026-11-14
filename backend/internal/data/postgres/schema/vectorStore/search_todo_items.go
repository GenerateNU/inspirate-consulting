package vectorRepository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

// returns todo items closest in meaning to the query embedding, most similar first.
// returns the student's todo items or the counselor's students' todo items, depending on the studentIDs provided.
func (r *VectorRepository) SearchTodoItems(ctx context.Context, queryEmbedding []float32, studentIDs []uuid.UUID, limit int) ([]models.TodoItemSearchResult, error) {
	searchQuery, err := schema.ReadSQLBaseScript("search_todo_items.sql", SqlVectorFiles)
	if err != nil {
		return nil, err
	}

	// pgx doesn't support passing a slice of UUIDs directly, so convert to a slice of strings.
	ids := make([]string, len(studentIDs))
	for i, id := range studentIDs {
		ids[i] = id.String()
	}

	rows, err := r.db.Query(ctx, searchQuery, formatVector(queryEmbedding), ids, limit)
	if err != nil {
		return nil, err
	}

	similarTodos, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.TodoItemSearchResult])
	if err != nil {
		return nil, err
	}

	return similarTodos, nil


}
