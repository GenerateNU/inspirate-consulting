package essays

import (
	"inspirate-consulting/internal/auth"
	storage "inspirate-consulting/internal/data"
)

type Handler struct {
	EssayRepository storage.EssayRepository
	Authorizer      *auth.Authorizer
}

func NewHandler(essayRepository storage.EssayRepository, authorizer *auth.Authorizer) *Handler {
	return &Handler{
		EssayRepository: essayRepository,
		Authorizer:      authorizer,
	}
}
