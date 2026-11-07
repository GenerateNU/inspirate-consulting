package todoItemRepository

import (
	"context"
	"testing"
	"time"

	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
)

// test marking a todo item complete with a valid id
func TestUpdateTodoItemCompletedAt(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	deadline := time.Now().Add(3 * 24 * time.Hour).UTC().Truncate(time.Second)

	input := models.CreateTodoItemRequestBody{
		TodoDescription: "Submit the FAFSA",
		Deadline:        &deadline,
	}
	created, err := repo.CreateTodoItem(ctx, uuid.NewString(), uuid.NewString(), input)
	if err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}

	completedAt := time.Now().UTC().Truncate(time.Second)

	updated, err := repo.UpdateTodoItemCompletedAt(ctx, created.ID, &completedAt)
	if err != nil {
		t.Fatalf("UpdateTodoItemCompletedAt failed: %v", err)
	}

	if updated.CompletedAt == nil {
		t.Fatal("expected non-nil CompletedAt")
	}
	if !updated.CompletedAt.Equal(completedAt) {
		t.Errorf("expected CompletedAt %v, got %v", completedAt, *updated.CompletedAt)
	}
	if updated.UpdatedAt.Before(created.UpdatedAt) {
		t.Errorf("expected UpdatedAt to advance, got %v (was %v)", updated.UpdatedAt, created.UpdatedAt)
	}

	// Everything other than completed_at/updated_at should be left alone.
	if updated.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, updated.ID)
	}
	if updated.StudentID != created.StudentID {
		t.Errorf("expected StudentID %q, got %q", created.StudentID, updated.StudentID)
	}
	if updated.UserID != created.UserID {
		t.Errorf("expected UserID %q, got %q", created.UserID, updated.UserID)
	}
	if updated.TodoDescription != created.TodoDescription {
		t.Errorf("expected TodoDescription %q, got %q", created.TodoDescription, updated.TodoDescription)
	}
	if !updated.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("expected CreatedAt %v, got %v", created.CreatedAt, updated.CreatedAt)
	}
	if updated.Deadline == nil || !updated.Deadline.Equal(deadline) {
		t.Errorf("expected Deadline %v, got %v", deadline, updated.Deadline)
	}

	// The change should be visible on a subsequent read.
	items, err := repo.ListTodoItems(ctx, models.TodoItemQuery{StudentID: created.StudentID})
	if err != nil {
		t.Fatalf("ListTodoItems failed: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 todo item for student %s, got %d", created.StudentID, len(items))
	}
	if items[0].CompletedAt == nil || !items[0].CompletedAt.Equal(completedAt) {
		t.Errorf("expected persisted CompletedAt %v, got %v", completedAt, items[0].CompletedAt)
	}
}

// test clearing completed_at, which un-completes a todo item
func TestUpdateTodoItemCompletedAt_Clear(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	input := models.CreateTodoItemRequestBody{
		TodoDescription: "Completed by mistake",
	}
	created, err := repo.CreateTodoItem(ctx, uuid.NewString(), uuid.NewString(), input)
	if err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}

	completedAt := time.Now().UTC().Truncate(time.Second)
	if _, err := repo.UpdateTodoItemCompletedAt(ctx, created.ID, &completedAt); err != nil {
		t.Fatalf("setup UpdateTodoItemCompletedAt failed: %v", err)
	}

	cleared, err := repo.UpdateTodoItemCompletedAt(ctx, created.ID, nil)
	if err != nil {
		t.Fatalf("expected clearing CompletedAt to succeed, got error: %v", err)
	}
	if cleared.CompletedAt != nil {
		t.Errorf("expected nil CompletedAt, got %v", cleared.CompletedAt)
	}
}

// test updating a todo item with a nonexistent id returns an error
func TestUpdateTodoItemCompletedAt_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	completedAt := time.Now().UTC().Truncate(time.Second)

	_, err := repo.UpdateTodoItemCompletedAt(ctx, uuid.NewString(), &completedAt)
	if err == nil {
		t.Fatal("expected an error for a nonexistent id, got nil")
	}
}

// test updating a todo item's description with a valid id
func TestUpdateTodoItemDescription(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	deadline := time.Now().Add(3 * 24 * time.Hour).UTC().Truncate(time.Second)

	input := models.CreateTodoItemRequestBody{
		TodoDescription: "Submit the FAFSA",
		Deadline:        &deadline,
	}
	created, err := repo.CreateTodoItem(ctx, uuid.NewString(), uuid.NewString(), input)
	if err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}

	newDescription := "Submit the FAFSA and the CSS Profile"

	updated, err := repo.UpdateTodoItemDescription(ctx, created.ID, newDescription)
	if err != nil {
		t.Fatalf("UpdateTodoItemDescription failed: %v", err)
	}

	if updated.TodoDescription != newDescription {
		t.Errorf("expected TodoDescription %q, got %q", newDescription, updated.TodoDescription)
	}
	if updated.UpdatedAt.Before(created.UpdatedAt) {
		t.Errorf("expected UpdatedAt to advance, got %v (was %v)", updated.UpdatedAt, created.UpdatedAt)
	}

	// Everything other than todo_description/updated_at should be left alone.
	if updated.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, updated.ID)
	}
	if updated.StudentID != created.StudentID {
		t.Errorf("expected StudentID %q, got %q", created.StudentID, updated.StudentID)
	}
	if updated.UserID != created.UserID {
		t.Errorf("expected UserID %q, got %q", created.UserID, updated.UserID)
	}
	if updated.CompletedAt != nil {
		t.Errorf("expected nil CompletedAt, got %v", updated.CompletedAt)
	}
	if !updated.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("expected CreatedAt %v, got %v", created.CreatedAt, updated.CreatedAt)
	}
	if updated.Deadline == nil || !updated.Deadline.Equal(deadline) {
		t.Errorf("expected Deadline %v, got %v", deadline, updated.Deadline)
	}

	// The change should be visible on a subsequent read.
	items, err := repo.ListTodoItems(ctx, models.TodoItemQuery{StudentID: created.StudentID})
	if err != nil {
		t.Fatalf("ListTodoItems failed: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 todo item for student %s, got %d", created.StudentID, len(items))
	}
	if items[0].TodoDescription != newDescription {
		t.Errorf("expected persisted TodoDescription %q, got %q", newDescription, items[0].TodoDescription)
	}
}

// test updating a todo item's description with a nonexistent id returns an error
func TestUpdateTodoItemDescription_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	_, err := repo.UpdateTodoItemDescription(ctx, uuid.NewString(), "Submit the FAFSA")
	if err == nil {
		t.Fatal("expected an error for a nonexistent id, got nil")
	}
}

// test changing a todo item's deadline with a valid id
func TestUpdateTodoItemDeadline(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	deadline := time.Now().Add(3 * 24 * time.Hour).UTC().Truncate(time.Second)

	input := models.CreateTodoItemRequestBody{
		TodoDescription: "Submit the FAFSA",
		Deadline:        &deadline,
	}
	created, err := repo.CreateTodoItem(ctx, uuid.NewString(), uuid.NewString(), input)
	if err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}

	newDeadline := time.Now().Add(10 * 24 * time.Hour).UTC().Truncate(time.Second)

	updated, err := repo.UpdateTodoItemDeadline(ctx, created.ID, &newDeadline)
	if err != nil {
		t.Fatalf("UpdateTodoItemDeadline failed: %v", err)
	}

	if updated.Deadline == nil {
		t.Fatal("expected non-nil Deadline")
	}
	if !updated.Deadline.Equal(newDeadline) {
		t.Errorf("expected Deadline %v, got %v", newDeadline, *updated.Deadline)
	}
	if updated.UpdatedAt.Before(created.UpdatedAt) {
		t.Errorf("expected UpdatedAt to advance, got %v (was %v)", updated.UpdatedAt, created.UpdatedAt)
	}

	// Everything other than deadline/updated_at should be left alone.
	if updated.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, updated.ID)
	}
	if updated.StudentID != created.StudentID {
		t.Errorf("expected StudentID %q, got %q", created.StudentID, updated.StudentID)
	}
	if updated.UserID != created.UserID {
		t.Errorf("expected UserID %q, got %q", created.UserID, updated.UserID)
	}
	if updated.TodoDescription != created.TodoDescription {
		t.Errorf("expected TodoDescription %q, got %q", created.TodoDescription, updated.TodoDescription)
	}
	if updated.CompletedAt != nil {
		t.Errorf("expected nil CompletedAt, got %v", updated.CompletedAt)
	}
	if !updated.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("expected CreatedAt %v, got %v", created.CreatedAt, updated.CreatedAt)
	}

	// The change should be visible on a subsequent read.
	items, err := repo.ListTodoItems(ctx, models.TodoItemQuery{StudentID: created.StudentID})
	if err != nil {
		t.Fatalf("ListTodoItems failed: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 todo item for student %s, got %d", created.StudentID, len(items))
	}
	if items[0].Deadline == nil || !items[0].Deadline.Equal(newDeadline) {
		t.Errorf("expected persisted Deadline %v, got %v", newDeadline, items[0].Deadline)
	}
}

// test clearing a todo item's deadline
func TestUpdateTodoItemDeadline_Clear(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	deadline := time.Now().Add(3 * 24 * time.Hour).UTC().Truncate(time.Second)

	input := models.CreateTodoItemRequestBody{
		TodoDescription: "Deadline set by mistake",
		Deadline:        &deadline,
	}
	created, err := repo.CreateTodoItem(ctx, uuid.NewString(), uuid.NewString(), input)
	if err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}

	cleared, err := repo.UpdateTodoItemDeadline(ctx, created.ID, nil)
	if err != nil {
		t.Fatalf("expected clearing Deadline to succeed, got error: %v", err)
	}
	if cleared.Deadline != nil {
		t.Errorf("expected nil Deadline, got %v", cleared.Deadline)
	}
}

// test changing a todo item's deadline with a nonexistent id returns an error
func TestUpdateTodoItemDeadline_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	deadline := time.Now().Add(3 * 24 * time.Hour).UTC().Truncate(time.Second)

	_, err := repo.UpdateTodoItemDeadline(ctx, uuid.NewString(), &deadline)
	if err == nil {
		t.Fatal("expected an error for a nonexistent id, got nil")
	}
}
