package user

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) FetchCurrentUser(ctx context.Context) (*models.SessionUser, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, errs.Unauthorized()
	}

	user, err := h.UserRepository.FetchUser(ctx, models.FetchUserInput{ID: userID})
	if err != nil {
		return nil, err
	}

	session := &models.SessionUser{User: *user.Body}
	if studentID, err := auth.GetStudentID(ctx); err == nil {
		session.StudentID = &studentID
	}
	if counselorID, err := auth.GetCounselorID(ctx); err == nil {
		session.CounselorID = &counselorID
	}

	return session, nil
}
