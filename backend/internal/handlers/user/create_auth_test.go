package user

import (
	"context"
	"errors"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
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
	if role == models.StudentRole {
		year := models.Senior
		gpa := 4
		counselorID := uuid.New()
		in.Body.Year = &year
		in.Body.GPA = &gpa
		in.Body.CounselorID = &counselorID
	}
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
			supabaseID := uuid.New()
			sb := &supabase.MockSupabase{
				SignupFunc: func(email, password, role string) (models.SignupResponse, error) {
					gotEmail, gotPassword, gotRole = email, password, role
					return models.SignupResponse{User: models.UserSignupResponse{ID: supabaseID}}, nil
				},
			}
			repo := mocks.NewUserRepository(t)
			// the users row is linked to the Supabase account just created
			repo.On("CreateUser", mock.Anything, *in, supabaseID).
				Return(&models.CreateUserOutput{Body: &models.CreateUserBody{User: &models.User{Name: "New User"}}}, nil).Once()

			out, err := NewHandler(repo).CreateUser(ctx, in, sb)
			require.NoError(t, err)

			assert.Equal(t, "new@example.com", gotEmail)
			assert.Equal(t, string(role), gotRole, "role must reach Supabase app_metadata")
			assert.NotEmpty(t, out.Body.TempPassword)
			assert.Equal(t, out.Body.TempPassword, gotPassword, "the returned temp password is the one Supabase got")
		})
	}

	missingFields := map[string]func(in *models.CreateUserInput){
		"year":         func(in *models.CreateUserInput) { in.Body.Year = nil },
		"gpa":          func(in *models.CreateUserInput) { in.Body.GPA = nil },
		"counselor_id": func(in *models.CreateUserInput) { in.Body.CounselorID = nil },
	}
	for field, clear := range missingFields {
		t.Run("student without "+field+" is rejected before signup", func(t *testing.T) {
			t.Parallel()
			signedUp := false
			sb := &supabase.MockSupabase{
				SignupFunc: func(string, string, string) (models.SignupResponse, error) {
					signedUp = true
					return models.SignupResponse{}, nil
				},
			}
			repo := mocks.NewUserRepository(t) // CreateUser must not be called

			in := createInput("X", "x@example.com", models.StudentRole)
			clear(in)

			out, err := NewHandler(repo).CreateUser(ctx, in, sb)
			assert.Nil(t, out)
			assert.Equal(t, errs.BadRequest("year, gpa, and counselor_id are required for students"), err)
			assert.False(t, signedUp, "no Supabase account should be created for an invalid request")
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
		repo.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("db down"))

		out, err := NewHandler(repo).CreateUser(ctx, createInput("X", "x@example.com", models.StudentRole), &supabase.MockSupabase{})
		assert.Nil(t, out)
		assert.EqualError(t, err, "db down")
	})
}
