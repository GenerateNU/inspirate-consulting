package media

import (
	storage "inspirate-consulting/internal/data"
)

type Handler struct {
	MediaRepository storage.MediaRepository
}

func NewHandler(mediaRepository storage.MediaRepository) *Handler {
	return &Handler{
		MediaRepository: mediaRepository,
	}
}
