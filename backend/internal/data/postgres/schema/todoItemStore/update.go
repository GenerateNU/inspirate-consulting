package todoItemRepository

import (
	"context"
	"time"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

func (r *TodoItemRepository) UpdateTodoItemCompletedAt(ctx context.Context, id string, completedAt *time.Time) (*models.TodoItem, error) {
	updatedItem := &models.TodoItem{}

	updateQuery, err := schema.ReadSQLBaseScript("update_todo_item_completed_at.sql", SqlTodoItemFiles)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(
		ctx,
		updateQuery,
		completedAt,
		id,
	).Scan(
		&updatedItem.ID,
		&updatedItem.CreatedAt,
		&updatedItem.UpdatedAt,
		&updatedItem.StudentID,
		&updatedItem.UserID,
		&updatedItem.TodoDescription,
		&updatedItem.CompletedAt,
		&updatedItem.Deadline,
	)
	if err != nil {
		return nil, err
	}

	return updatedItem, nil
}
