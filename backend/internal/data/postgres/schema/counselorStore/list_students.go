package counselorRepository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

func (r *CounselorRepository) ListStudents(ctx context.Context, counselorID uuid.UUID) ([]models.CounselorStudentRow, error) {
	listQuery, err := schema.ReadSQLBaseScript("list_students_by_counselor.sql", SqlCounselorFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, listQuery, counselorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.CounselorStudentRow])
	if err != nil {
		return nil, err
	}

	return students, nil
}
