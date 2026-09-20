package todoItemRepository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"inspirate-consulting/internal/models"
)

func (r *TodoItemRepository) GetTodoItemsByStudent(ctx context.Context, studentID string) ([]models.TodoItem, error) {
	const selectQuery = `
	SELECT id, created_at, updated_at, student_id, user_id, todo_description, completed_at, deadline
	FROM public.todo_items
	WHERE student_id = $1
	`

	rows, err := r.db.Query(ctx, selectQuery, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByPos[models.TodoItem])
}