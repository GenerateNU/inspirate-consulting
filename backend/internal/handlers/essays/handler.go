package essays

import (
	"inspirate-consulting/internal/auth"
	storage "inspirate-consulting/internal/data"
)

type Handler struct {
	EssayRepository       storage.EssayRepository
	EssayGroupsRepository storage.EssayGroupRepository
	Authorizer            *auth.Authorizer
}

func NewHandler(essayRepository storage.EssayRepository, authorizer *auth.Authorizer, essayGroupRepository storage.EssayGroupRepository) *Handler {
	return &Handler{
		EssayRepository:       essayRepository,
		EssayGroupsRepository: essayGroupRepository,
		Authorizer:            authorizer,
	}
}
