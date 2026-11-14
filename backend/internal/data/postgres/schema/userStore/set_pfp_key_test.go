package userRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func TestSetPfpKey(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	input := models.CreateUserInput{}
	input.Body.Name = "Pfp User"
	input.Body.Email = "pfpuser@gmail.com"
	input.Body.Role = models.CounselorRole

	created, err := repo.CreateUser(ctx, input, uuid.New())
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	id := created.Body.User.ID

	// first picture: nothing is replaced
	first := "pfps/user/first.jpg"
	updated, previous, err := repo.SetPfpKey(ctx, id, &first)
	if err != nil {
		t.Fatalf("SetPfpKey failed: %v", err)
	}
	if updated.PfpKey == nil || *updated.PfpKey != first {
		t.Errorf("expected pfp_key %q, got %v", first, updated.PfpKey)
	}
	if previous != nil {
		t.Errorf("expected no previous key, got %q", *previous)
	}

	// replacing returns the old key
	second := "pfps/user/second.png"
	updated, previous, err = repo.SetPfpKey(ctx, id, &second)
	if err != nil {
		t.Fatalf("SetPfpKey failed: %v", err)
	}
	if updated.PfpKey == nil || *updated.PfpKey != second {
		t.Errorf("expected pfp_key %q, got %v", second, updated.PfpKey)
	}
	if previous == nil || *previous != first {
		t.Errorf("expected previous key %q, got %v", first, previous)
	}

	// nil clears it and returns the old key
	updated, previous, err = repo.SetPfpKey(ctx, id, nil)
	if err != nil {
		t.Fatalf("SetPfpKey failed: %v", err)
	}
	if updated.PfpKey != nil {
		t.Errorf("expected pfp_key to be cleared, got %q", *updated.PfpKey)
	}
	if previous == nil || *previous != second {
		t.Errorf("expected previous key %q, got %v", second, previous)
	}
}

func TestSetPfpKey_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)

	key := "pfps/nobody/a.jpg"
	_, _, err := repo.SetPfpKey(context.Background(), uuid.New(), &key)
	if err == nil {
		t.Fatal("expected error for non-existent user, got nil")
	}
	httpErr, ok := err.(errs.HTTPError)
	if !ok || httpErr.Code != 404 {
		t.Errorf("expected 404 HTTPError, got %v", err)
	}
}
