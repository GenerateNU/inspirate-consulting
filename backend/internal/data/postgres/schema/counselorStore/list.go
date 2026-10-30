package counselorRepository

import (
	"context"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/pagination"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (c *CounselorRepository) ListCounselors(ctx context.Context, name *string, cursor *pagination.NameIDKey, limit int) ([]models.CounselorListItem, error) {
	selectQuery, err := schema.ReadSQLBaseScript("list_counselors.sql", SqlCounselorFiles)
	if err != nil {
		return nil, err
	}

	var cursorName *string
	var cursorID *uuid.UUID
	if cursor != nil {
		cursorName, cursorID = &cursor.Name, &cursor.ID
	}

	rows, err := c.db.Query(ctx, selectQuery, name, cursorName, cursorID, limit)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[models.CounselorListItem])
}
