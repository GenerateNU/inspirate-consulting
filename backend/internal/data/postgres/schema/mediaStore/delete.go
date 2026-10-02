package mediaRepository

import (
	"context"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
)

func (r *MediaRepository) DeleteMedia(ctx context.Context, id string) error {

	deleteQuery, err := schema.ReadSQLBaseScript("delete_media.sql", SqlMediaFiles)
	if err != nil {
		return err
	}
	result, err := r.db.Exec(ctx, deleteQuery, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errs.NotFound("media", "id", id)
	}

	return nil
}
