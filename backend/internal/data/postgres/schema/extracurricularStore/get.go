package extracurricularRepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// GetExtracurricular returns one extracurricular together with all of its
// achievements and notes. The three reads run in a single read-only,
// repeatable-read transaction so the children match the parent snapshot.
// Returns a 404 HTTPError if no extracurricular with the given id exists.
func (r *ExtracurricularRepository) GetExtracurricular(ctx context.Context, id int64) (*models.ExtracurricularDetail, error) {
	getQuery, err := schema.ReadSQLBaseScript("get_extracurricular.sql", SqlExtracurricularFiles)
	if err != nil {
		return nil, err
	}
	achievementsQuery, err := schema.ReadSQLBaseScript("list_extracurricular_achievements.sql", SqlExtracurricularFiles)
	if err != nil {
		return nil, err
	}
	notesQuery, err := schema.ReadSQLBaseScript("list_extracurricular_notes.sql", SqlExtracurricularFiles)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	rows, err := tx.Query(ctx, getQuery, id)
	if err != nil {
		return nil, err
	}
	extracurricular, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Extracurricular])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errs.NotFound("extracurricular", "id", id)
	}
	if err != nil {
		return nil, err
	}

	rows, err = tx.Query(ctx, achievementsQuery, id)
	if err != nil {
		return nil, err
	}
	achievements, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.ExtracurricularAchievement])
	if err != nil {
		return nil, err
	}

	rows, err = tx.Query(ctx, notesQuery, id)
	if err != nil {
		return nil, err
	}
	notes, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.ExtracurricularNote])
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &models.ExtracurricularDetail{
		Extracurricular: extracurricular,
		Achievements:    achievements,
		Notes:           notes,
	}, nil
}
