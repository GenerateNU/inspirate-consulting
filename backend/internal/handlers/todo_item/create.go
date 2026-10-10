package todoitem

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) CreateTodoItem(ctx context.Context, input *models.CreateTodoItemRequestBody) (*models.TodoItem, error) {
	input.StudentID = auth.GetStudentID(ctx)
	input.UserID = auth.GetUserID(ctx)

	linked := 0
	if input.EssayID != nil {
		linked++
	}
	if input.MediaID != nil {
		linked++
	}
	if input.GlobalCollegeID != nil {
		linked++
	}
	if linked > 1 {
		return nil, errs.BadRequest("a todo item can be linked to at most one of essay, media, or college")
	}

	createdTodoItem, err := h.TodoItemRepository.CreateTodoItem(ctx, input)
	if err != nil {
		return nil, err
	}
	return createdTodoItem, nil

}
