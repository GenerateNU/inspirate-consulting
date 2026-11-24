package todoItemRepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// UpdateTodoItemLink replaces a todo item's link: its link type and all four link ids.
func (r *TodoItemRepository) UpdateTodoItemLink(ctx context.Context, id string, link models.TodoItemLink) (*models.TodoItem, error) {
	query, err := schema.ReadSQLBaseScript("update_todo_item_link.sql", SqlTodoItemFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query,
		id, link.LinkType, link.EssayID, link.MediaID, link.GlobalCollegeID, link.ExtracurricularID)
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
