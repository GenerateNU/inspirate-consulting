package extracurricularRepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// CreateExtracurricularNote adds a note to an extracurricular and returns it.
func (r *ExtracurricularRepository) CreateExtracurricularNote(ctx context.Context, extracurricularID int64, input models.ExtracurricularNoteRequest) (*models.ExtracurricularNote, error) {
	query, err := schema.ReadSQLBaseScript("create_extracurricular_note.sql", SqlExtracurricularFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, extracurricularID, input.Body)
	if err != nil {
		return nil, err
	}
	note, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.ExtracurricularNote])
	if err != nil {
		return nil, err
	}
	return &note, nil
}

// UpdateExtracurricularNote replaces the body of a note belonging to the given
// extracurricular. Returns a 404 HTTPError if no such note exists under it.
func (r *ExtracurricularRepository) UpdateExtracurricularNote(ctx context.Context, extracurricularID, noteID int64, input models.ExtracurricularNoteRequest) (*models.ExtracurricularNote, error) {
	query, err := schema.ReadSQLBaseScript("update_extracurricular_note.sql", SqlExtracurricularFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, extracurricularID, noteID, input.Body)
	if err != nil {
		return nil, err
	}
	note, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.ExtracurricularNote])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errs.NotFound("note", "id", noteID)
	}
	if err != nil {
		return nil, err
	}
	return &note, nil
}

// DeleteExtracurricularNote removes a note belonging to the given
// extracurricular. Returns a 404 HTTPError if no such note exists under it.
func (r *ExtracurricularRepository) DeleteExtracurricularNote(ctx context.Context, extracurricularID, noteID int64) error {
	query, err := schema.ReadSQLBaseScript("delete_extracurricular_note.sql", SqlExtracurricularFiles)
	if err != nil {
		return err
	}

	result, err := r.db.Exec(ctx, query, extracurricularID, noteID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errs.NotFound("note", "id", noteID)
	}
	return nil
}
