package vectorRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchChatMessages_RanksBySimilarity(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	me := createTestUser(t, db, "me")
	other := createTestUser(t, db, "other")

	exact := createTestMessage(t, db, me, other, "exact match")
	unrelated := createTestMessage(t, db, me, other, "unrelated")
	received := createTestMessage(t, db, other, me, "close match I received")

	linkVector(t, db, unitVector(0), &exact, nil)
	linkVector(t, db, unitVector(1), &unrelated, nil)
	closeVector := unitVector(0)
	closeVector[1] = 1 // 45° from the query: similarity ≈ 0.707
	linkVector(t, db, closeVector, &received, nil)

	results, err := repo.SearchChatMessages(context.Background(), unitVector(0), me, 10)

	require.NoError(t, err)
	require.Len(t, results, 3)
	assert.Equal(t, exact, results[0].ID)
	assert.Equal(t, received, results[1].ID)
	assert.Equal(t, unrelated, results[2].ID)
	assert.InDelta(t, 1.0, results[0].Similarity, 0.0001)
	assert.InDelta(t, 0.7071, results[1].Similarity, 0.0001)
	assert.InDelta(t, 0.0, results[2].Similarity, 0.0001)

	// the embedded ChatMessage fields are filled in too
	assert.Equal(t, "exact match", results[0].Message)
	assert.Equal(t, me, results[0].SenderID)
	assert.Equal(t, other, results[0].RecipientID)
}

func TestSearchChatMessages_OnlyUsersOwnMessages(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	me := createTestUser(t, db, "me")
	other := createTestUser(t, db, "other")
	strangerA := createTestUser(t, db, "stranger a")
	strangerB := createTestUser(t, db, "stranger b")

	mine := createTestMessage(t, db, me, other, "mine")
	notMine := createTestMessage(t, db, strangerA, strangerB, "someone else's")
	linkVector(t, db, unitVector(1), &mine, nil)
	linkVector(t, db, unitVector(0), &notMine, nil) // a better match, but not visible to me

	results, err := repo.SearchChatMessages(context.Background(), unitVector(0), me, 10)

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, mine, results[0].ID)
}

func TestSearchChatMessages_RespectsLimit(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	me := createTestUser(t, db, "me")
	other := createTestUser(t, db, "other")
	for i := 0; i < 5; i++ {
		id := createTestMessage(t, db, me, other, "message")
		linkVector(t, db, unitVector(i), &id, nil)
	}

	results, err := repo.SearchChatMessages(context.Background(), unitVector(0), me, 2)

	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestSearchChatMessages_SkipsMessagesWithoutVectors(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	me := createTestUser(t, db, "me")
	other := createTestUser(t, db, "other")
	createTestMessage(t, db, me, other, "not embedded yet")

	results, err := repo.SearchChatMessages(context.Background(), unitVector(0), me, 10)

	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestSearchChatMessages_UnknownUser(t *testing.T) {
	t.Parallel()

	repo, _ := setupTestRepo(t)

	results, err := repo.SearchChatMessages(context.Background(), unitVector(0), uuid.New(), 10)

	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestFormatVector(t *testing.T) {
	assert.Equal(t, "[]", formatVector([]float32{}))
	assert.Equal(t, "[0.1,-2,3.5]", formatVector([]float32{0.1, -2, 3.5}))
}
