package studentRepository

import (
	"context"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/pagination"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *StudentRepository) ListStudents(ctx context.Context, name *string, cursor *pagination.NameIDKey, limit int) ([]models.StudentListItem, error) {
	selectQuery, err := schema.ReadSQLBaseScript("list_students.sql", SqlStudentFiles)
	if err != nil {
		return nil, err
	}

	var cursorName *string
	var cursorID *uuid.UUID
	if cursor != nil {
		cursorName, cursorID = &cursor.Name, &cursor.ID
	}

	rows, err := r.db.Query(ctx, selectQuery, name, cursorName, cursorID, limit)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[models.StudentListItem])
}
