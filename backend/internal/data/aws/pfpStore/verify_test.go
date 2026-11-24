package pfpRepository

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/sethvargo/go-envconfig"

	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data/aws/s3client"
)

// connects to the S3 bucket
// requires valid AWS credentials to be configured
func setupTestRepo(t *testing.T) (*PfpRepository, *s3client.Client, string) {
	t.Helper()

	var s3Config config.S3
	if err := envconfig.Process(context.Background(), &s3Config); err != nil {
		t.Fatalf("failed to load S3 config: %v", err)
	}

	client, err := s3client.NewClient(context.Background(), s3Config.Bucket)
	if err != nil {
		t.Fatalf("failed to create S3 client: %v", err)
	}

	return NewPfpRepository(client), client, s3Config.Bucket
}

// uploads directly with the SDK (bypassing the presigned URL), for test setup only
func seedTestObject(t *testing.T, client *s3client.Client, bucket string, key string, contentType string) {
	t.Helper()
	_, err := client.Raw().PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      awssdk.String(bucket),
		Key:         awssdk.String(key),
		Body:        bytes.NewReader([]byte("fake image")),
		ContentType: awssdk.String(contentType),
	})
	if err != nil {
		t.Fatalf("failed to seed test object key=%s: %v", key, err)
	}
	t.Cleanup(func() {
		if err := client.DeleteObject(context.Background(), key); err != nil {
			t.Logf("cleanup failed for key=%s: %v", key, err)
		}
	})
}

func TestVerifyUpload(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, client, bucket := setupTestRepo(t)
	ctx := context.Background()

	t.Run("accepts the user's uploaded image", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()
		key := userPrefix(userID) + uuid.NewString() + ".jpg"
		seedTestObject(t, client, bucket, key, "image/jpeg")

		if err := repo.VerifyUpload(ctx, userID, key); err != nil {
			t.Fatalf("VerifyUpload failed: %v", err)
		}
	})

	t.Run("missing upload is not found", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()
		assertHTTPStatus(t, repo.VerifyUpload(ctx, userID, userPrefix(userID)+uuid.NewString()+".jpg"), http.StatusNotFound)
	})

	t.Run("non-image upload is rejected", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()
		// a .jpg key holding HTML, as if the upload bypassed the signed Content-Type
		key := userPrefix(userID) + uuid.NewString() + ".jpg"
		seedTestObject(t, client, bucket, key, "text/html")

		assertHTTPStatus(t, repo.VerifyUpload(ctx, userID, key), http.StatusBadRequest)
	})
}
