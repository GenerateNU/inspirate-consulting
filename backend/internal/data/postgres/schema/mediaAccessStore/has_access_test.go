package mediaAccessRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// test HasMediaAccess is true only for the exact student + media pair that was granted
func TestHasMediaAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewMediaAccessRepository(db)
	ctx := context.Background()

	studentID := uuid.NewString()
	otherStudentID := uuid.NewString()
	grantedMedia := createTestMedia(t, db, "Granted Video")
	otherMedia := createTestMedia(t, db, "Not Granted Video")

	if _, err := repo.GrantMediaAccess(ctx, &models.GrantMediaAccessRequestBody{StudentID: studentID, MediaID: grantedMedia}); err != nil {
		t.Fatalf("setup GrantMediaAccess failed: %v", err)
	}

	tests := []struct {
		name      string
		studentID string
		mediaID   string
		want      bool
	}{
		{name: "granted student and media", studentID: studentID, mediaID: grantedMedia, want: true},
		{name: "same student, different media", studentID: studentID, mediaID: otherMedia, want: false},
		{name: "different student, granted media", studentID: otherStudentID, mediaID: grantedMedia, want: false},
		{name: "nonexistent media", studentID: studentID, mediaID: uuid.NewString(), want: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := repo.HasMediaAccess(ctx, tc.studentID, tc.mediaID)
			if err != nil {
				t.Fatalf("HasMediaAccess failed: %v", err)
			}
			if got != tc.want {
				t.Errorf("expected %v, got %v", tc.want, got)
			}
		})
	}
}

// test HasMediaAccess is false after the grant is revoked
func TestHasMediaAccess_Revoked(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewMediaAccessRepository(db)
	ctx := context.Background()

	studentID := uuid.NewString()
	mediaID := createTestMedia(t, db, "Revoked Video")

	grant, err := repo.GrantMediaAccess(ctx, &models.GrantMediaAccessRequestBody{StudentID: studentID, MediaID: mediaID})
	if err != nil {
		t.Fatalf("setup GrantMediaAccess failed: %v", err)
	}
	if err := repo.RevokeMediaAccess(ctx, grant.ID); err != nil {
		t.Fatalf("setup RevokeMediaAccess failed: %v", err)
	}

	got, err := repo.HasMediaAccess(ctx, studentID, mediaID)
	if err != nil {
		t.Fatalf("HasMediaAccess failed: %v", err)
	}
	if got {
		t.Error("expected no access after revoke")
	}
}

// test HasS3KeyAccess is true only for keys of media the student was granted
func TestHasS3KeyAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewMediaAccessRepository(db)
	ctx := context.Background()

	studentID := uuid.NewString()
	otherStudentID := uuid.NewString()
	grantedKey := "videos/" + uuid.NewString() + "-granted.mp4"
	otherKey := "videos/" + uuid.NewString() + "-other.mp4"
	grantedMedia := createTestMediaWithKey(t, db, "Granted Video", grantedKey)
	createTestMediaWithKey(t, db, "Other Video", otherKey)

	if _, err := repo.GrantMediaAccess(ctx, &models.GrantMediaAccessRequestBody{StudentID: studentID, MediaID: grantedMedia}); err != nil {
		t.Fatalf("setup GrantMediaAccess failed: %v", err)
	}

	tests := []struct {
		name      string
		studentID string
		s3Key     string
		want      bool
	}{
		{name: "granted student and key", studentID: studentID, s3Key: grantedKey, want: true},
		{name: "same student, key of media not granted", studentID: studentID, s3Key: otherKey, want: false},
		{name: "different student, granted key", studentID: otherStudentID, s3Key: grantedKey, want: false},
		{name: "key not attached to any media", studentID: studentID, s3Key: "videos/orphan.mp4", want: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := repo.HasS3KeyAccess(ctx, tc.studentID, tc.s3Key)
			if err != nil {
				t.Fatalf("HasS3KeyAccess failed: %v", err)
			}
			if got != tc.want {
				t.Errorf("expected %v, got %v", tc.want, got)
			}
		})
	}
}
