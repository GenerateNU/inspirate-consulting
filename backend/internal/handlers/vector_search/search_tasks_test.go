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

func TestHandler_SearchTasks(t *testing.T) {
	t.Parallel()

	t.Run("student searches their own todos", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.New()
		ctx := studentCtx(studentID)
		embedder := newFakeEmbeddingClient()
		expected := []models.TodoItemSearchResult{
			{TodoItem: models.TodoItem{ID: uuid.NewString(), StudentID: studentID.String(), TodoDescription: "Finish essay"}, Similarity: 0.8},
		}
		mockVector := mocks.NewVectorRepository(t)
		mockVector.On("SearchTodoItems", ctx, embedder.embedding, []uuid.UUID{studentID}, searchLimit).Return(expected, nil)

		res, err := NewHandler(mockVector, nil, embedder).SearchTasks(ctx, "essay deadlines")

		require.NoError(t, err)
		assert.Equal(t, expected, res)
		assert.Equal(t, "essay deadlines", embedder.gotQuery)
	})

	t.Run("counselor searches all their students' todos", func(t *testing.T) {
		t.Parallel()

		counselorID := uuid.New()
		studentA := uuid.New()
		studentB := uuid.New()
		ctx := counselorCtx(counselorID)
		embedder := newFakeEmbeddingClient()
		expected := []models.TodoItemSearchResult{
			{TodoItem: models.TodoItem{ID: uuid.NewString(), StudentID: studentB.String()}, Similarity: 0.7},
		}

		mockCounselor := mocks.NewCounselorRepository(t)
		mockCounselor.On("ListStudents", ctx, counselorID).Return([]models.CounselorStudentRow{
			{StudentID: studentA, Name: "A"},
			{StudentID: studentB, Name: "B"},
		}, nil)
		mockVector := mocks.NewVectorRepository(t)
		mockVector.On("SearchTodoItems", ctx, embedder.embedding, []uuid.UUID{studentA, studentB}, searchLimit).Return(expected, nil)

		res, err := NewHandler(mockVector, mockCounselor, embedder).SearchTasks(ctx, "who is behind on essays")

		require.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("counselor with no students gets no results without calling Cloudflare", func(t *testing.T) {
		t.Parallel()

		counselorID := uuid.New()
		ctx := counselorCtx(counselorID)
		embedder := newFakeEmbeddingClient()
		mockCounselor := mocks.NewCounselorRepository(t)
		mockCounselor.On("ListStudents", ctx, counselorID).Return([]models.CounselorStudentRow{}, nil)
		mockVector := mocks.NewVectorRepository(t)

		res, err := NewHandler(mockVector, mockCounselor, embedder).SearchTasks(ctx, "essays")

		require.NoError(t, err)
		assert.Empty(t, res)
		assert.Zero(t, embedder.calls)
	})

	t.Run("counselor student lookup error is returned", func(t *testing.T) {
		t.Parallel()

		counselorID := uuid.New()
		ctx := counselorCtx(counselorID)
		mockCounselor := mocks.NewCounselorRepository(t)
		mockCounselor.On("ListStudents", ctx, counselorID).Return(nil, errors.New("database error"))
		mockVector := mocks.NewVectorRepository(t)

		res, err := NewHandler(mockVector, mockCounselor, newFakeEmbeddingClient()).SearchTasks(ctx, "essays")

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})

	t.Run("not configured without an embedding client", func(t *testing.T) {
		t.Parallel()

		mockVector := mocks.NewVectorRepository(t)
		res, err := NewHandler(mockVector, nil, nil).SearchTasks(studentCtx(uuid.New()), "essays")

		assert.Nil(t, res)
		assertHTTPStatus(t, err, http.StatusServiceUnavailable)
	})

	for name, query := range map[string]string{"empty query": "", "whitespace-only query": "  \n\t "} {
		t.Run(name+" is rejected", func(t *testing.T) {
			t.Parallel()

			embedder := newFakeEmbeddingClient()
			mockVector := mocks.NewVectorRepository(t)
			res, err := NewHandler(mockVector, nil, embedder).SearchTasks(studentCtx(uuid.New()), query)

			assert.Nil(t, res)
			assertHTTPStatus(t, err, http.StatusBadRequest)
			assert.Zero(t, embedder.calls)
		})
	}

	t.Run("student role without a student id is unauthorized", func(t *testing.T) {
		t.Parallel()

		embedder := newFakeEmbeddingClient()
		mockVector := mocks.NewVectorRepository(t)
		ctx := withRole(context.Background(), models.StudentRole)

		res, err := NewHandler(mockVector, nil, embedder).SearchTasks(ctx, "essays")

		assert.Nil(t, res)
		assertHTTPStatus(t, err, http.StatusUnauthorized)
		assert.Zero(t, embedder.calls)
	})

	t.Run("counselor role without a counselor id is unauthorized", func(t *testing.T) {
		t.Parallel()

		embedder := newFakeEmbeddingClient()
		mockVector := mocks.NewVectorRepository(t)
		ctx := withRole(context.Background(), models.CounselorRole)

		res, err := NewHandler(mockVector, nil, embedder).SearchTasks(ctx, "essays")

		assert.Nil(t, res)
		assertHTTPStatus(t, err, http.StatusUnauthorized)
	})

	t.Run("no role is forbidden", func(t *testing.T) {
		t.Parallel()

		embedder := newFakeEmbeddingClient()
		mockVector := mocks.NewVectorRepository(t)
		ctx := auth.WithUserID(context.Background(), uuid.New())

		res, err := NewHandler(mockVector, nil, embedder).SearchTasks(ctx, "essays")

		assert.Nil(t, res)
		assertHTTPStatus(t, err, http.StatusForbidden)
		assert.Zero(t, embedder.calls)
	})

	t.Run("embedding error is returned", func(t *testing.T) {
		t.Parallel()

		embedder := newFakeEmbeddingClient()
		embedder.err = errors.New("cloudflare unavailable")
		mockVector := mocks.NewVectorRepository(t)

		res, err := NewHandler(mockVector, nil, embedder).SearchTasks(studentCtx(uuid.New()), "essays")

		assert.Nil(t, res)
		assert.ErrorIs(t, err, embedder.err)
	})

	t.Run("repository error is returned", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.New()
		ctx := studentCtx(studentID)
		embedder := newFakeEmbeddingClient()
		mockVector := mocks.NewVectorRepository(t)
		mockVector.On("SearchTodoItems", ctx, embedder.embedding, []uuid.UUID{studentID}, searchLimit).Return(nil, errors.New("database error"))

		res, err := NewHandler(mockVector, nil, embedder).SearchTasks(ctx, "essays")

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
