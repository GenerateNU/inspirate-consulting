package extracurricularRepository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

// CreateExtracurricular inserts a new extracurricular row into the
// public.extracurriculars table and returns the created record back to the caller.
func (r *ExtracurricularRepository) CreateExtracurricular(ctx context.Context, userID string, input models.CreateExtracurricularInput) (*models.CreateExtracurricularOutput, error) {
	body := input.Body

	startDate, err := models.ParseDate(body.StartDate)
	if err != nil {
		return nil, err
	}

	var endDate *time.Time
	if body.EndDate != nil {
		parsed, err := models.ParseDate(*body.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = &parsed
	}

	insertQuery, err := schema.ReadSQLBaseScript("create_extracurricular.sql", SqlExtracurricularFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(
		ctx,
		insertQuery,
		input.StudentID,
		userID,
		body.Name,
		body.Status,
		body.Type,
		body.Description,
		body.LeadershipRole,
		startDate,
		endDate,
		body.Organization,
	)
	if err != nil {
		return nil, err
	}

	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Extracurricular])
	if err != nil {
		return nil, err
	}

	return &models.CreateExtracurricularOutput{Body: created}, nil
}

// ListExtracurriculars retrieves extracurriculars for a student
func (r *ExtracurricularRepository) ListExtracurriculars(ctx context.Context, studentID string) ([]models.Extracurricular, error) {
	selectQuery, err := schema.ReadSQLBaseScript("list_extracurriculars.sql", SqlExtracurricularFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, selectQuery, studentID)
	if err != nil {
		return nil, err
	}

	// CollectRows closes rows, checks rows.Err(), and returns an empty
	// (non-nil) slice when there are no results, so the API returns [] not null.
	return pgx.CollectRows(rows, pgx.RowToStructByName[models.Extracurricular])
}

// UpdateExtracurricular updates an extracurricular
// Fields left nil in the request keep their existing value via COALESCE.
// The student ID is taken from the existing row, not the request.
// Returns pgx.ErrNoRows if no extracurricular with the given id exists.
func (r *ExtracurricularRepository) UpdateExtracurricular(ctx context.Context, id int64, input models.UpdateExtracurricularInput) (*models.Extracurricular, error) {
	body := input.Body

	var startDate *time.Time
	if body.StartDate != nil {
		parsed, err := models.ParseDate(*body.StartDate)
		if err != nil {
			return nil, err
		}
		startDate = &parsed
	}

	var endDate *time.Time
	if body.EndDate != nil {
		parsed, err := models.ParseDate(*body.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = &parsed
	}

	updateQuery, err := schema.ReadSQLBaseScript("update_extracurricular.sql", SqlExtracurricularFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(
		ctx,
		updateQuery,
		body.Name,
		body.Status,
		body.Type,
		body.Description,
		body.LeadershipRole,
		startDate,
		endDate,
		body.Organization,
		id,
	)
	if err != nil {
		return nil, err
	}

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Extracurricular])
	if err != nil {
		return nil, err
	}

	return &updated, nil
}
