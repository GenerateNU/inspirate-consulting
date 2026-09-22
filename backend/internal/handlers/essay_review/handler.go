package essayreview

import storage "inspirate-consulting/internal/data"

type Handler struct {
	EssayReviewRepository storage.EssayReviewRepository
}

func NewHandler(essayReviewRepository storage.EssayReviewRepository) *Handler {
	return &Handler{
		EssayReviewRepository: essayReviewRepository,
	}
}
