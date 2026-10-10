package userRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

func TestCreateUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	supabaseID := uuid.New()
	key := "pfp-key"

	input := models.CreateUserInput{}
	input.Body.Name = "Aleng123"
	input.Body.Email = "tt@gmail.com"
	input.Body.PfpKey = &key

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	output, err := repo.CreateUser(ctx, input, supabaseID)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if output.Body.User.Name != input.Body.Name {
		t.Errorf("expected name %q, got %q", input.Body.Name, output.Body.User.Name)
	}
	if output.Body.User.SupabaseID != supabaseID {
		t.Errorf("expected supabase_id %q, got %q", supabaseID, output.Body.User.SupabaseID)
	}
	if output.Body.User.PfpKey == nil || *output.Body.User.PfpKey != key {
		t.Errorf("expected pfp_key %q, got %v", key, output.Body.User.PfpKey)
	}
}

func TestCreateUserNullPfpKey(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	supabaseID := uuid.New()

	input := models.CreateUserInput{}
	input.Body.Name = "NoPfpUser"
	input.Body.Email = "nopfp@gmail.com"
	input.Body.PfpKey = nil

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	output, err := repo.CreateUser(ctx, input, supabaseID)
	if err != nil {
		t.Fatalf("CreateUser with nil pfp_key failed: %v", err)
	}

	if output.Body.User.PfpKey != nil {
		t.Errorf("expected nil pfp_key, got %q", *output.Body.User.PfpKey)
	}
}

func TestCreateUser_Duplicate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	supabaseID := uuid.New()

	input := models.CreateUserInput{}
	input.Body.Name = "DupeUser"
	input.Body.Email = "dupe@gmail.com"

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	_, err := repo.CreateUser(ctx, input, supabaseID)
	if err != nil {
		t.Fatalf("first CreateUser failed: %v", err)
	}

	_, err = repo.CreateUser(ctx, input, supabaseID)
	if err != nil {
		t.Fatal("expected error creating duplicate user with same supabase_id, got nil")
	}
}
