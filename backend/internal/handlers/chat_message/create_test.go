package chatmessage

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// assertHTTPStatus checks that err is an errs.HTTPError with the given status code.
func assertHTTPStatus(t *testing.T, err error, status int) {
	t.Helper()

	var httpErr errs.HTTPError
	require.True(t, errors.As(err, &httpErr), "expected errs.HTTPError, got %T: %v", err, err)
	assert.Equal(t, status, httpErr.Code)
}

func TestHandler_CreateChatMessage(t *testing.T) {
	t.Parallel()

	senderID := uuid.New()
	recipientID := uuid.New()
	ctx := auth.WithUserID(context.Background(), senderID)
	body := models.CreateChatMessageRequestBody{RecipientID: recipientID, Message: "Hi there"}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		expected := &models.ChatMessage{ID: uuid.New(), SenderID: senderID, RecipientID: recipientID, Message: "Hi there"}
		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("CreateChatMessage", ctx, senderID, body).Return(expected, nil)

		res, err := NewHandler(mockRepo).CreateChatMessage(ctx, body)

		require.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("sending to yourself is rejected", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewChatMessageRepository(t)
		res, err := NewHandler(mockRepo).CreateChatMessage(ctx, models.CreateChatMessageRequestBody{RecipientID: senderID, Message: "Hi me"})

		assert.Nil(t, res)
		assertHTTPStatus(t, err, http.StatusBadRequest)
	})

	for name, message := range map[string]string{"empty message": "", "whitespace-only message": "  \n\t "} {
		t.Run(name+" is rejected", func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewChatMessageRepository(t)
			res, err := NewHandler(mockRepo).CreateChatMessage(ctx, models.CreateChatMessageRequestBody{RecipientID: recipientID, Message: message})

			assert.Nil(t, res)
			assertHTTPStatus(t, err, http.StatusBadRequest)
		})
	}

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("CreateChatMessage", ctx, senderID, body).Return(nil, errors.New("database error"))

		res, err := NewHandler(mockRepo).CreateChatMessage(ctx, body)

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})

	t.Run("missing user is unauthorized", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewChatMessageRepository(t)
		res, err := NewHandler(mockRepo).CreateChatMessage(context.Background(), body)

		assert.Nil(t, res)
		assertHTTPStatus(t, err, http.StatusUnauthorized)
	})
}
