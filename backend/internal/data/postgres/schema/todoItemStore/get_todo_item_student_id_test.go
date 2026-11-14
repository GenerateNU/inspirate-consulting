package todoItemRepository

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
)

// test retrieving the student id of an existing todo item
func TestGetTodoItemStudentID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	studentID := uuid.New()
	created, err := repo.CreateTodoItem(ctx, studentID.String(), uuid.NewString(), models.CreateTodoItemRequestBody{
		TodoDescription: "Finish the Common App essay",
	})
	if err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}

	got, err := repo.GetTodoItemStudentID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetTodoItemStudentID failed: %v", err)
	}

	if got != studentID {
		t.Errorf("expected StudentID %v, got %v", studentID, got)
	}
}

// test retrieving the student id of a nonexistent todo item returns a not found error
func TestGetTodoItemStudentID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	got, err := repo.GetTodoItemStudentID(ctx, uuid.NewString())
	if err == nil {
		t.Fatal("expected an error getting the student ID of a nonexistent todo item, got nil")
	}
	if got != uuid.Nil {
		t.Errorf("expected the nil UUID on error, got %v", got)
	}

	var httpErr errs.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected an errs.HTTPError, got %T: %v", err, err)
	}
	if httpErr.GetStatus() != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, httpErr.GetStatus())
	}
}

// test retrieving the student id with a non-uuid todo item id returns an error, since the column is typed uuid
func TestGetTodoItemStudentID_InvalidID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	got, err := repo.GetTodoItemStudentID(ctx, "not-a-uuid")
	if err == nil {
		t.Fatal("expected an error for a non-uuid todo item id, got nil")
	}
	if got != uuid.Nil {
		t.Errorf("expected the nil UUID on error, got %v", got)
	}
}
