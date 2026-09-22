package student

import storage "inspirate-consulting/internal/data"

// Holds the essay review repository because setting a balance is a ledger write.
type Handler struct {
	EssayReviewRepository storage.EssayReviewRepository
}

func NewHandler(essayReviewRepository storage.EssayReviewRepository) *Handler {
	return &Handler{
		EssayReviewRepository: essayReviewRepository,
	}
}
