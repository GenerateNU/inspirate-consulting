package mediaAccessRepository

import (
	"context"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
)

func (r *MediaAccessRepository) RevokeMediaAccess(ctx context.Context, id string) error {

	deleteQuery, err := schema.ReadSQLBaseScript("revoke_media_access.sql", SqlMediaAccessFiles)
	if err != nil {
		return err
	}
	result, err := r.db.Exec(ctx, deleteQuery, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errs.NotFound("media_access", "id", id)
	}

	return nil
}
