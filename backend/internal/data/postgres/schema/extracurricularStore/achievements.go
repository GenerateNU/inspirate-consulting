package extracurricularRepository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// CreateExtracurricularAchievement adds an achievement to an extracurricular and returns it.
func (r *ExtracurricularRepository) CreateExtracurricularAchievement(ctx context.Context, extracurricularID int64, input models.CreateExtracurricularAchievementRequest) (*models.ExtracurricularAchievement, error) {
	entryDate, err := models.ParseDate(input.EntryDate)
	if err != nil {
		return nil, err
	}

	query, err := schema.ReadSQLBaseScript("create_extracurricular_achievement.sql", SqlExtracurricularFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, extracurricularID, entryDate, input.Body)
	if err != nil {
		return nil, err
	}
	achievement, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.ExtracurricularAchievement])
	if err != nil {
		return nil, err
	}
	return &achievement, nil
}

// UpdateExtracurricularAchievement applies a partial update to an achievement
// belonging to the given extracurricular. Omitted fields keep their value.
// Returns a 404 HTTPError if no such achievement exists under it.
func (r *ExtracurricularRepository) UpdateExtracurricularAchievement(ctx context.Context, extracurricularID, achievementID int64, input models.UpdateExtracurricularAchievementRequest) (*models.ExtracurricularAchievement, error) {
	var entryDate *time.Time
	if input.EntryDate != nil {
		parsed, err := models.ParseDate(*input.EntryDate)
		if err != nil {
			return nil, err
		}
		entryDate = &parsed
	}

	query, err := schema.ReadSQLBaseScript("update_extracurricular_achievement.sql", SqlExtracurricularFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, extracurricularID, achievementID, entryDate, input.Body)
	if err != nil {
		return nil, err
	}
	achievement, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.ExtracurricularAchievement])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errs.NotFound("achievement", "id", achievementID)
	}
	if err != nil {
		return nil, err
	}
	return &achievement, nil
}

// DeleteExtracurricularAchievement removes an achievement belonging to the
// given extracurricular. Returns a 404 HTTPError if no such achievement exists under it.
func (r *ExtracurricularRepository) DeleteExtracurricularAchievement(ctx context.Context, extracurricularID, achievementID int64) error {
	query, err := schema.ReadSQLBaseScript("delete_extracurricular_achievement.sql", SqlExtracurricularFiles)
	if err != nil {
		return err
	}

	result, err := r.db.Exec(ctx, query, extracurricularID, achievementID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errs.NotFound("achievement", "id", achievementID)
	}
	return nil
}
