package globalCollegeRepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *GlobalCollegeRepository) CreateGlobalCollege(ctx context.Context, inputGlobalCollege models.CreateGlobalCollegeInput) (*models.CreateGlobalCollegeOutput, error) {
	createdGlobalCollege := models.GlobalCollege{}

	const insertQuery = `
	INSERT INTO public.global_colleges (
		school_name, school_location, ea_deadline, ed_deadline, rd_deadline
	) VALUES (
		$1, $2, $3, $4, $5
	)
	RETURNING id, created_at, updated_at, school_name, school_location, ea_deadline, ed_deadline, rd_deadline
	`

	err := r.db.QueryRow(
		ctx,
		insertQuery,
		inputGlobalCollege.Body.SchoolName,
		inputGlobalCollege.Body.SchoolLocation,
		inputGlobalCollege.Body.EADeadline,
		inputGlobalCollege.Body.EDDeadline,
		inputGlobalCollege.Body.RDDeadline,
	).Scan(
		&createdGlobalCollege.ID,
		&createdGlobalCollege.CreatedAt,
		&createdGlobalCollege.UpdatedAt,
		&createdGlobalCollege.SchoolName,
		&createdGlobalCollege.SchoolLocation,
		&createdGlobalCollege.EADeadline,
		&createdGlobalCollege.EDDeadline,
		&createdGlobalCollege.RDDeadline,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		// to return a conflict error if the school_name and school_location combination already exists
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errs.Conflict("global college", "school_name/school_location", inputGlobalCollege.Body.SchoolName+" / "+inputGlobalCollege.Body.SchoolLocation)
		}
		
		return nil, err
	}

	return &models.CreateGlobalCollegeOutput{Body: createdGlobalCollege}, nil
}
