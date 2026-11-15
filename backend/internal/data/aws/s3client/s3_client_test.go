package s3client

import (
	"context"
	"net/url"
	"slices"
	"strings"
	"testing"
)

// presigning is local signing only, so fake credentials are enough
func TestPresignUpload_SignsContentTypeAndLength(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	t.Setenv("AWS_REGION", "us-east-1")

	client, err := NewClient(context.Background(), "test-bucket")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	uploadURL, err := client.PresignUpload(context.Background(), "videos/a.mp4", "video/mp4", 1024)
	if err != nil {
		t.Fatalf("PresignUpload failed: %v", err)
	}

	parsed, err := url.Parse(uploadURL)
	if err != nil {
		t.Fatalf("invalid presigned URL %q: %v", uploadURL, err)
	}

	signed := strings.Split(parsed.Query().Get("X-Amz-SignedHeaders"), ";")
	for _, header := range []string{"content-type", "content-length"} {
		if !slices.Contains(signed, header) {
			t.Errorf("expected %q in signed headers, got %v", header, signed)
		}
	}
}
