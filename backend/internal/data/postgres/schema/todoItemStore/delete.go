package todoItemRepository

import (
	"context"
	"errors"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"

	"github.com/jackc/pgx/v5"
)

func (r *TodoItemRepository) DeleteTodoItem(ctx context.Context, id string) error {
	deleteQuery, err := schema.ReadSQLBaseScript("delete_todo.sql", SqlTodoItemFiles)
	if err != nil {
		return err
	}
	revokeQuery, err := schema.ReadSQLBaseScript("revoke_media_access_for_todo.sql", SqlTodoItemFiles)
	if err != nil {
		return err
	}

	// Deleting a todo that linked media also revokes the access it granted
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var studentID string
	var mediaID *string
	err = tx.QueryRow(ctx, deleteQuery, id).Scan(&studentID, &mediaID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.NotFound("todo item", "id", id)
		}
		return err
	}

	if mediaID != nil {
		if _, err := tx.Exec(ctx, revokeQuery, studentID, *mediaID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
