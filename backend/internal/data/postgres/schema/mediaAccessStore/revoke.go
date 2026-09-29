package mediaAccessRepository

import(
	"context"
	"inspirate-consulting/internal/errs"
)

func(r *MediaAccessRepository) RevokeMediaAccess(ctx context.Context, id string)error {

	const deleteQuery = `
	DELETE FROM public.media_access
	WHERE id = $1
	`
	result, err := r.db.Exec(ctx, deleteQuery, id)
	if(err != nil){
		return err
	}

	if result.RowsAffected() == 0 {
		return errs.NotFound("media_access", "id", id)
	}

	return nil
}