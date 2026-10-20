package studentRepository

import (
	"context"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/pagination"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *StudentRepository) ListStudents(ctx context.Context, name *string, after *pagination.NameIDKey, limit int) ([]models.StudentListItem, error) {
	selectQuery, err := schema.ReadSQLBaseScript("list_students.sql", SqlStudentFiles)
	if err != nil {
		return nil, err
	}

	var afterName *string
	var afterID *uuid.UUID
	if after != nil {
		afterName, afterID = &after.Name, &after.ID
	}

	rows, err := r.db.Query(ctx, selectQuery, name, afterName, afterID, limit)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[models.StudentListItem])
}
