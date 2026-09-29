package mediaRepository

import(
	"context"
	"inspirate-consulting/internal/errs"
)

func(r *MediaRepository) DeleteMedia(ctx context.Context, id string)error{

	const deleteQuery = `
	DELETE FROM public.media
	WHERE id = $1
	`
	result, err := r.db.Exec(ctx, deleteQuery, id)
	if(err != nil){
		return err
	}
	
	if result.RowsAffected()==0{
		return errs.NotFound("media", "id", id)
	}

	return nil
}