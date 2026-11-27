package todoitem

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

func (h *Handler) GetTodoItemStats(ctx context.Context, input *models.GetTodoItemStatsInput) (*models.TodoItemStats, error) {
	requested, err := auth.ParseOptionalStudentID(input.StudentID)
	if err != nil {
		return nil, err
	}
	studentID, err := h.Authorizer.ResolveViewableStudentID(ctx, requested)
	if err != nil {
		return nil, err
	}

	return h.TodoItemRepository.GetTodoItemStats(ctx, studentID.String())
}
