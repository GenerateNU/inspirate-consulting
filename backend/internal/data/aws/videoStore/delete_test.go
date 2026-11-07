package videoRepository

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"inspirate-consulting/internal/errs"
)

func TestDeleteVideo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, client, bucket := setupTestRepo(t)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		key := "videos/test-delete-video.mp4"
		seedTestObject(t, client, bucket, key, []byte("fake video content"))
		// in case the delete under test fails, still remove the object
		cleanupObject(t, client, key)

		if err := repo.DeleteVideo(ctx, key); err != nil {
			t.Fatalf("DeleteVideo failed: %v", err)
		}

		exists, err := client.ObjectExists(ctx, key)
		if err != nil {
			t.Fatalf("ObjectExists failed: %v", err)
		}
		if exists {
			t.Errorf("expected object %q to be deleted, but it still exists", key)
		}
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		key := "videos/test-delete-video-does-not-exist.mp4"

		err := repo.DeleteVideo(ctx, key)
		if err == nil {
			t.Fatal("expected an error when deleting a nonexistent video, got nil")
		}

		var httpErr errs.HTTPError
		if !errors.As(err, &httpErr) {
			t.Fatalf("expected errs.HTTPError, got %T: %v", err, err)
		}
		if httpErr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, httpErr.Code)
		}
	})
}
