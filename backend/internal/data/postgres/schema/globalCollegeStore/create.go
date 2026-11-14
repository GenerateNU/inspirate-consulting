epackage globalCollegeRepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// creates a global college and its application cycles in one transaction
func (r *GlobalCollegeRepository) CreateGlobalCollege(ctx context.Context, inputGlobalCollege models.CreateGlobalCollegeRequestBody) (*models.GlobalCollege, error) {
	insertQuery, err := schema.ReadSQLBaseScript("create_global_college.sql", SqlGlobalCollegeFiles)
	if err != nil {
		return nil, errs.InternalServerError("Failed to read base query: ", err.Error())
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	rows, err := tx.Query(ctx, insertQuery, inputGlobalCollege.SchoolName, inputGlobalCollege.SchoolLocation)
	if err != nil {
		return nil, err
	}
	row, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[collegeRow])
	if err != nil {
		var pgErr *pgconn.PgError

		// to return a conflict error if the school_name and school_location combination already exists (exact case-insensitive match)
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolationCode {
			return nil, errs.Conflict("global college", "school_name/school_location", inputGlobalCollege.SchoolName+" / "+inputGlobalCollege.SchoolLocation)
		}
		return nil, err
	}

	if err := replaceCycles(ctx, tx, row.ID, inputGlobalCollege.Cycles); err != nil {
		return nil, mapCycleError(err)
	}

	colleges := []models.GlobalCollege{row.toModel()}
	if err := attachCycles(ctx, tx, colleges); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &colleges[0], nil
}
