package media

import (
	storage "inspirate-consulting/internal/data"
)

type Handler struct {
	MediaRepository       storage.MediaRepository
	MediaAccessRepository storage.MediaAccessRepository
}

func NewHandler(mediaRepository storage.MediaRepository, mediaAccessRepository storage.MediaAccessRepository) *Handler {
	return &Handler{
		MediaRepository:       mediaRepository,
		MediaAccessRepository: mediaAccessRepository,
	}
}
