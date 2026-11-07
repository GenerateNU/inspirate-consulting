package chatMessageRepository

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func TestCreateChatMessage(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	sender := createTestUser(t, db, "sender")
	recipient := createTestUser(t, db, "recipient")

	created, err := repo.CreateChatMessage(context.Background(), sender, models.CreateChatMessageRequestBody{
		RecipientID: recipient,
		Message:     "Hello",
	})

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, created.ID)
	assert.Equal(t, sender, created.SenderID)
	assert.Equal(t, recipient, created.RecipientID)
	assert.Equal(t, "Hello", created.Message)
	assert.False(t, created.CreatedAt.IsZero())
	assert.Nil(t, created.EditedAt)
	assert.Nil(t, created.ReadAt)
}

func TestCreateChatMessage_UnknownRecipient(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	sender := createTestUser(t, db, "sender")

	created, err := repo.CreateChatMessage(context.Background(), sender, models.CreateChatMessageRequestBody{
		RecipientID: uuid.New(),
		Message:     "Hello?",
	})

	assert.Nil(t, created)
	var httpErr errs.HTTPError
	require.True(t, errors.As(err, &httpErr), "expected errs.HTTPError, got %v", err)
	assert.Equal(t, http.StatusNotFound, httpErr.Code)
}
