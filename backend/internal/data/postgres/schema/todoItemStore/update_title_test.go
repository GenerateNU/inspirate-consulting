package todoItemRepository

import (
	"context"
	"errors"
	"testing"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
)

// test changing a todo item's title with a valid id
func TestUpdateTodoItemTitle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	created, err := repo.CreateTodoItem(ctx, uuid.NewString(), uuid.NewString(), models.CreateTodoItemRequestBody{
		TodoDescription: "Ask the registrar for transcripts",
		Title:           "Old title",
	})
	if err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}

	const newTitle = "Request official transcripts"

	updated, err := repo.UpdateTodoItemTitle(ctx, created.ID, newTitle)
	if err != nil {
		t.Fatalf("UpdateTodoItemTitle failed: %v", err)
	}

	if updated.Title != newTitle {
		t.Errorf("expected Title %q, got %q", newTitle, updated.Title)
	}
	if updated.UpdatedAt.Before(created.UpdatedAt) {
		t.Errorf("expected UpdatedAt to advance, got %v (was %v)", updated.UpdatedAt, created.UpdatedAt)
	}

	// Everything other than title/updated_at should be left alone.
	if updated.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, updated.ID)
	}
	if updated.StudentID != created.StudentID {
		t.Errorf("expected StudentID %q, got %q", created.StudentID, updated.StudentID)
	}
	if updated.TodoDescription != created.TodoDescription {
		t.Errorf("expected TodoDescription %q, got %q", created.TodoDescription, updated.TodoDescription)
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
	if items[0].Title != newTitle {
		t.Errorf("expected persisted Title %q, got %q", newTitle, items[0].Title)
	}
}

// test changing a todo item's title with a nonexistent id returns a not-found error
func TestUpdateTodoItemTitle_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)

	_, err := repo.UpdateTodoItemTitle(context.Background(), uuid.NewString(), "Anything")
	if err == nil {
		t.Fatal("expected an error for a nonexistent id, got nil")
	}
	var httpErr errs.HTTPError
	if !errors.As(err, &httpErr) || httpErr.Code != 404 {
		t.Errorf("expected a 404 not-found error, got %v", err)
	}
}
