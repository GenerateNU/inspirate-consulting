package todoItemRepository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"inspirate-consulting/internal/models"
)

// test marking a todo item complete with a valid id
func TestUpdateTodoItemCompletedAt(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	deadline := time.Now().Add(3 * 24 * time.Hour).UTC().Truncate(time.Second)

	input := &models.TodoItem{
		StudentID:       uuid.NewString(),
		UserID:          uuid.NewString(),
		TodoDescription: "Submit the FAFSA",
		Deadline:        &deadline,
	}
	created, err := repo.CreateTodoItem(ctx, input)
	if err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}
	cleanupTodoItem(t, db, created.ID)

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
	items, err := repo.GetTodoItemsByStudent(ctx, created.StudentID)
	if err != nil {
		t.Fatalf("GetTodoItemsByStudent failed: %v", err)
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

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	input := &models.TodoItem{
		StudentID:       uuid.NewString(),
		UserID:          uuid.NewString(),
		TodoDescription: "Completed by mistake",
	}
	created, err := repo.CreateTodoItem(ctx, input)
	if err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}
	cleanupTodoItem(t, db, created.ID)

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
