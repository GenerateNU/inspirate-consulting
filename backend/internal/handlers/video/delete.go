package video

import "context"

func (h *Handler) DeleteVideo(ctx context.Context, s3Key string) error {

	err := h.VideoRepository.DeleteVideo(ctx, s3Key)
	if err != nil {
		return err
	}
	return nil
}
