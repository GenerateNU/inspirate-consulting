package mediaRepository

import (
	"context"

	"inspirate-consulting/internal/data/postgres/schema"

	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/models"
)

func (r *MediaRepository) ListAllMedia(ctx context.Context) ([]models.Media, error) {
	listQuery, err := schema.ReadSQLBaseScript("list_media.sql", SqlMediaFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, listQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	allMedia, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.Media])

	if err != nil {
		return nil, err
	}
	return allMedia, nil
}
