package todoitem

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

func (h *Handler) GetTodoItemsByStudent(ctx context.Context) ([]models.TodoItem, error) {

	studentID := auth.GetStudentID(ctx)

	todoItems, err := h.TodoItemRepository.GetTodoItemsByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}
	return todoItems, nil
}
