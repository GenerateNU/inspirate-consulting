package globalcollege

import (
	"context"
)

func (h *Handler) DeleteGlobalCollege(ctx context.Context, id int64) error {
	return h.GlobalCollegeRepository.DeleteGlobalCollege(ctx, id)
}
