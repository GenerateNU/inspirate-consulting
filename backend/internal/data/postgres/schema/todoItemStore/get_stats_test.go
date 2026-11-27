package todoItemRepository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// currentWeekStart asks Postgres for the week start the stats query uses, so tests can't drift from it.
func currentWeekStart(t *testing.T, db *pgxpool.Pool) time.Time {
	t.Helper()
	var weekStart time.Time
	err := db.QueryRow(context.Background(),
		`SELECT date_trunc('week', now() AT TIME ZONE 'UTC') AT TIME ZONE 'UTC'`,
	).Scan(&weekStart)
	require.NoError(t, err)
	return weekStart
}

// insertStatsTodo inserts a to-do item with the given deadline and completion time.
func insertStatsTodo(t *testing.T, db *pgxpool.Pool, studentID, userID uuid.UUID, deadline, completedAt *time.Time) {
	t.Helper()
	_, err := db.Exec(context.Background(),
		`INSERT INTO public.todo_items (student_id, user_id, todo_description, title, deadline, completed_at)
		 VALUES ($1, $2, 'stats task', 'stats task', $3, $4)`,
		studentID, userID, deadline, completedAt,
	)
	require.NoError(t, err)
}

func TestGetTodoItemStats(t *testing.T) {
	t.Parallel()

	day := 24 * time.Hour

	t.Run("counts total, incomplete, and due this week", func(t *testing.T) {
		t.Parallel()
		repo, db := setupTestRepo(t)
		counselorID, counselorUserID := createTestCounselor(t, db)
		studentID := createTestStudentWithCounselor(t, db, counselorID, 2027)

		weekStart := currentWeekStart(t, db)
		at := func(d time.Duration) *time.Time { v := weekStart.Add(d); return &v }
		completed := at(-time.Hour)

		insertStatsTodo(t, db, studentID, counselorUserID, at(0), nil)                  // due at the start of this week
		insertStatsTodo(t, db, studentID, counselorUserID, at(3*day), completed)        // due this week, completed
		insertStatsTodo(t, db, studentID, counselorUserID, at(7*day-time.Second), nil)  // due at the last second of this week
		insertStatsTodo(t, db, studentID, counselorUserID, at(7*day), nil)              // due next week
		insertStatsTodo(t, db, studentID, counselorUserID, at(-time.Second), completed) // due last week, completed
		insertStatsTodo(t, db, studentID, counselorUserID, nil, nil)                    // no deadline

		stats, err := repo.GetTodoItemStats(context.Background(), studentID.String())
		require.NoError(t, err)
		assert.Equal(t, 6, stats.TotalCount)
		assert.Equal(t, 4, stats.IncompleteCount)
		assert.Equal(t, 3, stats.DueThisWeekCount)
	})

	t.Run("student with no tasks", func(t *testing.T) {
		t.Parallel()
		repo, db := setupTestRepo(t)
		counselorID, _ := createTestCounselor(t, db)
		studentID := createTestStudentWithCounselor(t, db, counselorID, 2027)

		stats, err := repo.GetTodoItemStats(context.Background(), studentID.String())
		require.NoError(t, err)
		assert.Equal(t, 0, stats.TotalCount)
		assert.Equal(t, 0, stats.IncompleteCount)
		assert.Equal(t, 0, stats.DueThisWeekCount)
	})

	t.Run("all tasks completed and due this week", func(t *testing.T) {
		t.Parallel()
		repo, db := setupTestRepo(t)
		counselorID, counselorUserID := createTestCounselor(t, db)
		studentID := createTestStudentWithCounselor(t, db, counselorID, 2027)

		weekStart := currentWeekStart(t, db)
		at := func(d time.Duration) *time.Time { v := weekStart.Add(d); return &v }
		completed := at(time.Hour)

		insertStatsTodo(t, db, studentID, counselorUserID, at(day), completed)
		insertStatsTodo(t, db, studentID, counselorUserID, at(4*day), completed)
		insertStatsTodo(t, db, studentID, counselorUserID, at(6*day), completed)

		stats, err := repo.GetTodoItemStats(context.Background(), studentID.String())
		require.NoError(t, err)
		assert.Equal(t, 3, stats.TotalCount)
		assert.Equal(t, 0, stats.IncompleteCount)
		assert.Equal(t, 3, stats.DueThisWeekCount)
	})
}
