package pfp

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

func TestHandler_PresignUpload(t *testing.T) {
	t.Parallel()

	input := models.PresignPfpUploadRequestBody{ContentType: "image/jpeg", ContentLength: 2048}

	t.Run("presigns for the logged-in user", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()
		ctx := auth.WithUserID(context.Background(), userID)
		expected := &models.PresignUploadResponse{S3Key: "pfps/" + userID.String() + "/a.jpg", UploadURL: "https://bucket/upload"}

		pfpRepo := mocks.NewPfpRepository(t)
		pfpRepo.On("PresignUpload", ctx, userID, input).Return(expected, nil)

		res, err := NewHandler(mocks.NewUserRepository(t), pfpRepo).PresignUpload(ctx, input)

		require.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("caller without a user id is unauthorized", func(t *testing.T) {
		t.Parallel()

		res, err := NewHandler(mocks.NewUserRepository(t), mocks.NewPfpRepository(t)).PresignUpload(context.Background(), input)

		assert.Nil(t, res)
		assert.Equal(t, errs.Unauthorized(), err)
	})
}
