package pfpRepository

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"

	"inspirate-consulting/internal/data/aws/s3client"
	"inspirate-consulting/internal/errs"
)

// these cases fail or presign before any S3 call, so fake credentials are enough
func setupOfflineRepo(t *testing.T) *PfpRepository {
	t.Helper()
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	t.Setenv("AWS_REGION", "us-east-1")

	client, err := s3client.NewClient(context.Background(), "test-bucket")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	return NewPfpRepository(client)
}

func assertHTTPStatus(t *testing.T, err error, status int) {
	t.Helper()
	var httpErr errs.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected errs.HTTPError, got %T: %v", err, err)
	}
	if httpErr.Code != status {
		t.Errorf("expected status %d, got %d", status, httpErr.Code)
	}
}

func TestPfpRepository_RejectsOtherUsersKeys(t *testing.T) {
	repo := setupOfflineRepo(t)
	userID := uuid.New()

	keys := map[string]string{
		"another user's picture": "pfps/" + uuid.NewString() + "/a.jpg",
		"a video":                "videos/uuid-tips.mp4",
		"the user's bare prefix": "pfps/" + userID.String() + "/",
		"dot-dot traversal":      "pfps/" + userID.String() + "/../" + uuid.NewString() + "/" + uuid.NewString() + ".jpg",
		"nested segment":         "pfps/" + userID.String() + "/x/" + uuid.NewString() + ".jpg",
		"non-uuid name":          "pfps/" + userID.String() + "/a.jpg",
		"disallowed extension":   "pfps/" + userID.String() + "/" + uuid.NewString() + ".html",
		"prefix without slash":   "pfps/" + userID.String() + "a.jpg",
		"empty":                  "",
	}
	for name, key := range keys {
		t.Run("verify "+name, func(t *testing.T) {
			assertHTTPStatus(t, repo.VerifyUpload(context.Background(), userID, key), http.StatusForbidden)
		})
		t.Run("delete "+name, func(t *testing.T) {
			assertHTTPStatus(t, repo.DeletePfp(context.Background(), userID, key), http.StatusForbidden)
		})
	}
}
