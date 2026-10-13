package user

import (
	storage "inspirate-consulting/internal/data"
)

type Handler struct {
	UserRepository storage.UserRepository
}

func NewHandler(userRepository storage.UserRepository) *Handler {
	return &Handler{
		UserRepository: userRepository,
	}
}
