package userRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

func TestFetchUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	supabaseID := uuid.New()
	key := "pfp-key"

	input := models.CreateUserInput{}
	input.Body.Name = "FetchMe"
	input.Body.Email = "fetchme@gmail.com"
	input.Body.Password = "2976$$Alen$$"
	input.Body.PfpKey = &key

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	created, err := repo.CreateUser(ctx, input, supabaseID)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	fetchInput := models.FetchUserInput{ID: created.Body.ID}
	output, err := repo.FetchUser(ctx, fetchInput)
	if err != nil {
		t.Fatalf("FetchUser failed: %v", err)
	}

	if output.Body.Name != input.Body.Name {
		t.Errorf("expected name %q, got %q", input.Body.Name, output.Body.Name)
	}
	if output.Body.SupabaseID != supabaseID {
		t.Errorf("expected supabase_id %q, got %q", supabaseID, output.Body.SupabaseID)
	}
}

func TestFetchUser_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	fetchInput := models.FetchUserInput{ID: uuid.New()}
	_, err := repo.FetchUser(ctx, fetchInput)
	if err == nil {
		t.Fatal("expected error fetching non-existent user, got nil")
	}
}
