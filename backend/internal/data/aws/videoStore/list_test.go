package videoRepository

import (
	"context"
	"testing"
)

// test list function returns all videos currently in the bucket
// intentionally not run in parallel to avoid race conditions with other tests
func TestListVideos(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	repo, client, bucket := setupTestRepo(t)
	ctx := context.Background()

	keyA := "videos/test-list-a.mp4"
	keyB := "videos/test-list-b.mp4"

	seedTestObject(t, client, bucket, keyA, []byte("fake video a"))
	cleanupObject(t, client, keyA)

	seedTestObject(t, client, bucket, keyB, []byte("fake video b"))
	cleanupObject(t, client, keyB)

	videos, err := repo.ListVideos(ctx)
	if err != nil {
		t.Fatalf("ListVideos failed: %v", err)
	}

	foundKeys := map[string]bool{}
	for _, v := range videos {
		foundKeys[v.S3Key] = true
		if v.DownloadURL == "" {
			t.Errorf("expected a non-empty URL for key=%s", v.S3Key)
		}
	}
	if !foundKeys[keyA] || !foundKeys[keyB] {
		t.Errorf("expected both seeded keys in the list, found: %v", foundKeys)
	}
}
