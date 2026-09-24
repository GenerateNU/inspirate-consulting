package todoitem

import (
	"context"
	"inspirate-consulting/internal/models"
)

func (h *Handler) CreateTodoItem(ctx context.Context, input *models.CreateTodoItemRequestBody) (*models.TodoItem, error) {

	createdTodoItem, err := h.TodoItemRepository.CreateTodoItem(ctx, input)
	if err != nil {
		return nil, err
	}
	return createdTodoItem, nil

}
