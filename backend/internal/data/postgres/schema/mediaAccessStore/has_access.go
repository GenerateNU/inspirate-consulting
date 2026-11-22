package mediaAccessRepository

import (
	"context"
	"inspirate-consulting/internal/data/postgres/schema"
)

// reports whether the student has been granted access to the media
func (r *MediaAccessRepository) HasMediaAccess(ctx context.Context, studentID string, mediaID string) (bool, error) {
	query, err := schema.ReadSQLBaseScript("has_media_access.sql", SqlMediaAccessFiles)
	if err != nil {
		return false, err
	}

	var hasAccess bool
	if err := r.db.QueryRow(ctx, query, studentID, mediaID).Scan(&hasAccess); err != nil {
		return false, err
	}
	return hasAccess, nil
}

// reports whether the S3 key belongs to media the student has been granted access to
// necessary to validate /videos/lookup endpoint, which only has the S3 key and not the media ID
func (r *MediaAccessRepository) HasS3KeyAccess(ctx context.Context, studentID string, s3Key string) (bool, error) {
	query, err := schema.ReadSQLBaseScript("has_s3_key_access.sql", SqlMediaAccessFiles)
	if err != nil {
		return false, err
	}

	var hasAccess bool
	if err := r.db.QueryRow(ctx, query, studentID, s3Key).Scan(&hasAccess); err != nil {
		return false, err
	}
	return hasAccess, nil
}
