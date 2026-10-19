package mediaaccess

import (
	"inspirate-consulting/internal/auth"
	storage "inspirate-consulting/internal/data"
)

type Handler struct {
	MediaAccessRepository storage.MediaAccessRepository
	Authorizer            *auth.Authorizer
}

func NewHandler(mediaAccessRepository storage.MediaAccessRepository, authorizer *auth.Authorizer) *Handler {
	return &Handler{
		MediaAccessRepository: mediaAccessRepository,
		Authorizer:            authorizer,
	}
}
