package essays

import (
	"context"
	"errors"
	"net/http"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Unit tests for the UpdateEssayGroup handler logic without HTTP or database dependencies.
func TestHandler_UpdateEssayGroup(t *testing.T) {
	t.Parallel()

	essayID := uuid.New()
	groupID := uuid.New()
	studentID := uuid.New()

	assertStatus := func(t *testing.T, err error, status int) {
		t.Helper()
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, status, httpErr.GetStatus())
	}

	t.Run("moves the essay into a group", func(t *testing.T) {
		t.Parallel()

		updatedEssay := &models.Essays{
			ID:           essayID,
			EssayGroupID: &groupID,
		}

		essayRepo := mocks.NewEssayRepository(t)
		groupRepo := mocks.NewEssayGroupRepository(t)
		essayRepo.On("GetEssayStudentID", mock.Anything, essayID).Return(studentID, nil)
		groupRepo.On("GetEssayGroupStudentID", mock.Anything, groupID).Return(studentID, nil)
		essayRepo.On("UpdateEssayGroup", mock.Anything, essayID, &groupID).Return(updatedEssay, nil)

		handler := NewHandler(essayRepo, realAuthorizer(), groupRepo)
		res, err := handler.UpdateEssayGroup(studentRoleCtx(studentID), &models.UpdateEssayGroupInput{
			EssayID: essayID,
			Body:    models.UpdateEssayGroupBody{EssayGroupID: &groupID},
		})

		assert.NoError(t, err)
		assert.Equal(t, &models.UpdateEssayGroupOutput{
			Body: models.EssayBody{Essay: updatedEssay},
		}, res)
	})

	// A null essay_group_id is how a client removes an essay from its group.
	// No group lookup is needed in that case.
	t.Run("removes the essay from its group", func(t *testing.T) {
		t.Parallel()

		updatedEssay := &models.Essays{
			ID:           essayID,
			EssayGroupID: nil,
		}

		essayRepo := mocks.NewEssayRepository(t)
		groupRepo := mocks.NewEssayGroupRepository(t)
		essayRepo.On("GetEssayStudentID", mock.Anything, essayID).Return(studentID, nil)
		essayRepo.On("UpdateEssayGroup", mock.Anything, essayID, (*uuid.UUID)(nil)).Return(updatedEssay, nil)

		handler := NewHandler(essayRepo, realAuthorizer(), groupRepo)
		res, err := handler.UpdateEssayGroup(studentRoleCtx(studentID), &models.UpdateEssayGroupInput{
			EssayID: essayID,
			Body:    models.UpdateEssayGroupBody{EssayGroupID: nil},
		})

		assert.NoError(t, err)
		require.NotNil(t, res.Body.Essay)
		assert.Nil(t, res.Body.Essay.EssayGroupID)
	})

	t.Run("student cannot move another student's essay", func(t *testing.T) {
		t.Parallel()

		essayRepo := mocks.NewEssayRepository(t)
		groupRepo := mocks.NewEssayGroupRepository(t)
		essayRepo.On("GetEssayStudentID", mock.Anything, essayID).Return(uuid.New(), nil)

		handler := NewHandler(essayRepo, realAuthorizer(), groupRepo)
		res, err := handler.UpdateEssayGroup(studentRoleCtx(studentID), &models.UpdateEssayGroupInput{
			EssayID: essayID,
			Body:    models.UpdateEssayGroupBody{EssayGroupID: &groupID},
		})

		assert.Nil(t, res)
		assertStatus(t, err, http.StatusForbidden)
		essayRepo.AssertNotCalled(t, "UpdateEssayGroup", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("cannot move an essay into another student's group", func(t *testing.T) {
		t.Parallel()

		essayRepo := mocks.NewEssayRepository(t)
		groupRepo := mocks.NewEssayGroupRepository(t)
		essayRepo.On("GetEssayStudentID", mock.Anything, essayID).Return(studentID, nil)
		groupRepo.On("GetEssayGroupStudentID", mock.Anything, groupID).Return(uuid.New(), nil)

		handler := NewHandler(essayRepo, realAuthorizer(), groupRepo)
		res, err := handler.UpdateEssayGroup(studentRoleCtx(studentID), &models.UpdateEssayGroupInput{
			EssayID: essayID,
			Body:    models.UpdateEssayGroupBody{EssayGroupID: &groupID},
		})

		assert.Nil(t, res)
		assertStatus(t, err, http.StatusForbidden)
		essayRepo.AssertNotCalled(t, "UpdateEssayGroup", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("essay not found", func(t *testing.T) {
		t.Parallel()

		essayRepo := mocks.NewEssayRepository(t)
		groupRepo := mocks.NewEssayGroupRepository(t)
		essayRepo.On("GetEssayStudentID", mock.Anything, essayID).
			Return(uuid.Nil, errs.NotFound("essay", "id", essayID.String()))

		handler := NewHandler(essayRepo, nil, groupRepo)
		res, err := handler.UpdateEssayGroup(context.Background(), &models.UpdateEssayGroupInput{
			EssayID: essayID,
			Body:    models.UpdateEssayGroupBody{EssayGroupID: &groupID},
		})

		assert.Nil(t, res)
		assertStatus(t, err, http.StatusNotFound)
	})

	t.Run("group not found", func(t *testing.T) {
		t.Parallel()

		essayRepo := mocks.NewEssayRepository(t)
		groupRepo := mocks.NewEssayGroupRepository(t)
		essayRepo.On("GetEssayStudentID", mock.Anything, essayID).Return(studentID, nil)
		groupRepo.On("GetEssayGroupStudentID", mock.Anything, groupID).
			Return(uuid.Nil, errs.NotFound("group", "id", groupID.String()))

		handler := NewHandler(essayRepo, nil, groupRepo)
		res, err := handler.UpdateEssayGroup(context.Background(), &models.UpdateEssayGroupInput{
			EssayID: essayID,
			Body:    models.UpdateEssayGroupBody{EssayGroupID: &groupID},
		})

		assert.Nil(t, res)
		assertStatus(t, err, http.StatusNotFound)
		essayRepo.AssertNotCalled(t, "UpdateEssayGroup", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		essayRepo := mocks.NewEssayRepository(t)
		groupRepo := mocks.NewEssayGroupRepository(t)
		essayRepo.On("GetEssayStudentID", mock.Anything, essayID).Return(studentID, nil)
		groupRepo.On("GetEssayGroupStudentID", mock.Anything, groupID).Return(studentID, nil)
		essayRepo.On("UpdateEssayGroup", mock.Anything, essayID, &groupID).
			Return((*models.Essays)(nil), errors.New("database error"))

		handler := NewHandler(essayRepo, nil, groupRepo)
		res, err := handler.UpdateEssayGroup(context.Background(), &models.UpdateEssayGroupInput{
			EssayID: essayID,
			Body:    models.UpdateEssayGroupBody{EssayGroupID: &groupID},
		})

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
