package videoRepository

import (
	"context"
	"testing"
)

func TestGetVideo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, client, bucket := setupTestRepo(t)
	ctx := context.Background()

	key := "videos/test-get-video.mp4"
	seedTestObject(t, client, bucket, key, []byte("fake video content"))
	cleanupObject(t, client, key)

	video, err := repo.GetVideo(ctx, key)
	if err != nil {
		t.Fatalf("GetVideo failed: %v", err)
	}

	if video.S3Key != key {
		t.Errorf("expected S3Key %q, got %q", key, video.S3Key)
	}
	if video.DownloadURL == "" {
		t.Error("expected a non-empty presigned URL")
	}
}