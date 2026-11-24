package todoItemRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inspirate-consulting/internal/models"
)

// removing a student keeps their task if they already completed it
func TestUpdateTodoAssignment_RemoveKeepsCompleted(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorID, counselorUserID := createTestCounselor(t, db)
	done := createTestStudentWithCounselor(t, db, counselorID, 2018)
	notDone := createTestStudentWithCounselor(t, db, counselorID, 2018)

	assignmentID, items, err := repo.BulkCreateTodoItems(ctx, models.BulkCreateTodoItemsParams{
		CounselorID: counselorID,
		UserID:      counselorUserID,
		Title:       "Request transcripts",
		StudentIDs:  []uuid.UUID{done, notDone},
	})
	require.NoError(t, err)

	// mark one student's task complete
	for _, item := range items {
		if item.StudentID == done.String() {
			_, err := db.Exec(ctx, `UPDATE public.todo_items SET completed_at = now() WHERE id = $1`, item.ID)
			require.NoError(t, err)
		}
	}

	newTitle := "Request official transcripts"
	updated, err := repo.UpdateTodoAssignment(ctx, models.UpdateTodoAssignmentParams{
		AssignmentID:     assignmentID,
		CounselorID:      counselorID,
		UserID:           counselorUserID,
		Title:            &newTitle,
		RemoveStudentIDs: []uuid.UUID{done, notDone},
	})
	require.NoError(t, err)

	// only the completed task remains, with the new title
	require.Len(t, updated, 1)
	assert.Equal(t, done.String(), updated[0].StudentID)
	assert.Equal(t, newTitle, updated[0].Title)
}

func TestUpdateTodoAssignment_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)

	counselorID, counselorUserID := createTestCounselor(t, db)
	_, err := repo.UpdateTodoAssignment(context.Background(), models.UpdateTodoAssignmentParams{
		AssignmentID: uuid.New(),
		CounselorID:  counselorID,
		UserID:       counselorUserID,
	})
	assert.Error(t, err)
}
