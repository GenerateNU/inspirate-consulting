package todoitem

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) UpdateTodoAssignment(ctx context.Context, input *models.UpdateTodoAssignmentInput) ([]models.TodoItem, error) {
	body := input.Body

	counselorID, err := auth.GetCounselorID(ctx)
	if err != nil {
		return nil, errs.Unauthorized()
	}
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, errs.Unauthorized()
	}

	hasCollege := body.CollegeID != nil
	hasMedia := body.MediaID != nil
	if hasCollege && hasMedia {
		return nil, errs.BadRequest("a todo item can be linked to at most one of essay, media, college, or extracurricular")
	}
	if err := validateLink(body.LinkType, false, hasMedia, hasCollege, false); err != nil {
		return nil, err
	}

	return h.TodoItemRepository.UpdateTodoAssignment(ctx, models.UpdateTodoAssignmentParams{
		AssignmentID:     input.AssignmentID,
		CounselorID:      counselorID,
		UserID:           userID,
		Title:            body.Title,
		Description:      body.Description,
		Deadline:         body.Deadline,
		LinkType:         body.LinkType,
		CollegeID:        body.CollegeID,
		MediaID:          body.MediaID,
		AddStudentIDs:    body.AddStudentIDs,
		RemoveStudentIDs: body.RemoveStudentIDs,
	})
}