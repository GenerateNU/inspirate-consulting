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

	input := models.CreateUserInput{}
	input.Body.Name = "Old Name"
	input.Body.Email = "updateme@gmail.com"
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

	input := models.CreateUserInput{}
	input.Body.Name = "Old Name"
	input.Body.Email = "partial@gmail.com"
	input.Body.Role = models.CounselorRole

	created, err := repo.CreateUser(ctx, input, uuid.New())
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	updated, err := repo.UpdateUser(ctx, created.Body.User.ID, models.UpdateUserBody{})
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	if updated.Name != input.Body.Name {
		t.Errorf("expected name %q to be unchanged, got %q", input.Body.Name, updated.Name)
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
