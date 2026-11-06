package chatMessageRepository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListChats(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()
	me := createTestUser(t, db, "me")
	alice := createTestUser(t, db, "alice")
	bob := createTestUser(t, db, "bob")
	carol := createTestUser(t, db, "carol")

	pfp := "avatars/alice.png"
	_, err := db.Exec(ctx, `UPDATE public.users SET pfp_key = $1 WHERE id = $2`, pfp, alice)
	require.NoError(t, err)

	// Alice: two unread from her, then my reply is the latest message.
	insertMessage(t, db, alice, me, "hi", baseTime)
	insertMessage(t, db, alice, me, "you there?", baseTime.Add(time.Minute))
	aliceLast := insertMessage(t, db, me, alice, "yes", baseTime.Add(2*time.Minute))

	// Bob: his message is the latest overall and I already read one of his messages.
	read := insertMessage(t, db, bob, me, "old", baseTime.Add(3*time.Minute))
	_, err = db.Exec(ctx, `UPDATE public.chat_messages SET read_at = now() WHERE id = $1`, read.ID)
	require.NoError(t, err)
	bobLast := insertMessage(t, db, bob, me, "new", baseTime.Add(4*time.Minute))

	// A conversation I'm not part of must not show up.
	insertMessage(t, db, alice, carol, "not mine", baseTime.Add(5*time.Minute))

	chats, err := repo.ListChats(ctx, me)
	require.NoError(t, err)
	require.Len(t, chats, 2)

	assert.Equal(t, bob, chats[0].OtherUserID)
	assert.Equal(t, "bob", chats[0].OtherUserName)
	assert.Nil(t, chats[0].OtherUserPfpKey)
	assert.Equal(t, bobLast.ID, chats[0].LastMessageID)
	assert.Equal(t, "new", chats[0].LastMessage)
	assert.Equal(t, bob, chats[0].LastSenderID)
	assert.Nil(t, chats[0].LastMessageReadAt)
	assert.Equal(t, 1, chats[0].UnreadCount)

	assert.Equal(t, alice, chats[1].OtherUserID)
	assert.Equal(t, "alice", chats[1].OtherUserName)
	require.NotNil(t, chats[1].OtherUserPfpKey)
	assert.Equal(t, pfp, *chats[1].OtherUserPfpKey)
	assert.Equal(t, aliceLast.ID, chats[1].LastMessageID)
	assert.Equal(t, me, chats[1].LastSenderID)
	// My own unread sent message doesn't count toward my unread total.
	assert.Equal(t, 2, chats[1].UnreadCount)
}

func TestListChats_OtherUserResolvedFromEitherSide(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	a := createTestUser(t, db, "a")
	b := createTestUser(t, db, "b")
	insertMessage(t, db, a, b, "hello", baseTime)

	for me, other := range map[uuid.UUID]uuid.UUID{a: b, b: a} {
		chats, err := repo.ListChats(context.Background(), me)
		require.NoError(t, err)
		require.Len(t, chats, 1)
		assert.Equal(t, other, chats[0].OtherUserID)
	}
}

func TestListChats_NoConversations(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	me := createTestUser(t, db, "lonely")

	chats, err := repo.ListChats(context.Background(), me)
	require.NoError(t, err)
	assert.Empty(t, chats)
}
