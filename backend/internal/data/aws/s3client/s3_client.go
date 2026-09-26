package s3client

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)


const (
	// how long a presigned upload (PUT) URL remains valid.
	UploadURLExpiry = 15 * time.Minute

	// how long a presigned download (GET) URL remains valid.
	DownloadURLExpiry = 1 * time.Hour
)

// wraps AWS S3 client for a single bucket
// methods for: generating/getting presigned URLs, generating object keys, listing objects, and deleting objects
type Client struct {
	s3Client *s3.Client
	presignClient *s3.PresignClient
	bucket string
}

// loads AWS config using the SDK's default credential chain (.env vars, ~/.aws/credentials & ~/.aws/config, etc.)
// and returns a Client scoped to the given bucket.
func NewClient(ctx context.Context, bucket string) (*Client, error) {
	if bucket == "" {
		return nil, fmt.Errorf("s3client: bucket name is required")
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("s3client: failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3Client)

	return &Client{
		s3Client: s3Client,
		presignClient: presignClient,
		bucket: bucket,
	}, nil
}

// creates a unique S3 object key for a new upload
// preserves the original filename for readability
func GenerateObjectKey(originalFilename string) string {
	return fmt.Sprintf("videos/%s-%s", uuid.NewString(), originalFilename)
}

// returns a presigned URL for the given object key (PUT: new upload).
func (c *Client) PresignUpload(ctx context.Context, key string) (string, error) {
	request, err := c.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key: aws.String(key),
	}, s3.WithPresignExpires(UploadURLExpiry))
	if err != nil {
		return "", fmt.Errorf("s3client: failed to upload object: %w", err)
	}
	return request.URL, nil
}

// returns a presigned URL for the given object key (GET: retrieve upload).
func (c *Client) PresignDownload(ctx context.Context, key string) (string, error) {
	request, err := c.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key: aws.String(key),
	}, s3.WithPresignExpires(DownloadURLExpiry))
	if err != nil {
		return "", fmt.Errorf("s3client: failed to retrieve object: %w", err)
	}
	return request.URL, nil
}

// returns the key of every object currently in the bucket.
func (c *Client) ListObjectKeys(ctx context.Context) ([]string, error) {
	output, err := c.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("s3client: failed to list objects: %w", err)
	}

	keys := make([]string, 0, len(output.Contents))
	for _, object := range output.Contents {
		keys = append(keys, aws.ToString(object.Key))
	}
	return keys, nil
}

// removes an object from the bucket.
func (c *Client) DeleteObject(ctx context.Context, key string) error {
	_, err := c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key: aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("s3client: failed to delete object: %w", err)
	}
	return nil
}

// returns the underlying AWS S3 client, for cases (mainly testing) that need direct S3 access
func (c *Client) Raw() *s3.Client {
	return c.s3Client
}