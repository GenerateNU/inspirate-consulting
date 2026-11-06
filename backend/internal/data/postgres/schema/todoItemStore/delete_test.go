package todoItemRepository

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// test deleting a todo item with a valid id
func TestDeleteTodoItem(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	input := models.CreateTodoItemRequestBody{
		TodoDescription: "Submit the FAFSA",
	}
	created, err := repo.CreateTodoItem(ctx, uuid.NewString(), uuid.NewString(), input)
	if err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}

	if err := repo.DeleteTodoItem(ctx, created.ID); err != nil {
		t.Fatalf("DeleteTodoItem failed: %v", err)
	}

	// The item should no longer be visible on a subsequent read.
	items, err := repo.ListTodoItems(ctx, models.TodoItemQuery{StudentID: created.StudentID})
	if err != nil {
		t.Fatalf("ListTodoItems failed: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 todo items for student %s after delete, got %d", created.StudentID, len(items))
	}
}

// test deleting one todo item leaves the student's other items alone
func TestDeleteTodoItem_LeavesOthers(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	studentID := uuid.NewString()
	userID := uuid.NewString()

	toDelete, err := repo.CreateTodoItem(ctx, studentID, userID, models.CreateTodoItemRequestBody{
		TodoDescription: "Delete me",
	})
	if err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}
	toKeep, err := repo.CreateTodoItem(ctx, studentID, userID, models.CreateTodoItemRequestBody{
		TodoDescription: "Keep me",
	})
	if err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}

	if err := repo.DeleteTodoItem(ctx, toDelete.ID); err != nil {
		t.Fatalf("DeleteTodoItem failed: %v", err)
	}

	items, err := repo.ListTodoItems(ctx, models.TodoItemQuery{StudentID: studentID})
	if err != nil {
		t.Fatalf("ListTodoItems failed: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 todo item for student %s after delete, got %d", studentID, len(items))
	}
	if items[0].ID != toKeep.ID {
		t.Errorf("expected remaining item %s, got %s", toKeep.ID, items[0].ID)
	}
}

// test deleting a todo item with a nonexistent id returns a not found error
func TestDeleteTodoItem_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	err := repo.DeleteTodoItem(ctx, uuid.NewString())
	if err == nil {
		t.Fatal("expected an error for a nonexistent id, got nil")
	}

	var httpErr errs.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected errs.HTTPError, got %T: %v", err, err)
	}
	if httpErr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, httpErr.Code)
	}
}

// test deleting the same todo item twice returns an error the second time
func TestDeleteTodoItem_AlreadyDeleted(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	created, err := repo.CreateTodoItem(ctx, uuid.NewString(), uuid.NewString(), models.CreateTodoItemRequestBody{
		TodoDescription: "Delete me twice",
	})
	if err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}

	if err := repo.DeleteTodoItem(ctx, created.ID); err != nil {
		t.Fatalf("first DeleteTodoItem failed: %v", err)
	}

	if err := repo.DeleteTodoItem(ctx, created.ID); err == nil {
		t.Fatal("expected an error deleting an already-deleted item, got nil")
	}
}
