package video

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// GetVideo retrieves a video from the S3 bucket.
// Counselors can fetch any key; students only keys that belong to media they have been granted access to.
func (h *Handler) GetVideo(ctx context.Context, s3Key string) (*models.Video, error) {
	if studentID, err := auth.GetStudentID(ctx); err == nil {
		hasAccess, err := h.MediaAccessRepository.HasS3KeyAccess(ctx, studentID.String(), s3Key)
		if err != nil {
			return nil, err
		}
		if !hasAccess {
			return nil, errs.NotFound("video", "s3_key", s3Key)
		}
	} else if _, err := auth.GetCounselorID(ctx); err != nil {
		return nil, errs.Forbidden("only students and counselors can view videos")
	}

	video, err := h.VideoRepository.GetVideo(ctx, s3Key)
	if err != nil {
		return nil, err
	}

	return video, nil
}
