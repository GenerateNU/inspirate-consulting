package vectorsearch

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_SearchChats(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	ctx := auth.WithUserID(context.Background(), userID)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		embedder := newFakeEmbeddingClient()
		expected := []models.ChatMessageSearchResult{
			{ChatMessage: models.ChatMessage{ID: uuid.New(), SenderID: userID, Message: "essay feedback"}, Similarity: 0.9},
		}
		mockVector := mocks.NewVectorRepository(t)
		mockVector.On("SearchChatMessages", ctx, embedder.embedding, userID, searchLimit).Return(expected, nil)

		res, err := NewHandler(mockVector, nil, embedder).SearchChats(ctx, "what did we say about my essay")

		require.NoError(t, err)
		assert.Equal(t, expected, res)
		assert.Equal(t, "what did we say about my essay", embedder.gotQuery)
	})

	t.Run("not configured without an embedding client", func(t *testing.T) {
		t.Parallel()

		mockVector := mocks.NewVectorRepository(t)
		res, err := NewHandler(mockVector, nil, nil).SearchChats(ctx, "essay")

		assert.Nil(t, res)
		assertHTTPStatus(t, err, http.StatusServiceUnavailable)
	})

	for name, query := range map[string]string{"empty query": "", "whitespace-only query": "  \n\t "} {
		t.Run(name+" is rejected", func(t *testing.T) {
			t.Parallel()

			embedder := newFakeEmbeddingClient()
			mockVector := mocks.NewVectorRepository(t)
			res, err := NewHandler(mockVector, nil, embedder).SearchChats(ctx, query)

			assert.Nil(t, res)
			assertHTTPStatus(t, err, http.StatusBadRequest)
			assert.Zero(t, embedder.calls, "should not call Cloudflare for a blank query")
		})
	}

	t.Run("unauthorized without a user", func(t *testing.T) {
		t.Parallel()

		embedder := newFakeEmbeddingClient()
		mockVector := mocks.NewVectorRepository(t)
		res, err := NewHandler(mockVector, nil, embedder).SearchChats(context.Background(), "essay")

		assert.Nil(t, res)
		assertHTTPStatus(t, err, http.StatusUnauthorized)
		assert.Zero(t, embedder.calls)
	})

	t.Run("embedding error is returned", func(t *testing.T) {
		t.Parallel()

		embedder := newFakeEmbeddingClient()
		embedder.err = errors.New("cloudflare unavailable")
		mockVector := mocks.NewVectorRepository(t)

		res, err := NewHandler(mockVector, nil, embedder).SearchChats(ctx, "essay")

		assert.Nil(t, res)
		assert.ErrorIs(t, err, embedder.err)
	})

	t.Run("repository error is returned", func(t *testing.T) {
		t.Parallel()

		embedder := newFakeEmbeddingClient()
		mockVector := mocks.NewVectorRepository(t)
		mockVector.On("SearchChatMessages", ctx, embedder.embedding, userID, searchLimit).Return(nil, errors.New("database error"))

		res, err := NewHandler(mockVector, nil, embedder).SearchChats(ctx, "essay")

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
