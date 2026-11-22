package media

import (
	"context"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) CreateMedia(ctx context.Context, item *models.CreateMediaRequestBody) (*models.Media, error) {
	hasS3Key := item.S3Key != nil
	hasYoutubeURL := item.YoutubeURL != nil
	if hasS3Key == hasYoutubeURL {
		return nil, errs.BadRequest("provide exactly one of s3_key or youtube_url")
	}

	var youtubeVideoID *string
	if hasYoutubeURL {
		id, err := parseYoutubeVideoID(*item.YoutubeURL)
		if err != nil {
			return nil, err
		}
		youtubeVideoID = &id
	}

	createdMedia, err := h.MediaRepository.CreateMedia(ctx, item, youtubeVideoID)
	if err != nil {
		return nil, err
	}
	return createdMedia, nil
}
