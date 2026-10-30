package user

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) GetCurrentUser(ctx context.Context) (*models.GetCurrentUserOutput, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, errs.Unauthorized()
	}

	user, err := h.UserRepository.FetchUser(ctx, models.FetchUserInput{ID: userID})
	if err != nil {
		return nil, err
	}

	current := &models.CurrentUser{User: *user.Body}
	// The identity middleware only sets the profile ID that matches the caller's role
	if studentID, err := auth.GetStudentID(ctx); err == nil {
		current.Role = models.StudentRole
		current.StudentID = &studentID
	} else if counselorID, err := auth.GetCounselorID(ctx); err == nil {
		current.Role = models.CounselorRole
		current.CounselorID = &counselorID
	}

	return &models.GetCurrentUserOutput{Body: current}, nil
}
