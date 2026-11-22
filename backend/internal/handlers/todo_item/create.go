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
	hasEssay := input.EssayID != nil
	hasMedia := input.MediaID != nil
	hasCollege := input.GlobalCollegeID != nil
	hasExtracurricular := input.ExtracurricularID != nil

	//todo item links to one thing at a time
	linked := 0
	for _, isSet := range []bool{hasEssay, hasMedia, hasCollege, hasExtracurricular} {
		if isSet {
			linked++
		}
	}
	if linked > 1 {
		return nil, errs.BadRequest("a todo item can be linked to at most one of essay, media, college, or extracurricular")
	}

	// If it links to something, link_type must say what
	if err := validateLink(input.LinkType, hasEssay, hasMedia, hasCollege, hasExtracurricular); err != nil {
		return nil, err
	}

	// A todo that links media grants the student access to it, so only counselors may link media.
	// A student linking media would either grant themselves a video or duplicate a video they were already assigned.
	if input.MediaID != nil {
		if _, err := auth.GetStudentID(ctx); err == nil {
			return nil, errs.Forbidden("only counselors can assign videos")
		}
	}

	createdTodoItem, err := h.TodoItemRepository.CreateTodoItem(ctx, studentID.String(), userID.String(), *input)
	if err != nil {
		return nil, err
	}
	return createdTodoItem, nil

}
