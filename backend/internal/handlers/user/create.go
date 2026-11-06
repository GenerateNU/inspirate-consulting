package user

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
)

func (h *Handler) CreateUser(ctx context.Context, input *models.CreateUserInput, supabase supabase.SupabaseInterface) (*models.CreateUserOutput, error) {
	// checked before signup so a bad request doesn't create a user signup in supabase without a user row in postgres
	if input.Body.Role == models.StudentRole && (input.Body.Year == nil || input.Body.GPA == nil || input.Body.CounselorID == nil) {
		return nil, errs.BadRequest("year, gpa, and counselor_id are required for students")
	}

	tempPassword, err := generateTempPassword()
	if err != nil {
		return nil, err
	}

	signup_response, err := supabase.Signup(input.Body.Email, tempPassword, auth.Client, string(input.Body.Role))
	if err != nil {
		return nil, err
	}

	user, err := h.UserRepository.CreateUser(ctx, *input, signup_response.User.ID)
	if err != nil {
		return nil, err
	}

	user.Body.TempPassword = tempPassword
	return user, nil
}
