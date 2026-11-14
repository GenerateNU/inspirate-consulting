package pfp

import (
	"context"

	"github.com/google/uuid"

	"inspirate-consulting/internal/auth"
)

// userCtx returns a new user id and a context logged in as that user.
func userCtx() (uuid.UUID, context.Context) {
	userID := uuid.New()
	return userID, auth.WithUserID(context.Background(), userID)
}
