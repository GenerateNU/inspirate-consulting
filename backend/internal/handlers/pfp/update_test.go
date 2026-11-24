package pfp

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func TestHandler_UpdatePfp(t *testing.T) {
	t.Parallel()

	t.Run("saves the key and deletes the replaced picture", func(t *testing.T) {
		t.Parallel()

		userID, ctx := userCtx()
		newKey := "pfps/" + userID.String() + "/new.jpg"
		oldKey := "pfps/" + userID.String() + "/old.jpg"

		pfpRepo := mocks.NewPfpRepository(t)
		pfpRepo.On("VerifyUpload", ctx, userID, newKey).Return(nil)
		pfpRepo.On("DeletePfp", ctx, userID, oldKey).Return(nil)
		pfpRepo.On("GetDownloadURL", ctx, newKey).Return("https://bucket/new.jpg", nil)

		userRepo := mocks.NewUserRepository(t)
		userRepo.On("SetPfpKey", ctx, userID, &newKey).Return(&models.User{ID: userID, PfpKey: &newKey}, &oldKey, nil)

		res, err := NewHandler(userRepo, pfpRepo).UpdatePfp(ctx, newKey)

		require.NoError(t, err)
		assert.Equal(t, &models.Pfp{S3Key: newKey, DownloadURL: "https://bucket/new.jpg"}, res)
	})

	t.Run("first picture deletes nothing", func(t *testing.T) {
		t.Parallel()

		userID, ctx := userCtx()
		newKey := "pfps/" + userID.String() + "/new.jpg"

		pfpRepo := mocks.NewPfpRepository(t)
		pfpRepo.On("VerifyUpload", ctx, userID, newKey).Return(nil)
		pfpRepo.On("GetDownloadURL", ctx, newKey).Return("https://bucket/new.jpg", nil)

		userRepo := mocks.NewUserRepository(t)
		userRepo.On("SetPfpKey", ctx, userID, &newKey).Return(&models.User{ID: userID, PfpKey: &newKey}, nil, nil)

		_, err := NewHandler(userRepo, pfpRepo).UpdatePfp(ctx, newKey)

		require.NoError(t, err)
		pfpRepo.AssertNotCalled(t, "DeletePfp", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("failed cleanup still saves the new picture", func(t *testing.T) {
		t.Parallel()

		userID, ctx := userCtx()
		newKey := "pfps/" + userID.String() + "/new.jpg"
		oldKey := "pfps/" + userID.String() + "/old.jpg"

		pfpRepo := mocks.NewPfpRepository(t)
		pfpRepo.On("VerifyUpload", ctx, userID, newKey).Return(nil)
		pfpRepo.On("DeletePfp", ctx, userID, oldKey).Return(errors.New("s3 unavailable"))
		pfpRepo.On("GetDownloadURL", ctx, newKey).Return("https://bucket/new.jpg", nil)

		userRepo := mocks.NewUserRepository(t)
		userRepo.On("SetPfpKey", ctx, userID, &newKey).Return(&models.User{ID: userID, PfpKey: &newKey}, &oldKey, nil)

		res, err := NewHandler(userRepo, pfpRepo).UpdatePfp(ctx, newKey)

		require.NoError(t, err)
		assert.Equal(t, newKey, res.S3Key)
	})

	t.Run("rejected upload is not saved", func(t *testing.T) {
		t.Parallel()

		userID, ctx := userCtx()
		otherKey := "pfps/" + uuid.NewString() + "/theirs.jpg"
		forbidden := errs.Forbidden("you can only use your own profile picture")

		pfpRepo := mocks.NewPfpRepository(t)
		pfpRepo.On("VerifyUpload", ctx, userID, otherKey).Return(forbidden)

		userRepo := mocks.NewUserRepository(t) // SetPfpKey must not be called

		res, err := NewHandler(userRepo, pfpRepo).UpdatePfp(ctx, otherKey)

		assert.Nil(t, res)
		assert.Equal(t, forbidden, err)
	})

	t.Run("caller without a user id is unauthorized", func(t *testing.T) {
		t.Parallel()

		res, err := NewHandler(mocks.NewUserRepository(t), mocks.NewPfpRepository(t)).UpdatePfp(context.Background(), "pfps/x/a.jpg")

		assert.Nil(t, res)
		assert.Equal(t, errs.Unauthorized(), err)
	})
}
