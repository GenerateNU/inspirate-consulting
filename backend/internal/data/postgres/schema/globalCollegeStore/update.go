package globalCollegeRepository

import (
	"context"
	"errors"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// updates a global college's name, location, and deadlines
func (r *GlobalCollegeRepository) UpdateGlobalCollege(ctx context.Context, id int64, input models.UpdateGlobalCollegeRequestBody, nearDuplicateThreshold float64) (*models.GlobalCollege, error) {

	updatedCollege := &models.GlobalCollege{}

	updateQuery, err := schema.ReadSQLBaseScript("update_global_college.sql", SqlGlobalCollegeFiles)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(
		ctx,
		updateQuery,
		id,
		input.SchoolName,
		input.SchoolLocation,
		input.EADeadline,
		input.EDDeadline,
		input.RDDeadline,
	).Scan(
		&updatedCollege.ID,
		&updatedCollege.SchoolName,
		&updatedCollege.SchoolLocation,
		&updatedCollege.EADeadline,
		&updatedCollege.EDDeadline,
		&updatedCollege.RDDeadline,
		&updatedCollege.CreatedAt,
		&updatedCollege.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errs.Conflict("global college", "school_name/school_location", input.SchoolName+" / "+input.SchoolLocation)
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("global college", "id", id)
		}
		return nil, err
	}

	return updatedCollege, nil
}
