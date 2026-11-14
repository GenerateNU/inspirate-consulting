package todoItemRepository

import (
	"context"
	"errors"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"

	"github.com/jackc/pgx/v5"

	"github.com/google/uuid"
)

func (r *TodoItemRepository) GetTodoItemStudentID(ctx context.Context, todoItemID string) (uuid.UUID, error) {
	query, err := schema.ReadSQLBaseScript("get_student_id.sql", SqlTodoItemFiles)
	if err != nil {
		e := errs.InternalServerError("Failed to read base query: ", err.Error())
		return uuid.Nil, &e
	}

	var studentID uuid.UUID
	if err := r.db.QueryRow(ctx, query, todoItemID).Scan(&studentID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, errs.NotFound("todo items", "id", todoItemID)
		}
		return uuid.Nil, err
	}

	return studentID, nil

}
