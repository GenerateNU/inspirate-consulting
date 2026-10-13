package todoItemRepository

import (
	"context"
	"strings"
	"time"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"

	"github.com/jackc/pgx/v5"
)

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

func (r *TodoItemRepository) ListTodoItems(ctx context.Context, query models.TodoItemQuery) ([]models.TodoItem, error) {
	selectQuery, err := schema.ReadSQLBaseScript("list_todo_items.sql", SqlTodoItemFiles)
	if err != nil {
		return nil, err
	}

	query.ApplyDefaults()

	// Escape LIKE wildcards so the search term matches literally
	var search *string
	if query.Search != nil {
		escaped := likeEscaper.Replace(*query.Search)
		search = &escaped
	}

	var afterValue *time.Time
	var afterID *string
	if query.After != nil {
		afterValue, afterID = &query.After.SortValue, &query.After.ID
	}

	rows, err := r.db.Query(
		ctx,
		selectQuery,
		query.StudentID,
		query.Status,
		query.EssayID,
		query.MediaID,
		query.GlobalCollegeID,
		query.LinkedTo,
		search,
		query.SortBy,
		query.SortOrder,
		query.NullSortValue(),
		afterValue,
		afterID,
		query.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByPos[models.TodoItem])
}
