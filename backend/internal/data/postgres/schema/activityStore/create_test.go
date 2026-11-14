package activityRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"inspirate-consulting/internal/models"
)

// test creating an activity with a todo item reference, no optional essay fields
func TestCreateActivity(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	studentID := createTestStudent(t, db, counselorID)
	todoItemID := createTestTodoItem(t, db, studentID)

	err := repo.CreateActivity(ctx, studentID, models.CreateActivityBody{
		Action:     models.ActionTaskPastDue,
		TodoItemID: &todoItemID,
	})
	if err != nil {
		t.Fatalf("CreateActivity failed: %v", err)
	}

	var count int
	if err := db.QueryRow(ctx,
		`SELECT count(*) FROM public.activity WHERE student_id = $1 AND action = $2`,
		studentID, models.ActionTaskPastDue,
	).Scan(&count); err != nil {
		t.Fatalf("failed to verify inserted row: %v", err)
	}
	if count != 1 {
		t.Errorf("expected exactly 1 matching activity row, got %d", count)
	}
}

// test that an activity with no reference set at all is rejected
func TestCreateActivity_NoReferenceFails(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	studentID := createTestStudent(t, db, counselorID)

	err := repo.CreateActivity(ctx, studentID, models.CreateActivityBody{Action: models.ActionTaskPastDue})
	if err == nil {
		t.Fatal("expected an error creating an activity with no reference set, got nil")
	}
}

// test that an activity with more than one reference set is rejected
func TestCreateActivity_MultipleReferencesFails(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	studentID := createTestStudent(t, db, counselorID)
	todoItemID := createTestTodoItem(t, db, studentID)
	essayID := createTestEssay(t, db, studentID)
	status := models.Review

	err := repo.CreateActivity(ctx, studentID, models.CreateActivityBody{
		Action:      models.ActionEssayStatusChanged,
		EssayID:     &essayID,
		EssayStatus: &status,
		TodoItemID:  &todoItemID,
	})
	if err == nil {
		t.Fatal("expected an error creating an activity with more than one reference set, got nil")
	}
}

// test that essay_id without essay_status (or vice versa) is rejected
func TestCreateActivity_EssayIDWithoutStatusFails(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	studentID := createTestStudent(t, db, counselorID)
	essayID := createTestEssay(t, db, studentID)

	err := repo.CreateActivity(ctx, studentID, models.CreateActivityBody{
		Action:  models.ActionEssayStatusChanged,
		EssayID: &essayID,
	})
	if err == nil {
		t.Fatal("expected an error creating an activity with essay_id but no essay_status, got nil")
	}
}

// test creating an activity that references an essay and its new status
func TestCreateActivity_WithEssayStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	studentID := createTestStudent(t, db, counselorID)
	essayID := createTestEssay(t, db, studentID)
	status := models.Review

	err := repo.CreateActivity(ctx, studentID, models.CreateActivityBody{
		Action:      models.ActionEssayStatusChanged,
		EssayID:     &essayID,
		EssayStatus: &status,
	})
	if err != nil {
		t.Fatalf("CreateActivity failed: %v", err)
	}

	var gotEssayID uuid.UUID
	var gotStatus models.Status
	if err := db.QueryRow(ctx,
		`SELECT essay_id, essay_status FROM public.activity WHERE student_id = $1 AND action = $2`,
		studentID, models.ActionEssayStatusChanged,
	).Scan(&gotEssayID, &gotStatus); err != nil {
		t.Fatalf("failed to read back inserted row: %v", err)
	}
	if gotEssayID != essayID {
		t.Errorf("expected essay_id %s, got %s", essayID, gotEssayID)
	}
	if gotStatus != status {
		t.Errorf("expected essay_status %q, got %q", status, gotStatus)
	}
}

// test that creating an activity for a student that doesn't exist fails
func TestCreateActivity_UnknownStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()
	todoItemID := uuid.New()

	err := repo.CreateActivity(ctx, uuid.New(), models.CreateActivityBody{
		Action:     models.ActionTaskPastDue,
		TodoItemID: &todoItemID,
	})
	if err == nil {
		t.Fatal("expected an error creating an activity for an unknown student, got nil")
	}
}

// test that referencing an essay that doesn't exist fails
func TestCreateActivity_UnknownEssay(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	studentID := createTestStudent(t, db, counselorID)
	unknownEssayID := uuid.New()
	status := models.Review

	err := repo.CreateActivity(ctx, studentID, models.CreateActivityBody{
		Action:      models.ActionEssayStatusChanged,
		EssayID:     &unknownEssayID,
		EssayStatus: &status,
	})
	if err == nil {
		t.Fatal("expected an error creating an activity with an unknown essay_id, got nil")
	}
}
