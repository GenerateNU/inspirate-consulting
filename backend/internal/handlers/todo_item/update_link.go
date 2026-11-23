package todoitem

import (
	"context"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// UpdateTodoItemLink replaces what a todo item links to. Only counselors can do this, and only
// for todos belonging to a student assigned to them. The route's role guard limits it to
// counselors, and CanAccessStudent checks the student is assigned to the caller.
func (h *Handler) UpdateTodoItemLink(ctx context.Context, id string, link models.TodoItemLink) (*models.TodoItem, error) {
	hasEssay := link.EssayID != nil
	hasMedia := link.MediaID != nil
	hasCollege := link.GlobalCollegeID != nil
	hasExtracurricular := link.ExtracurricularID != nil

	linked := 0
	for _, isSet := range []bool{hasEssay, hasMedia, hasCollege, hasExtracurricular} {
		if isSet {
			linked++
		}
	}
	if linked > 1 {
		return nil, errs.BadRequest("a todo item can be linked to at most one of essay, media, college, or extracurricular")
	}
	if err := validateLink(link.LinkType, hasEssay, hasMedia, hasCollege, hasExtracurricular); err != nil {
		return nil, err
	}

	studentID, err := h.TodoItemRepository.GetTodoItemStudentID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := h.Authorizer.CanAccessStudent(ctx, studentID); err != nil {
		return nil, err
	}

	return h.TodoItemRepository.UpdateTodoItemLink(ctx, id, link)
}
