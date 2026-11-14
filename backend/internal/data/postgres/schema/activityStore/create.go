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
