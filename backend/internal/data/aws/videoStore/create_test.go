package videoRepository

import (
	"context"
	"strings"
	"testing"

	"inspirate-consulting/internal/models"
)

func TestPresignUpload(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _, _ := setupTestRepo(t)
	ctx := context.Background()

	result, err := repo.PresignUpload(ctx, models.PresignUploadRequestBody{
		OriginalFilename: "test-video.mp4",
		ContentType:      "video/mp4",
		ContentLength:    1024,
	})
	if err != nil {
		t.Fatalf("PresignUpload failed: %v", err)
	}

	if result.S3Key == "" {
		t.Error("expected a non-empty S3 key")
	}
	if !strings.HasPrefix(result.S3Key, "videos/") {
		t.Errorf("expected key to start with 'videos/', got %q", result.S3Key)
	}
	if !strings.HasSuffix(result.S3Key, "test-video.mp4") {
		t.Errorf("expected key to end with the original filename, got %q", result.S3Key)
	}
	if result.UploadURL == "" {
		t.Error("expected a non-empty upload URL")
	}
	if !strings.Contains(result.UploadURL, result.S3Key) {
		t.Errorf("expected upload URL to reference the generated key, got %q", result.UploadURL)
	}
}
