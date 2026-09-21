package todoItemRepository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"inspirate-consulting/internal/models"
)

// test retrieval returns every todo item belonging to a student, and only that student's items
func TestGetTodoItemsByStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := uuid.NewString()
	userID := uuid.NewString()
	otherStudentID := uuid.NewString()

	deadline := time.Now().Add(14 * 24 * time.Hour).UTC().Truncate(time.Second)

	inputs := []*models.CreateTodoItemRequestBody{
		{
			StudentID:       studentID,
			UserID:          userID,
			TodoDescription: "Request a letter of recommendation",
			Deadline:        &deadline,
		},
		{
			StudentID:       studentID,
			UserID:          userID,
			TodoDescription: "Schedule a campus tour",
		},
		{
			StudentID:       otherStudentID,
			UserID:          userID,
			TodoDescription: "Belongs to a different student",
		},
	}

	created := make([]models.TodoItem, 0, len(inputs))
	for _, input := range inputs {
		output, err := repo.CreateTodoItem(ctx, input)
		if err != nil {
			t.Fatalf("setup CreateTodoItem(%q) failed: %v", input.TodoDescription, err)
		}
		cleanupTodoItem(t, db, output.ID)
		created = append(created, *output)
	}

	items, err := repo.GetTodoItemsByStudent(ctx, studentID)
	if err != nil {
		t.Fatalf("GetTodoItemsByStudent failed: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 todo items for student %s, got %d", studentID, len(items))
	}

	byID := make(map[string]models.TodoItem, len(items))
	for _, item := range items {
		byID[item.ID] = item
		if item.StudentID != studentID {
			t.Errorf("id=%s: expected StudentID %q, got %q", item.ID, studentID, item.StudentID)
		}
	}

	// The other student's item must not leak into these results.
	if _, ok := byID[created[2].ID]; ok {
		t.Errorf("expected item %s belonging to another student to be excluded", created[2].ID)
	}

	for _, want := range created[:2] {
		got, ok := byID[want.ID]
		if !ok {
			t.Errorf("expected todo item %q (id=%s) to appear in GetTodoItemsByStudent result", want.TodoDescription, want.ID)
			continue
		}
		if got.TodoDescription != want.TodoDescription {
			t.Errorf("id=%s: expected TodoDescription %q, got %q", want.ID, want.TodoDescription, got.TodoDescription)
		}
		if got.UserID != want.UserID {
			t.Errorf("id=%s: expected UserID %q, got %q", want.ID, want.UserID, got.UserID)
		}
		if got.CompletedAt != nil {
			t.Errorf("id=%s: expected nil CompletedAt, got %v", want.ID, got.CompletedAt)
		}
	}

	// Spot-check that the nullable deadline round-tripped per item.
	withDeadline := byID[created[0].ID]
	if withDeadline.Deadline == nil || !withDeadline.Deadline.Equal(deadline) {
		t.Errorf("expected Deadline %v, got %v", deadline, withDeadline.Deadline)
	}
	withoutDeadline := byID[created[1].ID]
	if withoutDeadline.Deadline != nil {
		t.Errorf("expected nil Deadline, got %v", withoutDeadline.Deadline)
	}
}

// test retrieval for a student with no todo items returns an empty result, not an error
func TestGetTodoItemsByStudent_NoItems(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	items, err := repo.GetTodoItemsByStudent(ctx, uuid.NewString())
	if err != nil {
		t.Fatalf("expected no error for a student with no todo items, got: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 todo items, got %d", len(items))
	}
}

// test retrieval with a non-uuid student id returns an error, since the column is typed uuid
func TestGetTodoItemsByStudent_InvalidID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	_, err := repo.GetTodoItemsByStudent(ctx, "not-a-uuid")
	if err == nil {
		t.Fatal("expected an error for a non-uuid student id, got nil")
	}
}
