package essay_groups

import (
	"inspirate-consulting/internal/auth"
	storage "inspirate-consulting/internal/data"
)

type Handler struct {
	EssayGroupRepository storage.EssayGroupRepository
	Authorizer           *auth.Authorizer
}

func NewHandler(essayGroupRepository storage.EssayGroupRepository, authorizer *auth.Authorizer) *Handler {
	return &Handler{
		EssayGroupRepository: essayGroupRepository,
		Authorizer:           authorizer,
	}
}
