package todoitem

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) CreateTodoItem(ctx context.Context, input *models.CreateTodoItemRequestBody) (*models.TodoItem, error) {
	requested, err := auth.ParseOptionalStudentID(input.StudentID)
	if err != nil {
		return nil, err
	}
	studentID, err := h.Authorizer.ResolveStudentID(ctx, requested)
	if err != nil {
		return nil, err
	}
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, errs.Unauthorized()
	}
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

	createdTodoItem, err := h.TodoItemRepository.CreateTodoItem(ctx, studentID.String(), userID.String(), *input)
	if err != nil {
		return nil, err
	}
	return createdTodoItem, nil

}
