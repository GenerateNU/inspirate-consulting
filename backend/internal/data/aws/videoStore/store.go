package videoRepository

import (
	"strings"

	"inspirate-consulting/internal/data/aws/s3client"
)

// every video lives under this prefix; other prefixes in the bucket are off limits to this store
const videoPrefix = "videos/"

type VideoRepository struct {
	s3 *s3client.Client
}

func NewVideoRepository(s3 *s3client.Client) *VideoRepository {
	return &VideoRepository{s3: s3}
}

func isVideoKey(key string) bool {
	return strings.HasPrefix(key, videoPrefix) && len(key) > len(videoPrefix)
}
