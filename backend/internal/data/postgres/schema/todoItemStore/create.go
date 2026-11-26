package todoItemRepository

import (
	"context"
	"errors"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
)

func (r *TodoItemRepository) CreateTodoItem(ctx context.Context, studentID string, userID string, item models.CreateTodoItemRequestBody) (*models.TodoItem, error) {
	createdItem := &models.TodoItem{}

	insertQuery, err := schema.ReadSQLBaseScript("create_todo_item.sql", SqlTodoItemFiles)
	if err != nil {
		return nil, err
	}
	grantQuery, err := schema.ReadSQLBaseScript("grant_media_access_for_todo.sql", SqlTodoItemFiles)
	if err != nil {
		return nil, err
	}

	// A todo that links media also gives the student access to it
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	err = tx.QueryRow(
		ctx,
		insertQuery,
		studentID,
		userID,
		item.TodoDescription,
		item.Deadline,
		item.EssayID,
		item.MediaID,
		item.GlobalCollegeID,
		item.Title,
		item.ExtracurricularID,
		item.LinkType,
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
		&createdItem.Title,
		&createdItem.ExtracurricularID,
		&createdItem.LinkType,
		&createdItem.AssignmentID,
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

	if createdItem.MediaID != nil {
		if _, err := tx.Exec(ctx, grantQuery, createdItem.StudentID, *createdItem.MediaID); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return nil, errs.Conflict("media access", "media_id", *createdItem.MediaID)
			}
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return createdItem, nil
}
