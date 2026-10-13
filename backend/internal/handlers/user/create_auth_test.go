package user

import (
	"context"
	"errors"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createInput(name, email string, role models.Role) *models.CreateUserInput {
	in := &models.CreateUserInput{}
	in.Body.Name = name
	in.Body.Email = email
	in.Body.Role = role
	return in
}

func TestHandler_CreateUser_Auth(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	for _, role := range []models.Role{models.StudentRole, models.CounselorRole} {
		t.Run("signs up a "+string(role)+" in Supabase and returns the temp password", func(t *testing.T) {
			t.Parallel()
			in := createInput("New User", "new@example.com", role)

			var gotEmail, gotPassword, gotRole string
			sb := &supabase.MockSupabase{
				SignupFunc: func(email, password, role string) (models.SignupResponse, error) {
					gotEmail, gotPassword, gotRole = email, password, role
					return models.SignupResponse{User: models.UserSignupResponse{ID: uuid.New()}}, nil
				},
			}
			repo := mocks.NewUserRepository(t)
			repo.On("CreateUser", mock.Anything, *in).
				Return(&models.CreateUserOutput{Body: &models.CreateUserBody{User: &models.User{Name: "New User"}}}, nil).Once()

			out, err := NewHandler(repo).CreateUser(ctx, in, sb)
			require.NoError(t, err)

			assert.Equal(t, "new@example.com", gotEmail)
			assert.Equal(t, string(role), gotRole, "role must reach Supabase app_metadata")
			assert.NotEmpty(t, out.Body.TempPassword)
			assert.Equal(t, out.Body.TempPassword, gotPassword, "the returned temp password is the one Supabase got")
		})
	}

	t.Run("signup failure creates no users row", func(t *testing.T) {
		t.Parallel()
		sb := &supabase.MockSupabase{
			SignupFunc: func(string, string, string) (models.SignupResponse, error) {
				return models.SignupResponse{}, errors.New("email already registered")
			},
		}
		repo := mocks.NewUserRepository(t) // CreateUser must not be called

		out, err := NewHandler(repo).CreateUser(ctx, createInput("X", "x@example.com", models.StudentRole), sb)
		assert.Nil(t, out)
		assert.EqualError(t, err, "email already registered")
	})

	t.Run("database failure does not return a temp password", func(t *testing.T) {
		t.Parallel()
		repo := mocks.NewUserRepository(t)
		repo.On("CreateUser", mock.Anything, mock.Anything).Return(nil, errors.New("db down"))

		out, err := NewHandler(repo).CreateUser(ctx, createInput("X", "x@example.com", models.StudentRole), &supabase.MockSupabase{})
		assert.Nil(t, out)
		assert.EqualError(t, err, "db down")
	})
}
