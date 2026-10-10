package essays

import storage "inspirate-consulting/internal/data"

type Handler struct {
	EssayRepository storage.EssayRepository
}

func NewHandler(essayRepository storage.EssayRepository) *Handler {
	return &Handler{
		EssayRepository: essayRepository,
	}
}
