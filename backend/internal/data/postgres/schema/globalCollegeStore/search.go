package globalCollegeRepository

import (
	"context"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"

	"github.com/jackc/pgx/v5"
)

// matches query against either school_name OR school_location
func (r *GlobalCollegeRepository) SearchGlobalColleges(ctx context.Context, query string, threshold float64, limit int) ([]models.GlobalCollege, error) {
	searchQuery, err := schema.ReadSQLBaseScript("search_global_colleges.sql", SqlGlobalCollegeFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, searchQuery, query, threshold, limit)
	if err != nil {
		return nil, err
	}

	colleges, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.GlobalCollege])
	if err != nil {
		return nil, err
	}

	return colleges, nil
}
