package activityRepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *ActivityRepository) CreateActivity(ctx context.Context, studentID uuid.UUID, body models.CreateActivityBody) error {
	if err := validateActivityBody(body); err != nil {
		return err
	}

	insertQuery, err := schema.ReadSQLBaseScript("create_activity.sql", SqlActivityFiles)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, insertQuery,
		studentID,
		body.Action,
		body.EssayID,
		body.EssayStatus,
		body.GlobalCollegeID,
		body.ExtracurricularID,
		body.TodoItemID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return errs.BadRequest("student, essay, college, extracurricular, or todo item does not exist")
		}
		return err
	}

	return nil
}

// validateActivityBody enforces that an activity is about exactly one of
// Essay, GlobalCollege, ExtraCurricular, TodoITemID
func validateActivityBody(body models.CreateActivityBody) error {
	if (body.EssayID == nil) != (body.EssayStatus == nil) {
		return errs.BadRequest("essay_id and essay_status must be set together")
	}

	slots := 0
	if body.EssayID != nil {
		slots++
	}
	if body.GlobalCollegeID != nil {
		slots++
	}
	if body.ExtracurricularID != nil {
		slots++
	}
	if body.TodoItemID != nil {
		slots++
	}

	if slots != 1 {
		return errs.BadRequest("an activity must reference exactly one of: essay, global college, extracurricular, or todo item")
	}

	return nil
}
