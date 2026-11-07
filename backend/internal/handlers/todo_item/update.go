package todoitem

import (
	"context"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"time"
)

func (h *Handler) UpdateTodoItemCompletedAt(ctx context.Context, id string, completed bool) (*models.TodoItem, error) {
	var completedAt *time.Time
	if completed {
		now := time.Now()
		completedAt = &now
	}
	updatedItem, err := h.TodoItemRepository.UpdateTodoItemCompletedAt(ctx, id, completedAt)
	if err != nil {
		return nil, err
	}
	return updatedItem, nil
}

func (h *Handler) UpdateTodoItemDescription(ctx context.Context, id string, newDescription string) (*models.TodoItem, error) {
	updatedItem, err := h.TodoItemRepository.UpdateTodoItemDescription(ctx, id, newDescription)
	if err != nil {
		return nil, err
	}

	return updatedItem, nil
}

func (h *Handler) UpdateTodoItemDeadline(ctx context.Context, id string, newDeadline *time.Time) (*models.TodoItem, error) {
	if newDeadline != nil && newDeadline.Before(time.Now()) {
		return nil, errs.BadRequest("deadline cannot be in the past")
	}
	updatedItem, err := h.TodoItemRepository.UpdateTodoItemDeadline(ctx, id, newDeadline)
	if err != nil {
		return nil, err
	}
	return updatedItem, nil
}
