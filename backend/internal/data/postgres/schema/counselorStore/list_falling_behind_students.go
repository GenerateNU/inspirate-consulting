package counselorRepository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
)

// returns the ids of the counselor's students with an incomplete to-do item past its deadline
func (r *CounselorRepository) ListFallingBehindStudentIDs(ctx context.Context, counselorID string) ([]uuid.UUID, error) {
	listQuery, err := schema.ReadSQLBaseScript("list_falling_behind_students_by_counselor.sql", SqlCounselorFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, listQuery, counselorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	studentIDs, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return nil, err
	}

	return studentIDs, nil
}
