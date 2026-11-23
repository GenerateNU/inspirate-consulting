package todoitem

import (
	"context"
	"inspirate-consulting/internal/models"
)

// Only counselors can update todo item titles
func (h *Handler) UpdateTodoItemTitle(ctx context.Context, id string, title string) (*models.TodoItem, error) {
	studentID, err := h.TodoItemRepository.GetTodoItemStudentID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := h.Authorizer.CanAccessStudent(ctx, studentID); err != nil {
		return nil, err
	}

	return h.TodoItemRepository.UpdateTodoItemTitle(ctx, id, title)
}
