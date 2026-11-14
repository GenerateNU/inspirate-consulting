package pfpRepository

import (
	"fmt"
	"maps"
	"path"
	"slices"
	"strings"

	"github.com/google/uuid"

	"inspirate-consulting/internal/data/aws/s3client"
)

// every profile picture lives under pfps/{userID}/
const pfpPrefix = "pfps/"

// 5 MiB, kept in sync with the maximum on PresignPfpUploadRequestBody
const maxPfpBytes = 5 << 20

// allowed image types and the file extension each one is stored with
var pfpExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

type PfpRepository struct {
	s3 *s3client.Client
}

func NewPfpRepository(s3 *s3client.Client) *PfpRepository {
	return &PfpRepository{s3: s3}
}

func userPrefix(userID uuid.UUID) string {
	return fmt.Sprintf("%s%s/", pfpPrefix, userID)
}

// only keys this store generates match, so no "..", extra "/", or other user's segment gets through
func isUserPfpKey(userID uuid.UUID, key string) bool {
	name, ok := strings.CutPrefix(key, userPrefix(userID))
	if !ok {
		return false
	}
	ext := path.Ext(name)
	if !slices.Contains(slices.Collect(maps.Values(pfpExtensions)), ext) {
		return false
	}
	base := strings.TrimSuffix(name, ext)
	id, err := uuid.Parse(base)
	return err == nil && id.String() == base
}
