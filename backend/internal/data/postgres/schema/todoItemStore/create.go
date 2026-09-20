package todoItemRepository

import (
	"context"
	"inspirate-consulting/internal/models"
)

func (r *TodoItemRepository) CreateTodoItem(ctx context.Context, item *models.TodoItem) (*models.TodoItem, error) {
	createdItem := &models.TodoItem{}

	const insertQuery = `
	INSERT INTO public.todo_items (
		student_id, user_id, todo_description, deadline
	) VALUES (
		$1, $2, $3, $4
	)
	RETURNING id, created_at, updated_at, student_id, user_id, todo_description, completed_at, deadline
	`

	err := r.db.QueryRow(
		ctx,
		insertQuery,
		item.StudentID,
		item.UserID,
		item.TodoDescription,
		item.Deadline,
	).Scan(
		&createdItem.ID,
		&createdItem.CreatedAt,
		&createdItem.UpdatedAt,
		&createdItem.StudentID,
		&createdItem.UserID,
		&createdItem.TodoDescription,
		&createdItem.CompletedAt,
		&createdItem.Deadline,
	)
	if err != nil {
		return nil, err
	}

	return createdItem, nil
}
