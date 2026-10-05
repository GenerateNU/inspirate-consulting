package userRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

func TestFetchUserBySupabaseID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	supabaseID := uuid.New()
	key := "pfp-key"

	input := models.CreateUserInput{}
	input.Body.Name = "FetchMeBySupabase"
	input.Body.Email = "fetchmebysupabase@gmail.com"
	input.Body.PfpKey = &key

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	created, err := repo.CreateUser(ctx, input, supabaseID)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	output, err := repo.FetchUserBySupabaseID(ctx, models.FetchUserBySupabaseIDInput{SupabaseID: supabaseID})
	if err != nil {
		t.Fatalf("FetchUserBySupabaseID failed: %v", err)
	}

	if output.Body.ID != created.Body.User.ID {
		t.Errorf("expected id %q, got %q", created.Body.User.ID, output.Body.ID)
	}
	if output.Body.SupabaseID != supabaseID {
		t.Errorf("expected supabase_id %q, got %q", supabaseID, output.Body.SupabaseID)
	}
}

func TestFetchUserBySupabaseID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	_, err := repo.FetchUserBySupabaseID(ctx, models.FetchUserBySupabaseIDInput{SupabaseID: uuid.New()})
	if err == nil {
		t.Fatal("expected error fetching non-existent user, got nil")
	}
}
