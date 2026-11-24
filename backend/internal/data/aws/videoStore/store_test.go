package videoRepository

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"inspirate-consulting/internal/data/aws/s3client"
	"inspirate-consulting/internal/errs"
)

// keys outside videos/ are rejected before any S3 call, so fake credentials are enough
func TestVideoRepository_RejectsKeysOutsideVideoPrefix(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	t.Setenv("AWS_REGION", "us-east-1")

	client, err := s3client.NewClient(context.Background(), "test-bucket")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	repo := NewVideoRepository(client)
	ctx := context.Background()

	keys := []string{"pfps/some-user/avatar.jpg", "videos/", "", "common-app-tips.mp4", "videosx/a.mp4"}
	for _, key := range keys {
		t.Run("get "+key, func(t *testing.T) {
			_, err := repo.GetVideo(ctx, key)
			assertNotFound(t, err)
		})
		t.Run("delete "+key, func(t *testing.T) {
			assertNotFound(t, repo.DeleteVideo(ctx, key))
		})
	}
}

func assertNotFound(t *testing.T, err error) {
	t.Helper()
	var httpErr errs.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected errs.HTTPError, got %T: %v", err, err)
	}
	if httpErr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, httpErr.Code)
	}
}
