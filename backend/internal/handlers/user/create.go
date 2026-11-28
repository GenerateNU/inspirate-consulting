package user

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/passwordgen"
	"inspirate-consulting/internal/supabase"
	"log/slog"
)

func (h *Handler) CreateUser(ctx context.Context, input *models.CreateUserInput, supabase supabase.SupabaseInterface) (*models.CreateUserOutput, error) {
	// checked before signup so a bad request doesn't create a user signup in supabase without a user row in postgres
	if input.Body.Role == models.StudentRole && (input.Body.GraduationYear == nil || input.Body.GPA == nil || input.Body.CounselorID == nil) {
		return nil, errs.BadRequest("graduation_year, gpa, and counselor_id are required for students")
	}

	tempPassword, err := passwordgen.GenerateTempPassword()
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

	notification := &models.AccountEmailNotificationRequestBody{
		UserEmail:    input.Body.Email,
		Name:         input.Body.Name,
		TempPassword: tempPassword,
		EmailType:    models.AccountCreated,
	}
	if _, err := h.SQSClient.SendAccountEmailNotification(ctx, notification); err != nil {
		// the user already exists, so don't fail the request; hand the temp password back
		// so the creator can pass it on, otherwise the new user could never sign in
		slog.Warn("failed to queue account-created email", "err", err, "email", input.Body.Email)
		user.Body.TempPassword = &tempPassword
		return user, nil
	}

	user.Body.EmailSent = true
	return user, nil
}
