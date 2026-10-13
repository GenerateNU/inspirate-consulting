package authtest

import (
	"context"
	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

var _ data.UserRepository = (*MemUserRepo)(nil)

// MemUserRepo is an in-memory UserRepository that mirrors the Postgres one,
// including new users starting with reset_time set (password reset required).
type MemUserRepo struct {
	mu    sync.Mutex
	users map[uuid.UUID]*models.User // by users.id
}

func NewMemUserRepo() *MemUserRepo {
	return &MemUserRepo{users: map[uuid.UUID]*models.User{}}
}

// Seed inserts a user for an existing auth account, with no pending reset.
func (r *MemUserRepo) Seed(name string, supabaseID uuid.UUID) *models.User {
	r.mu.Lock()
	defer r.mu.Unlock()
	u := &models.User{ID: uuid.New(), Name: name, SupabaseID: supabaseID}
	r.users[u.ID] = u
	return copyUser(u)
}

// BySupabaseID returns a snapshot of the stored user for an auth account.
func (r *MemUserRepo) BySupabaseID(id uuid.UUID) (*models.User, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.users {
		if u.SupabaseID == id {
			return copyUser(u), true
		}
	}
	return nil, false
}

func (r *MemUserRepo) Count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.users)
}

func (r *MemUserRepo) CreateUser(_ context.Context, in models.CreateUserInput, supabaseID uuid.UUID) (*models.CreateUserOutput, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	resetTime := time.Now().Add(-1 * time.Second)
	u := &models.User{
		ID:         uuid.New(),
		Name:       in.Body.Name,
		SupabaseID: supabaseID,
		PfpKey:     in.Body.PfpKey,
		ResetTime:  &resetTime,
	}
	r.users[u.ID] = u
	return &models.CreateUserOutput{Body: &models.CreateUserBody{User: copyUser(u)}}, nil
}

func (r *MemUserRepo) FetchUser(_ context.Context, in models.FetchUserInput) (*models.FetchUserOutput, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[in.ID]
	if !ok {
		return nil, errs.NotFound("User", "id", in.ID)
	}
	return &models.FetchUserOutput{Body: copyUser(u)}, nil
}

func (r *MemUserRepo) FetchUserBySupabaseID(_ context.Context, in models.FetchUserBySupabaseIDInput) (*models.FetchUserOutput, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.users {
		if u.SupabaseID == in.SupabaseID {
			return &models.FetchUserOutput{Body: copyUser(u)}, nil
		}
	}
	return nil, errs.NotFound("User", "supabase_id", in.SupabaseID)
}

func (r *MemUserRepo) UpdateResetTime(_ context.Context, supabaseID uuid.UUID, resetTime *time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.users {
		if u.SupabaseID == supabaseID {
			u.ResetTime = resetTime
			return nil
		}
	}
	// the Postgres UPDATE is a no-op for an unknown id; mirror that
	return nil
}

func copyUser(u *models.User) *models.User {
	c := *u
	return &c
}
