package essayreview

import (
	"inspirate-consulting/internal/auth"
	storage "inspirate-consulting/internal/data"
)

type Handler struct {
	EssayReviewRepository storage.EssayReviewRepository
	Authorizer            *auth.Authorizer
}

func NewHandler(essayReviewRepository storage.EssayReviewRepository, authorizer *auth.Authorizer) *Handler {
	return &Handler{
		EssayReviewRepository: essayReviewRepository,
		Authorizer:            authorizer,
	}
}
