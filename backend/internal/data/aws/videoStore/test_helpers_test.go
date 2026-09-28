package videoRepository

import (
	"bytes"
	"context"
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sethvargo/go-envconfig"

	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data/aws/s3client"
)

// connects to the S3 bucket
// requires valid AWS credentials to be configured
func setupTestRepo(t *testing.T) (*VideoRepository, *s3client.Client, string) {
	t.Helper()

	var s3Config config.S3
	if err := envconfig.Process(context.Background(), &s3Config); err != nil {
		t.Fatalf("failed to load S3 config: %v", err)
	}

	client, err := s3client.NewClient(context.Background(), s3Config.Bucket)
	if err != nil {
		t.Fatalf("failed to create S3 client: %v", err)
	}

	return NewVideoRepository(client), client, s3Config.Bucket
}

// seedTestObject uploads raw bytes directly via the underlying AWS SDK
// client (bypassing the presigned-URL flow), for test setup only.
func seedTestObject(t *testing.T, client *s3client.Client, bucket string, key string, body []byte) {
	t.Helper()
	_, err := client.Raw().PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: awssdk.String(bucket),
		Key:    awssdk.String(key),
		Body:   bytes.NewReader(body),
	})
	if err != nil {
		t.Fatalf("failed to seed test object key=%s: %v", key, err)
	}
}

// cleanupObject deletes a test object from the bucket after the test completes.
func cleanupObject(t *testing.T, client *s3client.Client, key string) {
	t.Helper()
	t.Cleanup(func() {
		if err := client.DeleteObject(context.Background(), key); err != nil {
			t.Logf("cleanup failed for key=%s: %v", key, err)
		}
	})
}