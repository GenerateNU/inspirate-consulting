package globalCollegeRepository

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *GlobalCollegeRepository) GetGlobalCollege(ctx context.Context, id int64) (*models.GlobalCollege, error) {

	const selectQuery = `
	SELECT id, created_at, updated_at, school_name, school_location, ea_deadline, ed_deadline, rd_deadline
	FROM public.global_colleges
	WHERE id = $1
	`

	rows, err := r.db.Query(ctx, selectQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
 
	globalCollege, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[models.GlobalCollege])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("global college", "id", strconv.FormatInt(id, 10))
		}
		return nil, err
	}
 
	return &globalCollege, nil
}