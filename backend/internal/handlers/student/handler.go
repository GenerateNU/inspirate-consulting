package student

import (
	"inspirate-consulting/internal/auth"
	storage "inspirate-consulting/internal/data"
)

type Handler struct {
	EssayReviewRepository storage.EssayReviewRepository
	StudentRepository     storage.StudentRepository
  Authorizer            *auth.Authorizer
}

func NewHandler(essayReviewRepository storage.EssayReviewRepository, studentRepository storage.StudentRepository, authorizer *auth.Authorizer) *Handler {
	return &Handler{
		EssayReviewRepository: essayReviewRepository,
		StudentRepository:     studentRepository,
    Authorizer:            authorizer,
	}
}
