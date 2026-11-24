package todoItemRepository

import (
	"context"
	"errors"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/jackc/pgx/v5"
)

func (r *TodoItemRepository) UpdateTodoItemTitle(ctx context.Context, id string, title string) (*models.TodoItem, error) {
	query, err := schema.ReadSQLBaseScript("update_todo_item_title.sql", SqlTodoItemFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, id, title)
	if err != nil {
		return nil, err
	}
	item, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByPos[models.TodoItem])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("todo item", "id", id)
		}
		return nil, err
	}
	return item, nil
}
