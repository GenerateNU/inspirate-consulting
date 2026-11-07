package chatMessageRepository

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inspirate-consulting/internal/errs"
)

func assertNotFound(t *testing.T, err error) {
	t.Helper()

	var httpErr errs.HTTPError
	require.True(t, errors.As(err, &httpErr), "expected errs.HTTPError, got %v", err)
	assert.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestEditChatMessage(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	sender := createTestUser(t, db, "sender")
	recipient := createTestUser(t, db, "recipient")
	msg := insertMessage(t, db, sender, recipient, "original", baseTime)

	edited, err := repo.EditChatMessage(context.Background(), msg.ID, sender, "edited")

	require.NoError(t, err)
	assert.Equal(t, msg.ID, edited.ID)
	assert.Equal(t, "edited", edited.Message)
	require.NotNil(t, edited.EditedAt)
	assert.True(t, edited.CreatedAt.Equal(msg.CreatedAt))
}

func TestEditChatMessage_OnlySenderCanEdit(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	sender := createTestUser(t, db, "sender")
	recipient := createTestUser(t, db, "recipient")
	msg := insertMessage(t, db, sender, recipient, "original", baseTime)

	for name, userID := range map[string]uuid.UUID{"recipient": recipient, "stranger": uuid.New()} {
		t.Run(name, func(t *testing.T) {
			_, err := repo.EditChatMessage(context.Background(), msg.ID, userID, "hijacked")
			assertNotFound(t, err)
		})
	}

	_, err := repo.EditChatMessage(context.Background(), uuid.New(), sender, "missing")
	assertNotFound(t, err)
}

func TestUpdateChatMessageReadAt(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()
	sender := createTestUser(t, db, "sender")
	recipient := createTestUser(t, db, "recipient")
	msg := insertMessage(t, db, sender, recipient, "hi", baseTime)

	firstRead := time.Now().UTC().Truncate(time.Microsecond)
	updated, err := repo.UpdateChatMessageReadAt(ctx, msg.ID, recipient, &firstRead)
	require.NoError(t, err)
	require.NotNil(t, updated.ReadAt)
	assert.True(t, updated.ReadAt.Equal(firstRead))

	// Marking an already-read message read again keeps the original read time.
	later := firstRead.Add(time.Hour)
	updated, err = repo.UpdateChatMessageReadAt(ctx, msg.ID, recipient, &later)
	require.NoError(t, err)
	require.NotNil(t, updated.ReadAt)
	assert.True(t, updated.ReadAt.Equal(firstRead))

	updated, err = repo.UpdateChatMessageReadAt(ctx, msg.ID, recipient, nil)
	require.NoError(t, err)
	assert.Nil(t, updated.ReadAt)
}

func TestUpdateChatMessageReadAt_OnlyRecipient(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	sender := createTestUser(t, db, "sender")
	recipient := createTestUser(t, db, "recipient")
	msg := insertMessage(t, db, sender, recipient, "hi", baseTime)
	now := time.Now()

	_, err := repo.UpdateChatMessageReadAt(context.Background(), msg.ID, sender, &now)
	assertNotFound(t, err)
	assert.Nil(t, getReadAt(t, db, msg.ID))

	_, err = repo.UpdateChatMessageReadAt(context.Background(), uuid.New(), recipient, &now)
	assertNotFound(t, err)
}

func TestMarkChatRead(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	me := createTestUser(t, db, "me")
	alice := createTestUser(t, db, "alice")
	bob := createTestUser(t, db, "bob")

	fromAlice1 := insertMessage(t, db, alice, me, "1", baseTime)
	fromAlice2 := insertMessage(t, db, alice, me, "2", baseTime.Add(time.Minute))
	toAlice := insertMessage(t, db, me, alice, "3", baseTime.Add(2*time.Minute))
	fromBob := insertMessage(t, db, bob, me, "4", baseTime.Add(3*time.Minute))

	require.NoError(t, repo.MarkChatRead(context.Background(), me, alice))

	assert.NotNil(t, getReadAt(t, db, fromAlice1.ID))
	assert.NotNil(t, getReadAt(t, db, fromAlice2.ID))
	// Messages I sent and other conversations are left alone.
	assert.Nil(t, getReadAt(t, db, toAlice.ID))
	assert.Nil(t, getReadAt(t, db, fromBob.ID))
}

func TestMarkChatRead_KeepsExistingReadAt(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	me := createTestUser(t, db, "me")
	alice := createTestUser(t, db, "alice")
	msg := insertMessage(t, db, alice, me, "hi", baseTime)

	earlier := baseTime.Add(time.Hour)
	_, err := db.Exec(context.Background(), `UPDATE public.chat_messages SET read_at = $1 WHERE id = $2`, earlier, msg.ID)
	require.NoError(t, err)

	require.NoError(t, repo.MarkChatRead(context.Background(), me, alice))

	readAt := getReadAt(t, db, msg.ID)
	require.NotNil(t, readAt)
	assert.True(t, readAt.Equal(earlier))
}
