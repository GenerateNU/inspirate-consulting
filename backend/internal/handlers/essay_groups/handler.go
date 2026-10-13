package essay_groups

import storage "inspirate-consulting/internal/data"

type Handler struct {
	EssayGroupRepository storage.EssayGroupRepository
}

func NewHandler(essayGroupRepository storage.EssayGroupRepository) *Handler {
	return &Handler{
		EssayGroupRepository: essayGroupRepository,
	}
}
