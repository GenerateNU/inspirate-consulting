package essayRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *EssayRepository) GetEssaysFromStudent(ctx context.Context, studentID uuid.UUID) ([]models.Essays, error) {
	selectQuery, err := schema.ReadSQLBaseScript("get_essays_from_student.sql", SqlEssayFiles)
	if err != nil {
		return nil, errs.InternalServerError("Failed to read base query: ", err.Error())
	}

	rows, err := r.db.Query(ctx, selectQuery, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	essays, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.Essays])

	if err != nil {
		return nil, fmt.Errorf("failed to collect essays: %w", err)
	}

	return essays, nil
}

func (r *EssayRepository) GetEssaysByGroup(ctx context.Context, essayGroupID uuid.UUID) ([]models.Essays, error) {
	selectQuery, err := schema.ReadSQLBaseScript("get_essays_by_group.sql", SqlEssayFiles)
	if err != nil {
		return nil, errs.InternalServerError("Failed to read base query: ", err.Error())
	}

	rows, err := r.db.Query(ctx, selectQuery, essayGroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	essays, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.Essays])

	if err != nil {
		return nil, fmt.Errorf("failed to collect essays: %w", err)
	}

	return essays, nil
}
