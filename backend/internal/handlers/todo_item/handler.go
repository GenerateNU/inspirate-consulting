package todoitem

import (
	"inspirate-consulting/internal/auth"
	storage "inspirate-consulting/internal/data"
)

type Handler struct {
	TodoItemRepository storage.TodoItemRepository
	Authorizer         *auth.Authorizer
}

func NewHandler(todoItemRepository storage.TodoItemRepository, authorizer *auth.Authorizer) *Handler {
	return &Handler{
		TodoItemRepository: todoItemRepository,
		Authorizer:         authorizer,
	}
}
