package todoItemRepository

import (
	"context"
	"errors"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/jackc/pgx/v5"
)

func (r *TodoItemRepository) GetTodoItemsByStudent(ctx context.Context, studentID string) ([]models.TodoItem, error) {
	selectQuery, err := schema.ReadSQLBaseScript("list_todo_items_by_student.sql", SqlTodoItemFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, selectQuery, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.TodoItem])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("todo_item", "student_id", studentID)
		}
		return nil, err
	}

	return items, nil
}
