package extracurricularRepository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"inspirate-consulting/internal/models"
)

func TestExtracurricularAchievements_CreateUpdateDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()
	parent := createTestExtracurricular(t, repo, uuid.NewString(), createTestUser(t, db), minimalCreateInput("Debate", "2025-01-01"))

	created, err := repo.CreateExtracurricularAchievement(ctx, parent.ID, models.CreateExtracurricularAchievementRequest{
		EntryDate: "2025-03-15", Body: "Won regionals",
	})
	if err != nil {
		t.Fatalf("CreateExtracurricularAchievement failed: %v", err)
	}
	if created.ExtracurricularID != parent.ID || created.Body != "Won regionals" ||
		!created.EntryDate.Equal(time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("created = %+v", *created)
	}

	t.Run("partial update keeps omitted fields", func(t *testing.T) {
		got, err := repo.UpdateExtracurricularAchievement(ctx, parent.ID, created.ID, models.UpdateExtracurricularAchievementRequest{
			Body: ptr("Won regionals (1st place)"),
		})
		if err != nil {
			t.Fatalf("UpdateExtracurricularAchievement failed: %v", err)
		}
		if got.Body != "Won regionals (1st place)" || !got.EntryDate.Equal(created.EntryDate) {
			t.Errorf("got = %+v", *got)
		}

		got, err = repo.UpdateExtracurricularAchievement(ctx, parent.ID, created.ID, models.UpdateExtracurricularAchievementRequest{
			EntryDate: ptr("2025-03-16"),
		})
		if err != nil {
			t.Fatalf("UpdateExtracurricularAchievement failed: %v", err)
		}
		if got.Body != "Won regionals (1st place)" || !got.EntryDate.Equal(time.Date(2025, 3, 16, 0, 0, 0, 0, time.UTC)) {
			t.Errorf("got = %+v", *got)
		}
	})

	t.Run("delete", func(t *testing.T) {
		if err := repo.DeleteExtracurricularAchievement(ctx, parent.ID, created.ID); err != nil {
			t.Fatalf("DeleteExtracurricularAchievement failed: %v", err)
		}
		if n := countRows(t, db, "extracurricular_achievements", parent.ID); n != 0 {
			t.Errorf("achievements remaining = %d, want 0", n)
		}
		assertNotFound(t, repo.DeleteExtracurricularAchievement(ctx, parent.ID, created.ID))
	})
}

// An achievement can't be edited or deleted through a different extracurricular's id.
func TestExtracurricularAchievements_ScopedToParent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()
	userID := createTestUser(t, db)
	owner := createTestExtracurricular(t, repo, uuid.NewString(), userID, minimalCreateInput("Owner", "2025-01-01"))
	other := createTestExtracurricular(t, repo, uuid.NewString(), userID, minimalCreateInput("Other", "2025-01-01"))
	achievementID := insertAchievement(t, db, owner.ID, "2025-02-01", "Mine")

	_, err := repo.UpdateExtracurricularAchievement(ctx, other.ID, achievementID, models.UpdateExtracurricularAchievementRequest{Body: ptr("Hijacked")})
	assertNotFound(t, err)
	assertNotFound(t, repo.DeleteExtracurricularAchievement(ctx, other.ID, achievementID))

	if n := countRows(t, db, "extracurricular_achievements", owner.ID); n != 1 {
		t.Errorf("owner's achievements = %d, want 1", n)
	}
}
