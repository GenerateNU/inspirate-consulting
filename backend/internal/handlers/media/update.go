package media

import (
	"context"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) UpdateMedia(ctx context.Context, id string, body *models.UpdateMediaInputBody) (*models.Media, error) {
	if body.Title == nil && body.Description == nil && body.S3Key == nil && body.YoutubeURL == nil {
		return nil, errs.BadRequest("provide at least one of title, description, s3_key, or youtube_url")
	}
	if body.S3Key != nil && body.YoutubeURL != nil {
		return nil, errs.BadRequest("provide at most one of s3_key or youtube_url")
	}

	fields := models.UpdateMediaFields{
		Title:       body.Title,
		Description: body.Description,
		S3Key:       body.S3Key,
	}
	if body.YoutubeURL != nil {
		youtubeVideoID, err := parseYoutubeVideoID(*body.YoutubeURL)
		if err != nil {
			return nil, err
		}
		fields.YoutubeVideoID = &youtubeVideoID
	}

	updatedMedia, err := h.MediaRepository.UpdateMedia(ctx, id, fields)
	if err != nil {
		return nil, err
	}
	return updatedMedia, nil
}
