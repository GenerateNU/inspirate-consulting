package globalCollegeRepository

import (
	"context"

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

	globalColleges, err := collectColleges(ctx, r.db, rows)
	if err != nil {
		return nil, err
	}

	return globalColleges, nil
}
