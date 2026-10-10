package essayRepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *EssayRepository) UpdateStatus(ctx context.Context, essayID uuid.UUID, status models.Status) (*models.Essays, error) {
	updateQuery, err := schema.ReadSQLBaseScript("update_essay_status.sql", SqlEssayFiles)
	if err != nil {
		return nil, errs.InternalServerError("Failed to read base query: ", err.Error())
	}

	rows, err := r.db.Query(ctx, updateQuery, essayID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	essay, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[models.Essays])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("essay", "id", essayID.String())
		}
		return nil, err
	}

	return &essay, nil
}

// UpdateEssayGroup moves an essay into a group. A nil essayGroupID removes the
// essay from whatever group it currently belongs to.
func (r *EssayRepository) UpdateEssayGroup(ctx context.Context, essayID uuid.UUID, essayGroupID *uuid.UUID) (*models.Essays, error) {
	updateQuery, err := schema.ReadSQLBaseScript("update_essay_group.sql", SqlEssayFiles)
	if err != nil {
		return nil, errs.InternalServerError("Failed to read base query: ", err.Error())
	}

	// pgx may surface a constraint violation either here or when the rows are
	// collected, so both paths go through the same mapping.
	rows, err := r.db.Query(ctx, updateQuery, essayID, essayGroupID)
	if err != nil {
		return nil, mapUpdateGroupErr(err, essayID, essayGroupID)
	}
	defer rows.Close()

	essay, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[models.Essays])
	if err != nil {
		return nil, mapUpdateGroupErr(err, essayID, essayGroupID)
	}

	return &essay, nil
}

func mapUpdateGroupErr(err error, essayID uuid.UUID, essayGroupID *uuid.UUID) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return errs.NotFound("essay", "id", essayID.String())
	}

	var pgErr *pgconn.PgError

	// essay_group_id references essay_groups(id), so an unknown group is a client error
	if errors.As(err, &pgErr) && pgErr.Code == "23503" && essayGroupID != nil {
		return errs.NotFound("essay group", "id", essayGroupID.String())
	}

	return err
}
