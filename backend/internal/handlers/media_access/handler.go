package mediaaccess

import (
	storage "inspirate-consulting/internal/data"
)

type Handler struct{
	MediaAccessRepository storage.MediaAccessRepository
}

func NewHandler(mediaAccessRepository storage.MediaAccessRepository) *Handler{
	return &Handler{
		MediaAccessRepository: mediaAccessRepository,
	}
}