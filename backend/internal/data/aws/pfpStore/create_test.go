package pfpRepository

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"

	"inspirate-consulting/internal/models"
)

func TestPresignUpload_KeyIsUnderUsersPrefix(t *testing.T) {
	repo := setupOfflineRepo(t)
	userID := uuid.New()

	result, err := repo.PresignUpload(context.Background(), userID, models.PresignPfpUploadRequestBody{
		ContentType:   "image/png",
		ContentLength: 2048,
	})
	if err != nil {
		t.Fatalf("PresignUpload failed: %v", err)
	}

	if !strings.HasPrefix(result.S3Key, "pfps/"+userID.String()+"/") {
		t.Errorf("expected key under the user's prefix, got %q", result.S3Key)
	}
	if !strings.HasSuffix(result.S3Key, ".png") {
		t.Errorf("expected a .png key, got %q", result.S3Key)
	}
	if !isUserPfpKey(userID, result.S3Key) {
		t.Errorf("expected the generated key to pass the ownership check, got %q", result.S3Key)
	}

	parsed, err := url.Parse(result.UploadURL)
	if err != nil {
		t.Fatalf("invalid upload URL %q: %v", result.UploadURL, err)
	}
	if signed := parsed.Query().Get("X-Amz-SignedHeaders"); !strings.Contains(signed, "content-type") {
		t.Errorf("expected content-type to be signed, got %q", signed)
	}
}

func TestPresignUpload_RejectsBadInput(t *testing.T) {
	repo := setupOfflineRepo(t)

	cases := map[string]models.PresignPfpUploadRequestBody{
		"non-image type":  {ContentType: "video/mp4", ContentLength: 2048},
		"svg":             {ContentType: "image/svg+xml", ContentLength: 2048},
		"zero length":     {ContentType: "image/jpeg", ContentLength: 0},
		"over 5 MiB":      {ContentType: "image/jpeg", ContentLength: maxPfpBytes + 1},
		"negative length": {ContentType: "image/jpeg", ContentLength: -1},
	}
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := repo.PresignUpload(context.Background(), uuid.New(), input)
			assertHTTPStatus(t, err, http.StatusBadRequest)
		})
	}
}
