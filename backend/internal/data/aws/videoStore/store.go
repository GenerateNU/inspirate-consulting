package videoRepository

import (
	"inspirate-consulting/internal/data/aws/s3client"
)

type VideoRepository struct {
	s3 *s3client.Client
}

func NewVideoRepository(s3 *s3client.Client) *VideoRepository {
	return &VideoRepository{s3: s3}
}
