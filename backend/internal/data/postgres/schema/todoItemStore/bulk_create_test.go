package todoItemRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inspirate-consulting/internal/models"
)

func TestBulkCreateTodoItems_DedupesAndScopesToCounselor(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorA, counselorAUserID := createTestCounselor(t, db)
	studentA1 := createTestStudentWithCounselor(t, db, counselorA, "junior")
	studentA2 := createTestStudentWithCounselor(t, db, counselorA, "junior")
	counselorB, _ := createTestCounselor(t, db)
	studentB := createTestStudentWithCounselor(t, db, counselorB, "junior")

	assignmentID, items, err := repo.BulkCreateTodoItems(ctx, models.BulkCreateTodoItemsParams{
		CounselorID: counselorA,
		UserID:      counselorAUserID,
		Title:       "Request transcripts",
		// studentA1 is picked individually AND by year: should get one task.
		// studentB belongs to counselor B: should get none.
		StudentIDs: []uuid.UUID{studentA1, studentB},
		Years:      []string{"junior"},
	})
	require.NoError(t, err)

	require.Len(t, items, 2) // A1 and A2, once each
	got := map[string]bool{}
	for _, item := range items {
		got[item.StudentID] = true
		require.NotNil(t, item.AssignmentID)
		assert.Equal(t, assignmentID, *item.AssignmentID)
	}
	assert.True(t, got[studentA1.String()])
	assert.True(t, got[studentA2.String()])
	assert.False(t, got[studentB.String()])
}
