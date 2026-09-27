package todoitem

import (
	"context"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/auth"
)

func (h *Handler) CreateTodoItem(ctx context.Context, input *models.CreateTodoItemRequestBody) (*models.TodoItem, error) {
	input.StudentID = auth.GetStudentID(ctx)
	input.UserID = auth.GetUserID(ctx)

	createdTodoItem, err := h.TodoItemRepository.CreateTodoItem(ctx, input)
	if err != nil {
		return nil, err
	}
	return createdTodoItem, nil

}
