package globalCollegeRepository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

func (r *GlobalCollegeRepository) ListGlobalColleges(ctx context.Context) ([]models.GlobalCollege, error) {

	listQuery, err := schema.ReadSQLBaseScript("list_global_college.sql", SqlGlobalCollegeFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, listQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	globalColleges, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.GlobalCollege])
	if err != nil {
		return nil, err
	}

	return globalColleges, nil
}
