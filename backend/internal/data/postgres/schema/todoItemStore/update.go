package todoItemRepository

import (
	"context"
	"time"

	"inspirate-consulting/internal/models"
)

func (r *TodoItemRepository) UpdateTodoItemCompletedAt(ctx context.Context, id string, completedAt *time.Time) (*models.TodoItem, error) {
	updatedItem := &models.TodoItem{}

	const updateQuery = `
	UPDATE public.todo_items
	SET completed_at = $1, updated_at = now()
	WHERE id = $2
	RETURNING id, created_at, updated_at, student_id, user_id, todo_description, completed_at, deadline
	`

	err := r.db.QueryRow(
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