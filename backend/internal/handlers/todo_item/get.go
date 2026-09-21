package todoitem

import (
	"context"
	"inspirate-consulting/internal/models"
)

func(h *Handler) GetTodoItemsByStudent(ctx context.Context, studentID string)([]models.TodoItem, error){

	todoItems, err := h.TodoItemRepository.GetTodoItemsByStudent(ctx, studentID)
	if(err != nil){
		return nil, err
	}
	return todoItems, nil
}