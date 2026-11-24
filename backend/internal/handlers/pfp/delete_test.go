package pfp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func TestHandler_DeletePfp(t *testing.T) {
	t.Parallel()

	t.Run("clears the key and deletes the picture", func(t *testing.T) {
		t.Parallel()

		userID, ctx := userCtx()
		oldKey := "pfps/" + userID.String() + "/old.jpg"

		userRepo := mocks.NewUserRepository(t)
		userRepo.On("SetPfpKey", ctx, userID, (*string)(nil)).Return(&models.User{ID: userID}, &oldKey, nil)

		pfpRepo := mocks.NewPfpRepository(t)
		pfpRepo.On("DeletePfp", ctx, userID, oldKey).Return(nil)

		require.NoError(t, NewHandler(userRepo, pfpRepo).DeletePfp(ctx))
	})

	t.Run("user without a picture gets not found", func(t *testing.T) {
		t.Parallel()

		userID, ctx := userCtx()

		userRepo := mocks.NewUserRepository(t)
		userRepo.On("SetPfpKey", ctx, userID, (*string)(nil)).Return(&models.User{ID: userID}, nil, nil)

		err := NewHandler(userRepo, mocks.NewPfpRepository(t)).DeletePfp(ctx)

		assert.Equal(t, errs.NotFound("profile picture", "user id", userID.String()), err)
	})

	t.Run("caller without a user id is unauthorized", func(t *testing.T) {
		t.Parallel()

		err := NewHandler(mocks.NewUserRepository(t), mocks.NewPfpRepository(t)).DeletePfp(context.Background())

		assert.Equal(t, errs.Unauthorized(), err)
	})
}
