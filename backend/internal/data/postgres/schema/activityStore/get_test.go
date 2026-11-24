package activityRepository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"inspirate-consulting/internal/models"
)

// test that a counselor only sees activity belonging to their own students,
// not another counselor's
func TestGetActivitiesForCounselor_ScopedToOwnStudents(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorA := createTestCounselor(t, db)
	studentA := createTestStudent(t, db, counselorA)
	essayA := createTestEssay(t, db, studentA)
	statusA := models.Review
	if err := repo.CreateActivity(ctx, studentA, models.CreateActivityBody{
		Action: models.ActionAddedSchool, EssayID: &essayA, EssayStatus: &statusA,
	}); err != nil {
		t.Fatalf("failed to create activity for studentA: %v", err)
	}

	counselorB := createTestCounselor(t, db)
	studentB := createTestStudent(t, db, counselorB)
	essayB := createTestEssay(t, db, studentB)
	statusB := models.Review
	if err := repo.CreateActivity(ctx, studentB, models.CreateActivityBody{
		Action: models.ActionRemovedSchool, EssayID: &essayB, EssayStatus: &statusB,
	}); err != nil {
		t.Fatalf("failed to create activity for studentB: %v", err)
	}

	activities, err := repo.GetActivitiesForCounselor(ctx, counselorA, nil)
	if err != nil {
		t.Fatalf("GetActivitiesForCounselor failed: %v", err)
	}

	if len(activities) != 1 {
		t.Fatalf("expected 1 activity for counselorA, got %d", len(activities))
	}
	if activities[0].StudentID != studentA {
		t.Errorf("expected activity for studentA (%s), got %s", studentA, activities[0].StudentID)
	}
	if activities[0].Action != models.ActionAddedSchool {
		t.Errorf("expected action %q, got %q", models.ActionAddedSchool, activities[0].Action)
	}
}

// test that activity across a counselor's multiple students is sorted newest first
func TestGetActivitiesForCounselor_SortedByCreatedAtDesc(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	studentOne := createTestStudent(t, db, counselorID)
	studentTwo := createTestStudent(t, db, counselorID)

	older := time.Now().Add(-1 * time.Hour)
	newer := time.Now()

	insertWithTimestamp := func(studentID uuid.UUID, action models.Action, createdAt time.Time) {
		t.Helper()
		if _, err := db.Exec(ctx,
			`INSERT INTO public.activity (student_id, action, created_at) VALUES ($1, $2, $3)`,
			studentID, action, createdAt,
		); err != nil {
			t.Fatalf("failed to seed activity: %v", err)
		}
	}

	insertWithTimestamp(studentOne, models.ActionAddedSchool, older)
	insertWithTimestamp(studentTwo, models.ActionTaskPastDue, newer)

	activities, err := repo.GetActivitiesForCounselor(ctx, counselorID, nil)
	if err != nil {
		t.Fatalf("GetActivitiesForCounselor failed: %v", err)
	}

	if len(activities) != 2 {
		t.Fatalf("expected 2 activities, got %d", len(activities))
	}
	if activities[0].Action != models.ActionTaskPastDue {
		t.Errorf("expected the newest activity (task_past_due) first, got %q", activities[0].Action)
	}
	if activities[1].Action != models.ActionAddedSchool {
		t.Errorf("expected the oldest activity (added_school) second, got %q", activities[1].Action)
	}
}

// test that a counselor with students but no recorded activity gets an empty slice, not an error
func TestGetActivitiesForCounselor_NoActivity(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	createTestStudent(t, db, counselorID)

	activities, err := repo.GetActivitiesForCounselor(ctx, counselorID, nil)
	if err != nil {
		t.Fatalf("GetActivitiesForCounselor failed: %v", err)
	}
	if len(activities) != 0 {
		t.Errorf("expected no activities, got %d", len(activities))
	}
}

// test that the studentID filter scopes results to just that student
func TestGetActivitiesForCounselor_FilteredByStudentID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	studentOne := createTestStudent(t, db, counselorID)
	studentTwo := createTestStudent(t, db, counselorID)

	essayOne := createTestEssay(t, db, studentOne)
	statusOne := models.Review
	if err := repo.CreateActivity(ctx, studentOne, models.CreateActivityBody{
		Action: models.ActionAddedSchool, EssayID: &essayOne, EssayStatus: &statusOne,
	}); err != nil {
		t.Fatalf("failed to create activity for studentOne: %v", err)
	}
	essayTwo := createTestEssay(t, db, studentTwo)
	statusTwo := models.Review
	if err := repo.CreateActivity(ctx, studentTwo, models.CreateActivityBody{
		Action: models.ActionTaskPastDue, EssayID: &essayTwo, EssayStatus: &statusTwo,
	}); err != nil {
		t.Fatalf("failed to create activity for studentTwo: %v", err)
	}

	activities, err := repo.GetActivitiesForCounselor(ctx, counselorID, &studentOne)
	if err != nil {
		t.Fatalf("GetActivitiesForCounselor failed: %v", err)
	}

	if len(activities) != 1 {
		t.Fatalf("expected 1 activity scoped to studentOne, got %d", len(activities))
	}
	if activities[0].StudentID != studentOne {
		t.Errorf("expected activity for studentOne (%s), got %s", studentOne, activities[0].StudentID)
	}
}

// test the access-control case: passing a studentID that belongs to a
// different counselor must return nothing, not that counselor's data
func TestGetActivitiesForCounselor_StudentIDFromAnotherCounselorReturnsNothing(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorA := createTestCounselor(t, db)
	createTestStudent(t, db, counselorA)

	counselorB := createTestCounselor(t, db)
	studentB := createTestStudent(t, db, counselorB)
	essayB := createTestEssay(t, db, studentB)
	statusB := models.Review
	if err := repo.CreateActivity(ctx, studentB, models.CreateActivityBody{
		Action: models.ActionAddedSchool, EssayID: &essayB, EssayStatus: &statusB,
	}); err != nil {
		t.Fatalf("failed to create activity for studentB: %v", err)
	}

	// counselorA asks for studentB's activity - studentB isn't theirs
	activities, err := repo.GetActivitiesForCounselor(ctx, counselorA, &studentB)
	if err != nil {
		t.Fatalf("GetActivitiesForCounselor failed: %v", err)
	}

	if len(activities) != 0 {
		t.Errorf("expected no activities when requesting another counselor's student, got %d", len(activities))
	}
}
