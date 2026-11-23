package todoItemRepository

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// moduleLink is the link_type a media todo must have (todo_link_type_matches check).
var moduleLink = models.TodoLinkModule

// insertTestMedia adds a media row directly so todos have a valid media_id to link.
func insertTestMedia(t *testing.T, db *pgxpool.Pool) string {
	t.Helper()
	var id string
	err := db.QueryRow(context.Background(),
		`INSERT INTO public.media (title, description, s3_key) VALUES ('Todo Media', 'test', 'media/todo.mp4') RETURNING id`,
	).Scan(&id)
	if err != nil {
		t.Fatalf("failed to create test media: %v", err)
	}
	return id
}

func countMediaAccess(t *testing.T, db *pgxpool.Pool, studentID string, mediaID string) int {
	t.Helper()
	var n int
	err := db.QueryRow(context.Background(),
		`SELECT count(*) FROM public.media_access WHERE student_id = $1 AND media_id = $2`, studentID, mediaID,
	).Scan(&n)
	if err != nil {
		t.Fatalf("failed to count media access: %v", err)
	}
	return n
}

// test creating a media-linked todo also grants the student access
func TestCreateTodoItem_GrantsMediaAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := uuid.NewString()
	otherStudentID := uuid.NewString()
	mediaID := insertTestMedia(t, db)
	deadline := time.Now().Add(7 * 24 * time.Hour).UTC().Truncate(time.Second)

	created, err := repo.CreateTodoItem(ctx, studentID, uuid.NewString(), models.CreateTodoItemRequestBody{
		TodoDescription: "Watch the video",
		Deadline:        &deadline,
		MediaID:         &mediaID,
		LinkType:        &moduleLink,
	})
	if err != nil {
		t.Fatalf("CreateTodoItem failed: %v", err)
	}
	if created.MediaID == nil || *created.MediaID != mediaID {
		t.Errorf("expected MediaID %q, got %v", mediaID, created.MediaID)
	}
	if got := countMediaAccess(t, db, studentID, mediaID); got != 1 {
		t.Errorf("expected 1 media_access row for the student, got %d", got)
	}
	if got := countMediaAccess(t, db, otherStudentID, mediaID); got != 0 {
		t.Errorf("expected no access for a different student, got %d", got)
	}
}

// test a second todo for media the student already has is rejected with 409, and the
// rejected todo is rolled back so only the first todo exists
func TestCreateTodoItem_ExistingMediaAccessConflicts(t *testing.T) {
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
		LinkType:        &moduleLink,
	}); err != nil {
		t.Fatalf("setup CreateTodoItem failed: %v", err)
	}

	_, err := repo.CreateTodoItem(ctx, studentID, uuid.NewString(), models.CreateTodoItemRequestBody{
		TodoDescription: "Watch the video again",
		MediaID:         &mediaID,
		LinkType:        &moduleLink,
	})
	var httpErr errs.HTTPError
	if !errors.As(err, &httpErr) || httpErr.Code != http.StatusConflict {
		t.Fatalf("expected a 409 conflict, got %v", err)
	}

	var todos int
	if err := db.QueryRow(ctx,
		`SELECT count(*) FROM public.todo_items WHERE student_id = $1 AND media_id = $2`, studentID, mediaID,
	).Scan(&todos); err != nil {
		t.Fatalf("count failed: %v", err)
	}
	if todos != 1 {
		t.Errorf("expected the rejected todo to be rolled back (1 todo), got %d", todos)
	}
	if got := countMediaAccess(t, db, studentID, mediaID); got != 1 {
		t.Errorf("expected exactly 1 media_access row, got %d", got)
	}
}

// test a todo without media grants nothing
func TestCreateTodoItem_NoMediaNoGrant(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := uuid.NewString()
	if _, err := repo.CreateTodoItem(ctx, studentID, uuid.NewString(), models.CreateTodoItemRequestBody{
		TodoDescription: "Draft essay",
	}); err != nil {
		t.Fatalf("CreateTodoItem failed: %v", err)
	}

	var n int
	if err := db.QueryRow(ctx, `SELECT count(*) FROM public.media_access WHERE student_id = $1`, studentID).Scan(&n); err != nil {
		t.Fatalf("count failed: %v", err)
	}
	if n != 0 {
		t.Errorf("expected no media_access rows, got %d", n)
	}
}

// test a todo linking nonexistent media writes nothing: no todo and no access row
func TestCreateTodoItem_InvalidMediaWritesNothing(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := uuid.NewString()
	missingMediaID := uuid.NewString()
	if _, err := repo.CreateTodoItem(ctx, studentID, uuid.NewString(), models.CreateTodoItemRequestBody{
		TodoDescription: "Watch a video that doesn't exist",
		MediaID:         &missingMediaID,
		LinkType:        &moduleLink,
	}); err == nil {
		t.Fatal("expected an error for nonexistent media, got nil")
	}

	var todos int
	if err := db.QueryRow(ctx, `SELECT count(*) FROM public.todo_items WHERE student_id = $1`, studentID).Scan(&todos); err != nil {
		t.Fatalf("count failed: %v", err)
	}
	if todos != 0 {
		t.Errorf("expected no todo rows after a failed create, got %d", todos)
	}
	if got := countMediaAccess(t, db, studentID, missingMediaID); got != 0 {
		t.Errorf("expected no media_access rows after a failed create, got %d", got)
	}
}
