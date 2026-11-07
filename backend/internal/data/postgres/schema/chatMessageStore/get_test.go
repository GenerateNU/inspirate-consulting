package chatMessageRepository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/pagination"
)

func messageIDs(msgs []models.ChatMessage) []uuid.UUID {
	ids := make([]uuid.UUID, len(msgs))
	for i, m := range msgs {
		ids[i] = m.ID
	}
	return ids
}

func TestListChatMessages_OnlyConversationNewestFirst(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	alice := createTestUser(t, db, "alice")
	bob := createTestUser(t, db, "bob")
	carol := createTestUser(t, db, "carol")

	m1 := insertMessage(t, db, alice, bob, "1", baseTime)
	m2 := insertMessage(t, db, bob, alice, "2", baseTime.Add(time.Minute))
	m3 := insertMessage(t, db, alice, bob, "3", baseTime.Add(2*time.Minute))
	insertMessage(t, db, alice, carol, "other conversation", baseTime.Add(3*time.Minute))
	insertMessage(t, db, carol, bob, "other conversation", baseTime.Add(4*time.Minute))

	got, err := repo.ListChatMessages(context.Background(), alice, bob, nil, 10)
	require.NoError(t, err)
	assert.Equal(t, []uuid.UUID{m3.ID, m2.ID, m1.ID}, messageIDs(got))

	// The same conversation is returned from the other participant's side.
	got, err = repo.ListChatMessages(context.Background(), bob, alice, nil, 10)
	require.NoError(t, err)
	assert.Equal(t, []uuid.UUID{m3.ID, m2.ID, m1.ID}, messageIDs(got))
}

func TestListChatMessages_RespectsLimit(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	alice := createTestUser(t, db, "alice")
	bob := createTestUser(t, db, "bob")
	for i := range 5 {
		insertMessage(t, db, alice, bob, "msg", baseTime.Add(time.Duration(i)*time.Minute))
	}

	got, err := repo.ListChatMessages(context.Background(), alice, bob, nil, 3)
	require.NoError(t, err)
	assert.Len(t, got, 3)
}

func TestListChatMessages_CursorPagesWithoutGapsOrDuplicates(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	alice := createTestUser(t, db, "alice")
	bob := createTestUser(t, db, "bob")

	// Several messages share a timestamp so the id tie-breaker is exercised.
	var want []models.ChatMessage
	for i := range 7 {
		want = append(want, insertMessage(t, db, alice, bob, "msg", baseTime.Add(time.Duration(i/2)*time.Minute)))
	}

	all, err := repo.ListChatMessages(context.Background(), alice, bob, nil, 100)
	require.NoError(t, err)
	require.Len(t, all, len(want))

	var paged []models.ChatMessage
	var before *pagination.TimeIDKey
	for {
		page, err := repo.ListChatMessages(context.Background(), alice, bob, before, 2)
		require.NoError(t, err)
		if len(page) == 0 {
			break
		}
		paged = append(paged, page...)
		last := page[len(page)-1]
		before = &pagination.TimeIDKey{CreatedAt: last.CreatedAt, ID: last.ID}
	}

	assert.Equal(t, messageIDs(all), messageIDs(paged))
}

func TestListChatMessages_EmptyConversation(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	alice := createTestUser(t, db, "alice")

	got, err := repo.ListChatMessages(context.Background(), alice, uuid.New(), nil, 10)
	require.NoError(t, err)
	assert.Empty(t, got)
}
