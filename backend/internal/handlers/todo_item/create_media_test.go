package todoitem

import (
	"context"
	"net/http"
	"testing"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Creating a media-linked todo grants access, so only counselors may link media.
func TestHandler_CreateTodoItem_LinkedMedia(t *testing.T) {
	t.Parallel()

	studentID := uuid.New()
	userID := uuid.New()
	mediaID := uuid.NewString()
	module := models.TodoLinkModule
	studentCtx := auth.WithUserID(auth.WithStudentID(context.Background(), studentID), userID)
	counselorCtx := auth.WithUserID(auth.WithCounselorID(context.Background(), uuid.New()), userID)

	t.Run("counselor can assign a video", func(t *testing.T) {
		t.Parallel()

		input := &models.CreateTodoItemRequestBody{StudentID: studentID.String(), TodoDescription: "Watch the FAFSA video", MediaID: &mediaID, LinkType: &module}
		created := &models.TodoItem{ID: uuid.NewString(), StudentID: studentID.String(), MediaID: &mediaID, LinkType: &module}

		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("CreateTodoItem", mock.Anything, studentID.String(), userID.String(), *input).Return(created, nil)

		handler := NewHandler(mockRepo, nil)
		res, err := handler.CreateTodoItem(counselorCtx, input)

		assert.NoError(t, err)
		assert.Equal(t, created, res)
	})

	t.Run("student cannot link a video", func(t *testing.T) {
		t.Parallel()

		input := &models.CreateTodoItemRequestBody{TodoDescription: "Watch a video", MediaID: &mediaID, LinkType: &module}

		// No CreateTodoItem expectation: the todo (and the grant) must never be written
		mockRepo := mocks.NewTodoItemRepository(t)

		handler := NewHandler(mockRepo, nil)
		res, err := handler.CreateTodoItem(studentCtx, input)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
	})

	t.Run("student can still create a todo without media", func(t *testing.T) {
		t.Parallel()

		input := &models.CreateTodoItemRequestBody{TodoDescription: "Draft essay"}
		created := &models.TodoItem{ID: uuid.NewString(), StudentID: studentID.String()}

		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("CreateTodoItem", mock.Anything, studentID.String(), userID.String(), *input).Return(created, nil)

		handler := NewHandler(mockRepo, nil)
		_, err := handler.CreateTodoItem(studentCtx, input)

		assert.NoError(t, err)
	})
}
