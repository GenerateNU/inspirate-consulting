package pfp

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// presigns an upload under the logged-in user's own prefix; the key is saved only after UpdatePfp
func (h *Handler) PresignUpload(ctx context.Context, input models.PresignPfpUploadRequestBody) (*models.PresignUploadResponse, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, errs.Unauthorized()
	}

	return h.PfpRepository.PresignUpload(ctx, userID, input)
}
