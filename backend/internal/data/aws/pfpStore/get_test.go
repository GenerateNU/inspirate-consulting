package pfpRepository

import (
	"context"
	"net/http"
	"testing"
)

func TestGetDownloadURL_RejectsNonPfpKeys(t *testing.T) {
	repo := setupOfflineRepo(t)

	_, err := repo.GetDownloadURL(context.Background(), "videos/uuid-tips.mp4")
	assertHTTPStatus(t, err, http.StatusNotFound)
}
