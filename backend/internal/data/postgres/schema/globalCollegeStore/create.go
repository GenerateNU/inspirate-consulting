package globalCollegeRepository

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *GlobalCollegeRepository) CreateGlobalCollege(ctx context.Context, inputGlobalCollege models.CreateGlobalCollegeRequestBody, nearDuplicateThreshold float64) (*models.GlobalCollege, error) {

	// near-duplicate check: skipped if the caller has already confirmed they want to proceed despite a prior warning
	if !inputGlobalCollege.ConfirmDuplicate {
		similarColleges, err := r.FindNearDuplicateGlobalColleges(ctx, inputGlobalCollege.SchoolName, inputGlobalCollege.SchoolLocation, nearDuplicateThreshold)
		if err != nil {
			return nil, err
		}
		if len(similarColleges) > 0 {
			return nil, errs.HTTPError{
				Code: http.StatusConflict,
				Message: fmt.Sprintf(
					"%q is very similar to an existing college: %q. Resubmit with confirm_duplicate=true to create it anyway.",
					inputGlobalCollege.SchoolName, similarColleges[0].SchoolName,
				),
			}
		}
	}

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

		// to return a conflict error if the school_name and school_location combination already exists (exact case-insensitive match)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errs.Conflict("global college", "school_name/school_location", inputGlobalCollege.SchoolName+" / "+inputGlobalCollege.SchoolLocation)
		}

		return nil, err
	}

	return createdGlobalCollege, nil
}
