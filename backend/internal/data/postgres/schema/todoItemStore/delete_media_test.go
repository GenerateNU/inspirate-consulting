package todoItemRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"inspirate-consulting/internal/models"
)

// test deleting the only todo for a media revokes the access it granted
func TestDeleteTodoItem_RevokesMediaAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := uuid.NewString()
	mediaID := insertTestMedia(t, db)

	created, err := repo.CreateTodoItem(ctx, studentID, uuid.NewString(), models.CreateTodoItemRequestBody{
		TodoDescription: "Watch the video",
		MediaID:         &mediaID,
	})
	if err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}
	if got := countMediaAccess(t, db, studentID, mediaID); got != 1 {
		t.Fatalf("setup: expected 1 media_access row, got %d", got)
	}

	if err := repo.DeleteTodoItem(ctx, created.ID); err != nil {
		t.Fatalf("DeleteTodoItem failed: %v", err)
	}

	if got := countMediaAccess(t, db, studentID, mediaID); got != 0 {
		t.Errorf("expected access to be revoked, got %d media_access rows", got)
	}
}

// test deleting one student's todo doesn't touch another student's access to the same media
func TestDeleteTodoItem_OnlyRevokesThatStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := uuid.NewString()
	otherStudentID := uuid.NewString()
	mediaID := insertTestMedia(t, db)

	todo, err := repo.CreateTodoItem(ctx, studentID, uuid.NewString(), models.CreateTodoItemRequestBody{
		TodoDescription: "Watch the video",
		MediaID:         &mediaID,
	})
	if err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}
	if _, err := repo.CreateTodoItem(ctx, otherStudentID, uuid.NewString(), models.CreateTodoItemRequestBody{
		TodoDescription: "Watch the video",
		MediaID:         &mediaID,
	}); err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}

	if err := repo.DeleteTodoItem(ctx, todo.ID); err != nil {
		t.Fatalf("DeleteTodoItem failed: %v", err)
	}

	if got := countMediaAccess(t, db, studentID, mediaID); got != 0 {
		t.Errorf("expected the student's access to be revoked, got %d rows", got)
	}
	if got := countMediaAccess(t, db, otherStudentID, mediaID); got != 1 {
		t.Errorf("expected the other student's access to remain, got %d rows", got)
	}
}

// test deleting a todo without media leaves the student's media access alone
func TestDeleteTodoItem_NoMediaLeavesAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := uuid.NewString()
	mediaID := insertTestMedia(t, db)

	if _, err := repo.CreateTodoItem(ctx, studentID, uuid.NewString(), models.CreateTodoItemRequestBody{
		TodoDescription: "Watch the video",
		MediaID:         &mediaID,
	}); err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}
	essayTodo, err := repo.CreateTodoItem(ctx, studentID, uuid.NewString(), models.CreateTodoItemRequestBody{
		TodoDescription: "Draft essay",
	})
	if err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}

	if err := repo.DeleteTodoItem(ctx, essayTodo.ID); err != nil {
		t.Fatalf("DeleteTodoItem failed: %v", err)
	}
	if got := countMediaAccess(t, db, studentID, mediaID); got != 1 {
		t.Errorf("expected media access to remain, got %d rows", got)
	}
}
