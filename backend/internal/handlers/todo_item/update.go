package todoitem
import (
	"time"
	"context"
	"inspirate-consulting/internal/models"
)

func(h *Handler) UpdateTodoItemCompletedAt(ctx context.Context, id string, completed bool)(*models.TodoItem, error){
	var completedAt *time.Time
	if completed{
		now := time.Now()
		completedAt = &now
	}
	updatedItem, err := h.TodoItemRepository.UpdateTodoItemCompletedAt(ctx, id, completedAt)
	if(err != nil){
		return nil, err
	}
	return updatedItem, nil
}