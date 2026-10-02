package globalCollegeRepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *GlobalCollegeRepository) CreateGlobalCollege(ctx context.Context, inputGlobalCollege models.CreateGlobalCollegeRequestBody) (*models.GlobalCollege, error) {
	createdGlobalCollege := &models.GlobalCollege{}

	insertQuery, err := schema.ReadSQLBaseScript("create_global_college.sql", SqlGlobalCollegeFiles)
	if err != nil {
		fmt.Println("[db] ReadSQLBaseScript error:", err)
		err := errs.InternalServerError("Failed to read base query: ", err.Error())
		return nil, err
	}

	err = r.db.QueryRow(
		ctx,
		insertQuery,
		inputGlobalCollege.SchoolName,
		inputGlobalCollege.SchoolLocation,
		inputGlobalCollege.EADeadline,
		inputGlobalCollege.EDDeadline,
		inputGlobalCollege.RDDeadline,
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
			return nil, errs.Conflict("global college", "school_name/school_location", inputGlobalCollege.SchoolName+" / "+inputGlobalCollege.SchoolLocation)
		}

		return nil, err
	}

	return createdGlobalCollege, nil
}
