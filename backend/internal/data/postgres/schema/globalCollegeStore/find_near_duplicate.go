package globalCollegeRepository

import (
	"context"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"

	"github.com/jackc/pgx/v5"
)

// matches query against school_name AND school_location
func (r *GlobalCollegeRepository) FindNearDuplicateGlobalColleges(ctx context.Context, name string, location string, threshold float64) ([]models.GlobalCollege, error) {
	findDuplicatesQuery, err := schema.ReadSQLBaseScript("find_near_duplicate_global_colleges.sql", SqlGlobalCollegeFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, findDuplicatesQuery, name, location, threshold)
	if err != nil {
		return nil, err
	}

	colleges, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.GlobalCollege])
	if err != nil {
		return nil, err
	}

	return colleges, nil
}
