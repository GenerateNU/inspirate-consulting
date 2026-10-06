package mediaaccess

import(
	"context"
)

func (h *Handler) RevokeMediaAccess(ctx context.Context, id string)error{
	err := h.MediaAccessRepository.RevokeMediaAccess(ctx, id)
	if(err != nil){
		return err
	}

	return nil
}