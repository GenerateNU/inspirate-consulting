package extracurricularRepository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/models"
)

// extracurricularColumns is the column list shared by every query that returns
// a full extracurricular row. Order doesn't matter for RowToStructByName, but
// every column must match a db tag on models.Extracurricular.
const extracurricularColumns = `
		id,
		created_at,
		modified_at,
		student_id,
		user_id,
		name,
		status,
		type,
		description,
		leadership_role,
		start_date,
		end_date,
		organization`

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

	const insertQuery = `
	INSERT INTO public.extracurriculars (
		student_id,
		user_id,
		name,
		status,
		type,
		description,
		leadership_role,
		start_date,
		end_date,
		organization
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
	)
	RETURNING` + extracurricularColumns

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
	const selectQuery = `
	SELECT` + extracurricularColumns + `
	FROM public.extracurriculars
	WHERE student_id = $1
	ORDER BY start_date DESC
	`

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

	const updateQuery = `
	UPDATE public.extracurriculars SET
		name            = COALESCE($1, name),
		status          = COALESCE($2, status),
		type            = COALESCE($3, type),
		description     = COALESCE($4, description),
		leadership_role = COALESCE($5, leadership_role),
		start_date      = COALESCE($6, start_date),
		end_date        = COALESCE($7, end_date),
		organization    = COALESCE($8, organization),
		modified_at      = now()
	WHERE id = $9
	RETURNING` + extracurricularColumns

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