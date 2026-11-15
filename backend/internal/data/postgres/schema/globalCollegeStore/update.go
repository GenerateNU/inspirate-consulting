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

// updates a global college's name and location, and replaces its application cycles, in one transaction
func (r *GlobalCollegeRepository) UpdateGlobalCollege(ctx context.Context, id int64, input models.UpdateGlobalCollegeRequestBody) (*models.GlobalCollege, error) {
	updateQuery, err := schema.ReadSQLBaseScript("update_global_college.sql", SqlGlobalCollegeFiles)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	rows, err := tx.Query(ctx, updateQuery, id, input.SchoolName, input.SchoolLocation)
	if err != nil {
		return nil, err
	}
	row, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[collegeRow])
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolationCode {
			return nil, errs.Conflict("global college", "school_name/school_location", input.SchoolName+" / "+input.SchoolLocation)
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("global college", "id", id)
		}
		return nil, err
	}

	// existing applications under a removed cycle are kept; they come back with cycle_offered=false
	if err := replaceCycles(ctx, tx, id, input.Cycles); err != nil {
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
