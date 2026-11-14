package vectorRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchTodoItems_RanksBySimilarity(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	studentID := uuid.New()

	exact := createTestTodo(t, db, studentID, "Finish Common App essay")
	unrelated := createTestTodo(t, db, studentID, "Buy a graphing calculator")
	linkVector(t, db, unitVector(0), nil, &exact)
	linkVector(t, db, unitVector(1), nil, &unrelated)

	results, err := repo.SearchTodoItems(context.Background(), unitVector(0), []uuid.UUID{studentID}, 10)

	require.NoError(t, err)
	require.Len(t, results, 2)
	assert.Equal(t, exact.String(), results[0].ID)
	assert.Equal(t, unrelated.String(), results[1].ID)
	assert.InDelta(t, 1.0, results[0].Similarity, 0.0001)
	assert.InDelta(t, 0.0, results[1].Similarity, 0.0001)

	// the embedded TodoItem fields are filled in too
	assert.Equal(t, "Finish Common App essay", results[0].TodoDescription)
	assert.Equal(t, studentID.String(), results[0].StudentID)
}

func TestSearchTodoItems_OnlyGivenStudents(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	studentA := uuid.New()
	studentB := uuid.New()
	otherStudent := uuid.New()

	todoA := createTestTodo(t, db, studentA, "a")
	todoB := createTestTodo(t, db, studentB, "b")
	todoOther := createTestTodo(t, db, otherStudent, "not visible")
	linkVector(t, db, unitVector(1), nil, &todoA)
	linkVector(t, db, unitVector(2), nil, &todoB)
	linkVector(t, db, unitVector(0), nil, &todoOther) // a better match, but not one of the given students

	results, err := repo.SearchTodoItems(context.Background(), unitVector(0), []uuid.UUID{studentA, studentB}, 10)

	require.NoError(t, err)
	require.Len(t, results, 2)
	ids := []string{results[0].ID, results[1].ID}
	assert.ElementsMatch(t, []string{todoA.String(), todoB.String()}, ids)
}

func TestSearchTodoItems_NoStudents(t *testing.T) {
	t.Parallel()

	repo, db := setupTestRepo(t)
	todo := createTestTodo(t, db, uuid.New(), "a")
	linkVector(t, db, unitVector(0), nil, &todo)

	results, err := repo.SearchTodoItems(context.Background(), unitVector(0), []uuid.UUID{}, 10)

	require.NoError(t, err)
	assert.Empty(t, results)
}
