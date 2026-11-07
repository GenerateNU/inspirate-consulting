package user

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func TestHandler_UpdateUser(t *testing.T) {
	t.Parallel()

	name := "New Name"
	input := models.UpdateUserBody{Name: &name}

	t.Run("updates the logged-in user", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()
		ctx := auth.WithUserID(context.Background(), userID)

		repo := mocks.NewUserRepository(t)
		repo.On("UpdateUser", ctx, userID, input).Return(&models.User{ID: userID, Name: name}, nil)

		res, err := NewHandler(repo, nil).UpdateUser(ctx, input)

		require.NoError(t, err)
		assert.Equal(t, name, res.Name)
	})

	t.Run("caller without a user id is unauthorized", func(t *testing.T) {
		t.Parallel()

		repo := mocks.NewUserRepository(t) // UpdateUser should not be called

		res, err := NewHandler(repo, nil).UpdateUser(context.Background(), input)

		assert.Nil(t, res)
		assert.Equal(t, errs.Unauthorized(), err)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()
		ctx := auth.WithUserID(context.Background(), userID)
		notFound := errs.NotFound("user", "id", userID.String())

		repo := mocks.NewUserRepository(t)
		repo.On("UpdateUser", ctx, userID, input).Return(nil, notFound)

		res, err := NewHandler(repo, nil).UpdateUser(ctx, input)

		assert.Nil(t, res)
		assert.Equal(t, notFound, err)
	})
}
