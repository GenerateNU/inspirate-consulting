package todoitem

import (
	storage "inspirate-consulting/internal/data"
)

type Handler struct {
	TodoItemRepository storage.TodoItemRepository
}

// Handler with reference to TodoItemRepository
func NewHandler(todoItemRepository storage.TodoItemRepository) *Handler {
	return &Handler{
		TodoItemRepository: todoItemRepository,
	}
}
