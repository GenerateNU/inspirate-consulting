package aws

import (
	"context"
	"log"

	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/data/aws/s3client"
	videoRepository "inspirate-consulting/internal/data/aws/videoStore"
)

// NewRepository creates a new Repository instance connected to S3
// at the top level of the the aws package to prevent import loops between the repository implmentation and the client.
func NewRepository(ctx context.Context, s3Config config.S3) *data.Repository {
	client, err := s3client.NewClient(ctx, s3Config.Bucket)
	if err != nil {
		log.Fatalf("Failed to connect to S3: %v", err)
	}
	return &data.Repository{
		Video: videoRepository.NewVideoRepository(client),
	}
}