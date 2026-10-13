package essayGroupRepository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *EssayGroupRepository) ListEssayGroups(ctx context.Context, studentID uuid.UUID) ([]models.EssayGroups, error) {
	listQuery, err := schema.ReadSQLBaseScript("list_essay_groups.sql", SqlEssayGroupFiles)
	if err != nil {
		return nil, errs.InternalServerError("Failed to read base query: ", err.Error())
	}

	rows, err := r.db.Query(ctx, listQuery, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	essayGroups, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.EssayGroups])
	if err != nil {
		return nil, err
	}

	return essayGroups, nil
}
