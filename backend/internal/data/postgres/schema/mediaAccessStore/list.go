package mediaAccessRepository

import (
	"context"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"

	"github.com/jackc/pgx/v5"
)

func (r *MediaAccessRepository) ListAccessibleMedia(ctx context.Context, studentID string) ([]models.Media, error) {
	listQuery, err := schema.ReadSQLBaseScript("list_accessible_media.sql", SqlMediaAccessFiles)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, listQuery, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accessibleMedia, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.Media])

	if err != nil {
		return nil, err
	}

	return accessibleMedia, nil
}
