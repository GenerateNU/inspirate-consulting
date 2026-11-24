package mediaAccessRepository

import (
	"context"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"

	"github.com/jackc/pgx/v5"
)

func (r *MediaAccessRepository) ListAccessibleMedia(ctx context.Context, studentID string, query string, threshold float64, limit int, offset int) ([]models.AccessibleMedia, error) {
	listQuery, err := schema.ReadSQLBaseScript("list_accessible_media.sql", SqlMediaAccessFiles)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, listQuery, studentID, limit, offset, query, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accessibleMedia, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.AccessibleMedia])

	if err != nil {
		return nil, err
	}

	return accessibleMedia, nil
}
