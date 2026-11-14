package extracurricular

import (
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

// childOp runs one of the note/achievement handlers against a parent id and
// returns only the error, so the shared authorization cases can be table-driven.
type childOp struct {
	name   string
	expect func(m *mocks.ExtracurricularRepository, parentID int64)
	run    func(h *Handler, caller uuid.UUID, parentID int64) error
}

func childOps() []childOp {
	return []childOp{
		{"create note", func(m *mocks.ExtracurricularRepository, id int64) {
			m.On("CreateExtracurricularNote", mock.Anything, id, models.ExtracurricularNoteRequest{Body: "n"}).Return(&models.ExtracurricularNote{}, nil)
		}, func(h *Handler, c uuid.UUID, id int64) error {
			_, err := h.CreateExtracurricularNote(studentRoleCtx(c), &models.CreateExtracurricularNoteInput{
				ExtracurricularID: id, Body: models.ExtracurricularNoteRequest{Body: "n"},
			})
			return err
		}},
		{"update note", func(m *mocks.ExtracurricularRepository, id int64) {
			m.On("UpdateExtracurricularNote", mock.Anything, id, int64(3), models.ExtracurricularNoteRequest{Body: "n"}).Return(&models.ExtracurricularNote{}, nil)
		}, func(h *Handler, c uuid.UUID, id int64) error {
			_, err := h.UpdateExtracurricularNote(studentRoleCtx(c), &models.UpdateExtracurricularNoteInput{
				ExtracurricularID: id, NoteID: 3, Body: models.ExtracurricularNoteRequest{Body: "n"},
			})
			return err
		}},
		{"delete note", func(m *mocks.ExtracurricularRepository, id int64) {
			m.On("DeleteExtracurricularNote", mock.Anything, id, int64(3)).Return(nil)
		}, func(h *Handler, c uuid.UUID, id int64) error {
			return h.DeleteExtracurricularNote(studentRoleCtx(c), &models.DeleteExtracurricularNoteInput{ExtracurricularID: id, NoteID: 3})
		}},
		{"create achievement", func(m *mocks.ExtracurricularRepository, id int64) {
			m.On("CreateExtracurricularAchievement", mock.Anything, id, mock.Anything).Return(&models.ExtracurricularAchievement{}, nil)
		}, func(h *Handler, c uuid.UUID, id int64) error {
			_, err := h.CreateExtracurricularAchievement(studentRoleCtx(c), &models.CreateExtracurricularAchievementInput{
				ExtracurricularID: id, Body: models.CreateExtracurricularAchievementRequest{EntryDate: "2025-01-01", Body: "a"},
			})
			return err
		}},
		{"update achievement", func(m *mocks.ExtracurricularRepository, id int64) {
			m.On("UpdateExtracurricularAchievement", mock.Anything, id, int64(4), mock.Anything).Return(&models.ExtracurricularAchievement{}, nil)
		}, func(h *Handler, c uuid.UUID, id int64) error {
			_, err := h.UpdateExtracurricularAchievement(studentRoleCtx(c), &models.UpdateExtracurricularAchievementInput{
				ExtracurricularID: id, AchievementID: 4, Body: models.UpdateExtracurricularAchievementRequest{Body: ptr("a")},
			})
			return err
		}},
		{"delete achievement", func(m *mocks.ExtracurricularRepository, id int64) {
			m.On("DeleteExtracurricularAchievement", mock.Anything, id, int64(4)).Return(nil)
		}, func(h *Handler, c uuid.UUID, id int64) error {
			return h.DeleteExtracurricularAchievement(studentRoleCtx(c), &models.DeleteExtracurricularAchievementInput{ExtracurricularID: id, AchievementID: 4})
		}},
	}
}

// Every note/achievement handler checks access to the parent before touching the child.
func TestHandler_ExtracurricularChildren_Authorization(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	const parentID int64 = 30
	parent := &models.ExtracurricularDetail{Extracurricular: models.Extracurricular{ID: parentID, StudentID: ownerID.String()}}

	for _, op := range childOps() {
		t.Run(op.name+": owner reaches the repository", func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewExtracurricularRepository(t)
			mockRepo.On("GetExtracurricular", mock.Anything, parentID).Return(parent, nil)
			op.expect(mockRepo, parentID)

			require.NoError(t, op.run(NewHandler(mockRepo, realAuthorizer()), ownerID, parentID))
		})

		t.Run(op.name+": another student is forbidden", func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewExtracurricularRepository(t)
			mockRepo.On("GetExtracurricular", mock.Anything, parentID).Return(parent, nil)

			err := op.run(NewHandler(mockRepo, realAuthorizer()), uuid.New(), parentID)

			var httpErr errs.HTTPError
			require.ErrorAs(t, err, &httpErr)
			assert.Equal(t, http.StatusForbidden, httpErr.Code)
			assert.Len(t, mockRepo.Calls, 1, "child repository method must not be called")
		})

		t.Run(op.name+": missing parent is 404", func(t *testing.T) {
			t.Parallel()

			notFound := errs.NotFound("extracurricular", "id", parentID)
			mockRepo := mocks.NewExtracurricularRepository(t)
			mockRepo.On("GetExtracurricular", mock.Anything, parentID).Return(nil, notFound)

			err := op.run(NewHandler(mockRepo, realAuthorizer()), ownerID, parentID)

			assert.Equal(t, notFound, err)
			assert.Len(t, mockRepo.Calls, 1, "child repository method must not be called")
		})

		t.Run(op.name+": lookup error is propagated", func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewExtracurricularRepository(t)
			mockRepo.On("GetExtracurricular", mock.Anything, parentID).Return(nil, errors.New("database error"))

			assert.EqualError(t, op.run(NewHandler(mockRepo, realAuthorizer()), ownerID, parentID), "database error")
		})
	}
}

// The handler forwards the parent id, child id and body unchanged.
func TestHandler_ExtracurricularChildren_ForwardsArguments(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	parent := &models.ExtracurricularDetail{Extracurricular: models.Extracurricular{ID: 8, StudentID: ownerID.String()}}

	t.Run("update note", func(t *testing.T) {
		t.Parallel()

		want := &models.ExtracurricularNote{ID: 5, ExtracurricularID: 8, Body: "new"}
		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, int64(8)).Return(parent, nil)
		mockRepo.On("UpdateExtracurricularNote", mock.Anything, int64(8), int64(5), models.ExtracurricularNoteRequest{Body: "new"}).Return(want, nil)

		got, err := NewHandler(mockRepo, realAuthorizer()).UpdateExtracurricularNote(studentRoleCtx(ownerID), &models.UpdateExtracurricularNoteInput{
			ExtracurricularID: 8, NoteID: 5, Body: models.ExtracurricularNoteRequest{Body: "new"},
		})
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("update achievement", func(t *testing.T) {
		t.Parallel()

		body := models.UpdateExtracurricularAchievementRequest{EntryDate: ptr("2025-02-02")}
		want := &models.ExtracurricularAchievement{ID: 6, ExtracurricularID: 8}
		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, int64(8)).Return(parent, nil)
		mockRepo.On("UpdateExtracurricularAchievement", mock.Anything, int64(8), int64(6), body).Return(want, nil)

		got, err := NewHandler(mockRepo, realAuthorizer()).UpdateExtracurricularAchievement(studentRoleCtx(ownerID), &models.UpdateExtracurricularAchievementInput{
			ExtracurricularID: 8, AchievementID: 6, Body: body,
		})
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
