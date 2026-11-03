package todoItemRepository

import (
	"context"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
)

func (r *TodoItemRepository) DeleteTodoItem(ctx context.Context, id string) error {
	deleteQuery, err := schema.ReadSQLBaseScript("delete_todo.sql", SqlTodoItemFiles)
	if err != nil {
		return err
	}

	result, err := r.db.Exec(ctx, deleteQuery, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errs.NotFound("todo item", "id", id)
	}

	return nil
}
