package counselorstore

import (
	"context"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/pagination"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (c *CounselorRepository) ListCounselors(ctx context.Context, name *string, after *pagination.NameIDKey, limit int) ([]models.CounselorListItem, error) {
	selectQuery, err := schema.ReadSQLBaseScript("list_counselors.sql", SqlCounselorFiles)
	if err != nil {
		return nil, err
	}

	var afterName *string
	var afterID *uuid.UUID
	if after != nil {
		afterName, afterID = &after.Name, &after.ID
	}

	rows, err := c.db.Query(ctx, selectQuery, name, afterName, afterID, limit)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[models.CounselorListItem])
}
