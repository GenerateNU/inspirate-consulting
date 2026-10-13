package todoItemRepository

import (
	"context"
	"errors"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
)

func (r *TodoItemRepository) CreateTodoItem(ctx context.Context, item *models.CreateTodoItemRequestBody) (*models.TodoItem, error) {
	createdItem := &models.TodoItem{}

	insertQuery, err := schema.ReadSQLBaseScript("create_todo_item.sql", SqlTodoItemFiles)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(
		ctx,
		insertQuery,
		item.StudentID,
		item.UserID,
		item.TodoDescription,
		item.Deadline,
		item.EssayID,
		item.MediaID,
		item.GlobalCollegeID,
	).Scan(
		&createdItem.ID,
		&createdItem.CreatedAt,
		&createdItem.UpdatedAt,
		&createdItem.StudentID,
		&createdItem.UserID,
		&createdItem.TodoDescription,
		&createdItem.CompletedAt,
		&createdItem.Deadline,
		&createdItem.EssayID,
		&createdItem.MediaID,
		&createdItem.GlobalCollegeID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return nil, errs.BadRequest("linked essay, media, or college does not exist")
			case "23514":
				return nil, errs.BadRequest("a todo item can be linked to at most one of essay, media, or college")
			}
		}
		return nil, err
	}

	return createdItem, nil
}
