package mediaaccess

import(
	"context"
	"inspirate-consulting/internal/models"
)

func(h *Handler) GrantMediaAccess(ctx context.Context, body *models.GrantMediaAccessRequestBody)(*models.MediaAccess,error){
	createdAccess, err := h.MediaAccessRepository.GrantMediaAccess(ctx, body)
	if(err != nil){
		return nil, err
	}

	return createdAccess, nil
}
