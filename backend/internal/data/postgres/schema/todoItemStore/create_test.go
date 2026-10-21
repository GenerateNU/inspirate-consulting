package todoItemRepository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"inspirate-consulting/internal/models"
)

// test normal todo item creation with valid input
func TestCreateTodoItem(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	deadline := time.Now().Add(7 * 24 * time.Hour).UTC().Truncate(time.Second)

	studentID := uuid.NewString()
	userID := uuid.NewString()
	input := models.CreateTodoItemRequestBody{
		TodoDescription: "Finish the Common App essay",
		Deadline:        &deadline,
	}

	created, err := repo.CreateTodoItem(ctx, studentID, userID, input)
	if err != nil {
		t.Fatalf("CreateTodoItem failed: %v", err)
	}

	if created.ID == "" {
		t.Error("expected a generated ID, got empty string")
	}
	if created.StudentID != studentID {
		t.Errorf("expected StudentID %q, got %q", studentID, created.StudentID)
	}
	if created.UserID != userID {
		t.Errorf("expected UserID %q, got %q", userID, created.UserID)
	}
	if created.TodoDescription != input.TodoDescription {
		t.Errorf("expected TodoDescription %q, got %q", input.TodoDescription, created.TodoDescription)
	}
	if created.CompletedAt != nil {
		t.Errorf("expected nil CompletedAt on a new item, got %v", created.CompletedAt)
	}
	if created.Deadline == nil {
		t.Fatal("expected non-nil Deadline")
	}
	if !created.Deadline.Equal(deadline) {
		t.Errorf("expected Deadline %v, got %v", deadline, *created.Deadline)
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if created.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

// test creation with no deadline, since deadline is a nullable column
func TestCreateTodoItem_NilDeadline(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	input := models.CreateTodoItemRequestBody{
		TodoDescription: "Someday task with no deadline",
	}

	created, err := repo.CreateTodoItem(ctx, uuid.NewString(), uuid.NewString(), input)
	if err != nil {
		t.Fatalf("expected create with a nil deadline to succeed, got error: %v", err)
	}

	if created.Deadline != nil {
		t.Errorf("expected nil Deadline, got %v", created.Deadline)
	}
	if created.CompletedAt != nil {
		t.Errorf("expected nil CompletedAt, got %v", created.CompletedAt)
	}
}

// test creation with a non-uuid student_id fails, since the column is typed uuid
func TestCreateTodoItem_InvalidStudentID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	input := models.CreateTodoItemRequestBody{
		TodoDescription: "Should never be inserted",
	}

	_, err := repo.CreateTodoItem(ctx, "not-a-uuid", uuid.NewString(), input)
	if err == nil {
		t.Fatal("expected an error creating a todo item with a non-uuid student_id, got nil")
	}
}
