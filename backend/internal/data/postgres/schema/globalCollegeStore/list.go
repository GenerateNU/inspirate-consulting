package globalCollegeRepository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/models"
)

func (r *GlobalCollegeRepository) ListGlobalColleges(ctx context.Context) ([]models.GlobalCollege, error) {

	const listQuery = `
	SELECT id, created_at, updated_at, school_name, school_location, ea_deadline, ed_deadline, rd_deadline
	FROM public.global_colleges
	`

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