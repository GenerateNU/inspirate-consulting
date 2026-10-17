package counselorRepository

import (
	"context"
	"testing"
	"time"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
)

// test list function returns each of the counselor's students with an overdue incomplete to-do item, once
func TestListFallingBehindStudentIDs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewCounselorRepository(db)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	otherCounselorID := createTestCounselor(t, db)

	past := time.Now().Add(-24 * time.Hour)
	future := time.Now().Add(24 * time.Hour)

	behind := createTestStudent(t, db, counselorID, "Behind", "junior")
	createTestTodoItem(t, db, behind, &past, nil)
	createTestTodoItem(t, db, behind, &past, nil)

	// overdue but completed, not yet due, and no deadline should not count
	completed := createTestStudent(t, db, counselorID, "Completed", "junior")
	createTestTodoItem(t, db, completed, &past, &past)
	notDue := createTestStudent(t, db, counselorID, "Not Due", "junior")
	createTestTodoItem(t, db, notDue, &future, nil)
	noDeadline := createTestStudent(t, db, counselorID, "No Deadline", "junior")
	createTestTodoItem(t, db, noDeadline, nil, nil)

	// another counselor's student who is behind should not show up
	otherStudent := createTestStudent(t, db, otherCounselorID, "Other", "senior")
	createTestTodoItem(t, db, otherStudent, &past, nil)

	results, err := repo.ListFallingBehindStudentIDs(ctx, counselorID.String())
	if err != nil {
		t.Fatalf("ListFallingBehindStudentIDs failed: %v", err)
	}

	if len(results) != 1 || results[0] != behind {
		t.Errorf("expected only %v, got %v", behind, results)
	}
}
