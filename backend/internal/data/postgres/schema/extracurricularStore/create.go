package extracurricularRepository

import (
	"context"
	"fmt"
	"time"

	"inspirate-consulting/internal/models"
)

// CreateExtracurricular inserts a new extracurricular row into the
// public.extracurriculars table and returns the created record back to the caller.
func (r *ExtracurricularRepository) CreateExtracurricular(ctx context.Context, extracurricular models.CreateExtracurricularInput) (*models.CreateExtracurricularOutput, error) {
	createdExtracurricular := &models.CreateExtracurricularOutput{}

	studentID := extracurricular.Body.StudentID
	if studentID == nil || *studentID == "" {
		return nil, fmt.Errorf("student id is required")
	}

	startDate, err := models.ParseDate(extracurricular.Body.StartDate)
	if err != nil {
		return nil, err
	}

	var endDate *time.Time
	if extracurricular.Body.EndDate != nil {
		parsedEndDate, err := models.ParseDate(*extracurricular.Body.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = &parsedEndDate
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
	RETURNING
		id,
		created_at,
		updated_at,
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
	`

	var (
		id           int64
		createdAt    time.Time
		updatedAt    time.Time
		studentIDRow string
		userID       string
		name         string
		status       string
		extracType   string
		description  string
		leadership   *string
		startDateRow time.Time
		endDateRow   *time.Time
		organization string
	)

	err = r.db.QueryRow(
		ctx,
		insertQuery,
		*studentID,
		extracurricular.Body.UserID,
		extracurricular.Body.Name,
		extracurricular.Body.Status,
		extracurricular.Body.Type,
		extracurricular.Body.Description,
		extracurricular.Body.LeadershipRole,
		startDate,
		endDate,
		extracurricular.Body.Organization,
	).Scan(
		&id,
		&createdAt,
		&updatedAt,
		&studentIDRow,
		&userID,
		&name,
		&status,
		&extracType,
		&description,
		&leadership,
		&startDateRow,
		&endDateRow,
		&organization,
	)
	if err != nil {
		return nil, err
	}

	createdExtracurricular.Body = models.Extracurricular{
		ID:             id,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		StudentID:      studentIDRow,
		UserID:         userID,
		Name:           name,
		Status:         models.ExtracurricularStatus(status),
		Type:           models.ExtracurricularType(extracType),
		Description:    description,
		LeadershipRole: leadership,
		StartDate:      startDateRow,
		EndDate:        endDateRow,
		Organization:   organization,
	}

	return createdExtracurricular, nil
}

// ListExtracurriculars retrieves extracurriculars for a student
func (r *ExtracurricularRepository) ListExtracurriculars(ctx context.Context, studentID string) ([]models.Extracurricular, error) {
	const selectQuery = `
	SELECT
		id,
		created_at,
		updated_at,
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
	FROM public.extracurriculars
	WHERE student_id = $1
	ORDER BY start_date DESC
	`

	rows, err := r.db.Query(ctx, selectQuery, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Extracurricular

	for rows.Next() {
		var (
			id           int64
			createdAt    time.Time
			updatedAt    time.Time
			studentIDRow string
			userID       string
			name         string
			status       string
			extracType   string
			description  string
			leadership   *string
			startDate    time.Time
			endDate      *time.Time
			organization string
		)

		if err := rows.Scan(
			&id,
			&createdAt,
			&updatedAt,
			&studentIDRow,
			&userID,
			&name,
			&status,
			&extracType,
			&description,
			&leadership,
			&startDate,
			&endDate,
			&organization,
		); err != nil {
			return nil, err
		}

		out = append(out, models.Extracurricular{
			ID:             id,
			CreatedAt:      createdAt,
			UpdatedAt:      updatedAt,
			StudentID:      studentIDRow,
			UserID:         userID,
			Name:           name,
			Status:         models.ExtracurricularStatus(status),
			Type:           models.ExtracurricularType(extracType),
			Description:    description,
			LeadershipRole: leadership,
			StartDate:      startDate,
			EndDate:        endDate,
			Organization:   organization,
		})
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return out, nil
}

// UpdateExtracurricular updates an extracurricular
func (r *ExtracurricularRepository) UpdateExtracurricular(ctx context.Context, id int64, extracurricular models.UpdateExtracurricularInput) (*models.Extracurricular, error) {
	studentID := extracurricular.Body.StudentID
	if studentID == nil || *studentID == "" {
		return nil, fmt.Errorf("student id is required")
	}

	name := extracurricular.Body.Name
	status := extracurricular.Body.Status
	extracType := extracurricular.Body.Type
	description := extracurricular.Body.Description
	leadership := extracurricular.Body.LeadershipRole
	organization := extracurricular.Body.Organization
	startDate := extracurricular.Body.StartDate
	endDate := extracurricular.Body.EndDate

	if name == nil {
		name = new(string)
	}
	if status == nil {
		status = new(models.ExtracurricularStatus)
	}
	if extracType == nil {
		extracType = new(models.ExtracurricularType)
	}
	if description == nil {
		description = new(string)
	}
	if leadership == nil {
		leadership = new(string)
	}
	if organization == nil {
		organization = new(string)
	}

	var parsedStartDate *time.Time
	if startDate != nil {
		parsed, err := models.ParseDate(*startDate)
		if err != nil {
			return nil, err
		}
		parsedStartDate = &parsed
	}

	var parsedEndDate *time.Time
	if endDate != nil {
		parsed, err := models.ParseDate(*endDate)
		if err != nil {
			return nil, err
		}
		parsedEndDate = &parsed
	}

	const updateQuery = `
	UPDATE public.extracurriculars SET
		name = $1,
		status = $2,
		type = $3,
		description = $4,
		leadership_role = $5,
		start_date = $6,
		end_date = $7,
		organization = $8,
		updated_at = now()
	WHERE id = $9 AND student_id = $10
	RETURNING
		id,
		created_at,
		updated_at,
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
	`

	var (
		out        models.Extracurricular
		leadershipRow *string
		endDateRow *time.Time
		extracTypeRow string
		statusRow string
	)

	err := r.db.QueryRow(
		ctx,
		updateQuery,
		*name,
		*status,
		*extracType,
		*description,
		leadership,
		parsedStartDate,
		parsedEndDate,
		*organization,
		id,
		*studentID,
	).Scan(
		&out.ID,
		&out.CreatedAt,
		&out.UpdatedAt,
		&out.StudentID,
		&out.UserID,
		&out.Name,
		&statusRow,
		&extracTypeRow,
		&out.Description,
		&leadershipRow,
		&out.StartDate,
		&endDateRow,
		&out.Organization,
	)
	if err != nil {
		return nil, err
	}

	out.Status = models.ExtracurricularStatus(statusRow)
	out.Type = models.ExtracurricularType(extracTypeRow)
	out.LeadershipRole = leadershipRow
	out.EndDate = endDateRow

	return &out, nil
}
