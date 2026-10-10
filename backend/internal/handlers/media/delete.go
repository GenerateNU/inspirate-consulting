package media

import (
	"context"
)

func (h *Handler) DeleteMedia(ctx context.Context, id string) error {
	err := h.MediaRepository.DeleteMedia(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
