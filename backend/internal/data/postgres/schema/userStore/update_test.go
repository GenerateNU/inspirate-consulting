package userRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func TestUpdateUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	key := "old-key"
	input := models.CreateUserInput{}
	input.Body.Name = "Old Name"
	input.Body.Email = "updateme@gmail.com"
	input.Body.PfpKey = &key
	input.Body.Role = models.CounselorRole

	created, err := repo.CreateUser(ctx, input, uuid.New())
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	name := "New Name"
	newKey := "new-key"
	updated, err := repo.UpdateUser(ctx, created.Body.User.ID, models.UpdateUserBody{Name: &name, PfpKey: &newKey})
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	if updated.Name != name {
		t.Errorf("expected name %q, got %q", name, updated.Name)
	}
	if updated.PfpKey == nil || *updated.PfpKey != newKey {
		t.Errorf("expected pfp_key %q, got %v", newKey, updated.PfpKey)
	}
	if updated.Email != input.Body.Email {
		t.Errorf("expected email %q to be unchanged, got %q", input.Body.Email, updated.Email)
	}
}

// fields left nil keep their current values
func TestUpdateUser_Partial(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	key := "kept-key"
	input := models.CreateUserInput{}
	input.Body.Name = "Old Name"
	input.Body.Email = "partial@gmail.com"
	input.Body.PfpKey = &key
	input.Body.Role = models.CounselorRole

	created, err := repo.CreateUser(ctx, input, uuid.New())
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	name := "New Name"
	updated, err := repo.UpdateUser(ctx, created.Body.User.ID, models.UpdateUserBody{Name: &name})
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	if updated.Name != name {
		t.Errorf("expected name %q, got %q", name, updated.Name)
	}
	if updated.PfpKey == nil || *updated.PfpKey != key {
		t.Errorf("expected pfp_key %q to be unchanged, got %v", key, updated.PfpKey)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	name := "Nobody"
	_, err := repo.UpdateUser(ctx, uuid.New(), models.UpdateUserBody{Name: &name})
	if err == nil {
		t.Fatal("expected error updating non-existent user, got nil")
	}
	httpErr, ok := err.(errs.HTTPError)
	if !ok || httpErr.Code != 404 {
		t.Errorf("expected 404 HTTPError, got %v", err)
	}
}
