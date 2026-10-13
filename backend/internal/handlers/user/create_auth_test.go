package user

import (
	"context"
	"errors"
	"inspirate-consulting/internal/auth/authtest"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"
	"testing"

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
		t.Run("signs up a "+string(role)+" in Supabase with a working temp password", func(t *testing.T) {
			t.Parallel()
			sb := authtest.NewFakeSupabase(authtest.SharedSigner())
			in := createInput("New User", "new@example.com", role)

			// the generated UserRepository mock drops the supabase_id argument,
			// so use the in-memory repo to see what the row is linked to
			repo := authtest.NewMemUserRepo()

			out, err := NewHandler(repo).CreateUser(ctx, in, sb)
			require.NoError(t, err)

			require.Len(t, sb.SignupCalls, 1)
			call := sb.SignupCalls[0]
			assert.Equal(t, "new@example.com", call.Email)
			assert.Equal(t, string(role), call.Role, "role must reach Supabase app_metadata")
			assert.Equal(t, out.Body.TempPassword, call.Password, "the returned temp password is the one Supabase got")

			acct, ok := sb.Account("new@example.com")
			require.True(t, ok)
			row, ok := repo.BySupabaseID(acct.ID)
			require.True(t, ok, "users row is linked to the new Supabase account")
			assert.Equal(t, "New User", row.Name)
			assert.NotNil(t, row.ResetTime, "new accounts start with a pending reset")
			assert.Equal(t, row.ID, out.Body.User.ID)

			// the counselor hands this password over, so it must actually log in
			_, err = sb.SupabaseLogin("new@example.com", out.Body.TempPassword, nil)
			assert.NoError(t, err)
		})
	}

	t.Run("signup failure creates no users row", func(t *testing.T) {
		t.Parallel()
		sb := authtest.NewFakeSupabase(authtest.SharedSigner())
		sb.SignupErr = errors.New("email already registered")
		repo := mocks.NewUserRepository(t) // CreateUser must not be called

		out, err := NewHandler(repo).CreateUser(ctx, createInput("X", "x@example.com", models.StudentRole), sb)
		assert.Nil(t, out)
		assert.EqualError(t, err, "email already registered")
	})

	t.Run("database failure does not return a temp password", func(t *testing.T) {
		t.Parallel()
		sb := authtest.NewFakeSupabase(authtest.SharedSigner())
		repo := mocks.NewUserRepository(t)
		repo.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("db down"))

		out, err := NewHandler(repo).CreateUser(ctx, createInput("X", "x@example.com", models.StudentRole), sb)
		assert.Nil(t, out)
		assert.EqualError(t, err, "db down")
	})
}
