package pfpRepository

import (
	"context"
	"strings"

	"inspirate-consulting/internal/errs"
)

// any user's picture can be viewed, so this only checks that the key is a profile picture
func (r *PfpRepository) GetDownloadURL(ctx context.Context, s3Key string) (string, error) {
	if !strings.HasPrefix(s3Key, pfpPrefix) {
		return "", errs.NotFound("profile picture", "s3Key", s3Key)
	}
	return r.s3.PresignDownload(ctx, s3Key)
}
