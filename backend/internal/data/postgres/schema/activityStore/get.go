package activityRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// Gets all the activities from all of the given counselor students.
func (r *ActivityRepository) GetActivitiesForCounselor(ctx context.Context, counselorID uuid.UUID) ([]models.Activity, error) {
	selectQuery, err := schema.ReadSQLBaseScript("get_activities_for_counselor.sql", SqlActivityFiles)
	if err != nil {
		return nil, errs.InternalServerError("Failed to read base query: ", err.Error())
	}

	rows, err := r.db.Query(ctx, selectQuery, counselorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activities, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.Activity])
	if err != nil {
		return nil, fmt.Errorf("failed to collect activities: %w", err)
	}

	return activities, nil
}
