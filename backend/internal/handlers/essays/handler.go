package essays

import (
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/data"
	storage "inspirate-consulting/internal/data"
)

type Handler struct {
	EssayRepository       storage.EssayRepository
	EssayGroupsRepository data.EssayGroupRepository
	Authorizer            *auth.Authorizer
}

func NewHandler(essayRepository storage.EssayRepository, authorizer *auth.Authorizer, essayGroupRepository data.EssayGroupRepository) *Handler {
	return &Handler{
		EssayRepository:       essayRepository,
		EssayGroupsRepository: essayGroupRepository,
		Authorizer:            authorizer,
	}
}
