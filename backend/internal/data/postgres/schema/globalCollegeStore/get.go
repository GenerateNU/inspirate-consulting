package globalCollegeRepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *GlobalCollegeRepository) GetGlobalCollege(ctx context.Context, id int64) (*models.GetGlobalCollegeOutput, error) {
	globalCollege := models.GlobalCollege{}

	const selectQuery = `
	SELECT id, created_at, updated_at, school_name, school_location, ea_deadline, ed_deadline, rd_deadline
	FROM public.global_colleges
	WHERE id = $1
	`

	err := r.db.QueryRow(
		ctx,
		selectQuery,
		id,
	).Scan(
		&globalCollege.ID,
		&globalCollege.CreatedAt,
		&globalCollege.UpdatedAt,
		&globalCollege.SchoolName,
		&globalCollege.SchoolLocation,
		&globalCollege.EADeadline,
		&globalCollege.EDDeadline,
		&globalCollege.RDDeadline,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("global college", "id", id)
		}
		return nil, err
	}

	return &models.GetGlobalCollegeOutput{Body: globalCollege}, nil
}